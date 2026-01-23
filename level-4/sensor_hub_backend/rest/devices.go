package rest

import (
	"sensor_hub_backend/elastic/buffer"
	"sensor_hub_backend/sql"

	"github.com/gin-gonic/gin"
)

func RegisterAuthorizedDevicesRoutes(router gin.IRouter) {
	router.GET("", getDevices)
	router.POST("/authorize/:device_id", postToggleAuthorizeDevice)
}

type deviceListDto struct {
	Device     sql.DeviceEntity `json:"device"`
	BufferSize int              `json:"buffer_size"`
}

func getDevices(c *gin.Context) {
	devices, err := sql.SelectDevices()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	dtos := make([]deviceListDto, len(devices))
	for i, device := range devices {
		dtos[i] = deviceListDto{
			Device:     device,
			BufferSize: buffer.GetBufferLength(device.DeviceId),
		}
	}

	c.JSON(200, dtos)
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
