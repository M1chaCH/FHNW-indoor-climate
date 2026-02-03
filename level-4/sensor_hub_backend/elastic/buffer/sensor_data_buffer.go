package buffer

import (
	"sensor_hub_backend/elastic"
	"sensor_hub_backend/logs"
)

const sensorDataBufferSizeLimit = 1000

var sensorDataBuffer = make(map[string][]*elastic.SensorDataDocument)
var sensorDataBufferOverflow = make(map[string]int)

func PutSensorDataToBuffer(typedDocs []*elastic.SensorDataDocument) {
	for _, doc := range typedDocs {
		v, ok := sensorDataBuffer[doc.DeviceId]
		if !ok {
			v = make([]*elastic.SensorDataDocument, 0, sensorDataBufferSizeLimit)
		}

		if len(v) == sensorDataBufferSizeLimit {
			logs.LogWarn("Sensor data buffer size limit reached, overwriting data...")

			i := sensorDataBufferOverflow[doc.DeviceId]

			if i >= sensorDataBufferSizeLimit {
				i = 0
			}

			v[i] = doc
			i++
			sensorDataBufferOverflow[doc.DeviceId] = i
		} else {
			v = append(v, doc)
		}

		sensorDataBuffer[doc.DeviceId] = v
	}
}

func FlushBufferToElastic(deviceId string) {
	data, ok := sensorDataBuffer[deviceId]
	if !ok || len(data) == 0 {
		return
	}

	elastic.SendSensorDataToElastic(data)
	sensorDataBuffer[deviceId] = make([]*elastic.SensorDataDocument, 0)
	sensorDataBufferOverflow[deviceId] = 0
}

func GetBufferLength(deviceId string) int {
	v, ok := sensorDataBuffer[deviceId]
	if !ok {
		return 0
	}

	return len(v)
}
