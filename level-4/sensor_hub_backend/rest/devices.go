package rest

import (
	"fmt"
	"io"
	"sensor_hub_backend/elastic/buffer"
	"sensor_hub_backend/lifecycle"
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
			fmt.Printf("Subscription to devices failed due to an error: %s\n", err)
		} else {
			fmt.Println("Subscription to devices closed")
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
				fmt.Printf("Failed to render device list: %s\n", err)
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
