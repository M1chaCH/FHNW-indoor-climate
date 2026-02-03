package sensor_data

import (
	"fmt"
	"sensor_hub_backend/elastic"
	"sensor_hub_backend/elastic/buffer"
	"sensor_hub_backend/logs"
	"sensor_hub_backend/obs"
	"sensor_hub_backend/proto_types"
	"sensor_hub_backend/sql"
	"time"

	"github.com/eclipse/paho.golang/paho"
	"google.golang.org/protobuf/proto"
)

var Observable = obs.NewObservable[*proto_types.SensorData]("Sensor Data")

func HandleSensorDataReceived(p *paho.Publish) {
	sensorData := &proto_types.SensorData{}

	if err := proto.Unmarshal(p.Payload, sensorData); err != nil {
		logs.LogErr("Failed to unmarshal received sensor data", err)
		return
	}

	logs.LogInfo("Received sensor data: %s\n", sensorData)

	entity := createDeviceEntity(sensorData)
	sql.UpsertSensorReadingThrottled(&entity)

	docs := groupSensorReadings(sensorData)

	ok, err := sql.IsDeviceAuthorizedCached(entity.DeviceId)
	if err == nil && ok {
		elastic.SendSensorDataToElasticDebounced(docs)
	} else {
		if err != nil {
			logs.LogErr("Failed to check device authorization (sending data to buffer)", err)
		} else {
			logs.LogInfo("Device not authorized, sending data to buffer")
		}
		buffer.PutSensorDataToBuffer(docs)
	}
	Observable.Emit(sensorData)
}

func createDeviceEntity(sensorData *proto_types.SensorData) sql.DeviceEntity {
	entity := sql.DeviceEntity{
		DeviceId:        sensorData.DeviceId,
		Name:            *sensorData.DeviceName,
		LastIp:          *sensorData.Ip,
		LastReading:     "|",
		LastReadingTime: time.Time{},
		Authorized:      false,
	}

	firstTime := time.Time{}

	for _, measurement := range sensorData.Measurements {
		entity.LastReading += measurement.GetSensorValueName() + ":" + ReadMeasurementValueHumanized(measurement) + "|"

		measurementTime, err := parseTimestamp(measurement)
		if err != nil {
			continue
		}

		if firstTime.IsZero() || firstTime.After(measurementTime) {
			firstTime = measurementTime
		}
	}

	entity.LastReadingTime = firstTime

	return entity
}

func groupSensorReadings(sensorData *proto_types.SensorData) []*elastic.SensorDataDocument {
	sensorReadings := make(map[string]*elastic.SensorDataDocument)

	for _, measurement := range sensorData.GetMeasurements() {
		sensorReading := sensorReadings[measurement.GetSensorType()]
		if sensorReading == nil {
			t, err := parseTimestamp(measurement)
			if err != nil {
				logs.LogErrCustom("Failed to parse timestamp: %s, skipping\n", err)
				continue
			}

			sensorReading = &elastic.SensorDataDocument{
				DeviceId:   sensorData.GetDeviceId(),
				SensorType: measurement.GetSensorType(),
				Timestamp:  t,
				DeviceName: sensorData.GetDeviceName(),
				Values:     make(map[string]interface{}),
			}
			sensorReadings[measurement.GetSensorType()] = sensorReading
		}
		sensorReading.Values[measurement.GetSensorValueName()] = ReadMeasurementValue(measurement)
	}

	docs := make([]*elastic.SensorDataDocument, len(sensorReadings))
	i := 0
	for _, sensorReading := range sensorReadings {
		docs[i] = sensorReading
		i++
	}

	return docs
}

func ReadMeasurementValue(measurement *proto_types.SensorData_Measurement) interface{} {
	switch measurement.SensorValueType {
	case proto_types.SensorData_STRING:
		return measurement.GetStringValue()
	case proto_types.SensorData_BOOL:
		return measurement.GetFlagValue()
	case proto_types.SensorData_INT32:
		return measurement.GetIntValue()
	case proto_types.SensorData_DOUBLE:
		return measurement.GetDoubleValue()
	}

	return nil
}

func ReadMeasurementValueHumanized(measurement *proto_types.SensorData_Measurement) string {
	switch measurement.SensorValueType {
	case proto_types.SensorData_STRING:
		return measurement.GetStringValue()
	case proto_types.SensorData_BOOL:
		if measurement.GetFlagValue() {
			return "✅"
		}

		return "❌"
	case proto_types.SensorData_INT32:
		return fmt.Sprintf("%d", measurement.GetIntValue())
	case proto_types.SensorData_DOUBLE:
		return fmt.Sprintf("%.2f", measurement.GetDoubleValue())
	}

	return ""
}

func parseTimestamp(measurement *proto_types.SensorData_Measurement) (time.Time, error) {
	timeLayout := "2006-01-02T15:04:05"
	return time.Parse(timeLayout, measurement.GetReadTimestamp())
}
