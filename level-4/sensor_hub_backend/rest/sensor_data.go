package rest

import (
	"io"
	"sensor_hub_backend/lifecycle"
	"sensor_hub_backend/logs"
	"sensor_hub_backend/mqtt/sensor/sensor_data"
	"sensor_hub_backend/proto_types"
	"sensor_hub_backend/rest/renderer"
	"sensor_hub_backend/sql"

	"github.com/gin-gonic/gin"
)

func RegisterSensorRoutes(router gin.IRouter) {
	router.GET("/live", getSensorDataStream)
}

func getSensorDataStream(c *gin.Context) {
	stopContext := lifecycle.GetStopContext()

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	sensorDataChannel := make(chan *proto_types.SensorData)
	i := sensor_data.Observable.Subscribe(sensorDataChannel)
	defer sensor_data.Observable.Unsubscribe(i)

	c.Stream(func(w io.Writer) bool {
		select {
		case data := <-sensorDataChannel:
			measurements := make([]*renderer.SensorMeasurementTemplateDto, len(data.Measurements))

			for i, measurement := range data.GetMeasurements() {
				measurements[i] = &renderer.SensorMeasurementTemplateDto{
					SensorType:      measurement.GetSensorType(),
					SensorValueType: measurement.GetSensorValueType().String(),
					SensorValueName: measurement.GetSensorValueName(),
					ReadTimestamp:   measurement.GetReadTimestamp(),
					Value:           sensor_data.ReadMeasurementValueHumanized(measurement),
				}
			}

			buffered := ""
			authorized, err := sql.IsDeviceAuthorizedCached(data.GetDeviceId())
			if err == nil {
				if authorized {
					buffered = "❌"
				} else {
					buffered = "✅"
				}
			} else {
				logs.LogErr("Failed to check device authorization", err)
			}

			dto := &renderer.SensorDataTemplateDto{
				Id:           data.GetDeviceId(),
				Name:         data.GetDeviceName(),
				Ip:           data.GetIp(),
				Buffered:     buffered,
				Measurements: measurements,
			}

			htmlString, err := renderer.RenderSensorDataHtml(dto)
			if err != nil {
				logs.LogErr("Failed to render sensor data", err)
				return false
			}

			c.SSEvent("sensor-data", htmlString)
			return true
		case <-c.Request.Context().Done():
			return false
		case <-stopContext.Done():
			return false
		}
	})
}
