package rest

import (
	"io"
	"sensor_hub_backend/elastic/buffer"
	"sensor_hub_backend/lifecycle"
	"sensor_hub_backend/logs"
	"sensor_hub_backend/rest/renderer"
	"sensor_hub_backend/sql"

	"github.com/gin-gonic/gin"
)

func RegisterDevicesRoutes(router gin.IRouter) {
	router.GET("", getDevicesStream)
	router.POST("/authorize/:device_id", postToggleAuthorizeDevice)
}

func getDevicesStream(c *gin.Context) {
	stopContext := lifecycle.GetStopContext()

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	deviceChan := make(chan []sql.DeviceEntity)

	go func() {
		err := sql.SubscribeToDevices(deviceChan, c.Request.Context())
		if err != nil {
			logs.LogErr("Subscription to devices failed due to an error", err)
		} else {
			logs.LogInfo("Subscription to devices closed")
		}
	}()

	c.Stream(func(w io.Writer) bool {
		select {
		case devices := <-deviceChan:
			dtos := make([]*renderer.DeviceDto, len(devices))
			for i, device := range devices {
				dtos[i] = &renderer.DeviceDto{
					Device:     device,
					BufferSize: buffer.GetBufferLength(device.DeviceId),
				}
			}

			htmlString, err := renderer.RenderDeviceHtml(&renderer.DeviceListDto{Devices: dtos})
			if err != nil {
				logs.LogErr("Failed to render device list", err)
				return false
			}

			c.SSEvent("devices-changed", htmlString)
			return true
		case <-c.Request.Context().Done():
			return false
		case <-stopContext.Done():
			return false
		}
	})
}

func postToggleAuthorizeDevice(c *gin.Context) {
	deviceId := c.Param("device_id")
	active, err := sql.ToggleDeviceAuthorization(deviceId)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if active {
		buffer.FlushBufferToElastic(deviceId)
	}

	c.JSON(200, gin.H{"message": "toggled"})
}
