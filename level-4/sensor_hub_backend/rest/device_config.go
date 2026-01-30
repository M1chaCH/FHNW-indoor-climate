package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sensor_hub_backend/config"
	"sensor_hub_backend/mqtt"
	"sensor_hub_backend/proto_types"
	"sensor_hub_backend/rest/renderer"
	"sensor_hub_backend/sql"

	"github.com/gin-gonic/gin"
)

func RegisterDeviceConfigRoutes(router gin.IRouter) {
	router.GET("/:device_id", getConfigOfDevice)
	router.POST("/:device_id", postConfigOfDevice)
	router.POST("/:device_id/push", postPushConfigToDevice)
}

func getConfigOfDevice(c *gin.Context) {
	deviceId := c.Param("device_id")

	optionsJson, err := sql.SelectDeviceConfigJson(deviceId)
	if err != nil {
		fmt.Printf("Failed to get device config: %s\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	options, err := config.DeviceConfigOptionsFromJsonString(optionsJson)
	if err != nil {
		fmt.Printf("Failed to unmarshal device config options: %s\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	renderingDto := createRenderingDto(deviceId, options)
	htmlString, err := renderer.RenderDeviceConfigHtml(&renderingDto)
	if err != nil {
		fmt.Printf("Failed to render device config: %s\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, "text/html", []byte(htmlString))
}

func createRenderingDto(deviceId string, options []config.DeviceConfigOption) renderer.DeviceConfigRenderingDto {
	intOptions := make([]renderer.DeviceConfigOptionRenderingDto, 0)
	stringOptions := make([]renderer.DeviceConfigOptionRenderingDto, 0)
	doubleOptions := make([]renderer.DeviceConfigOptionRenderingDto, 0)
	flagOptions := make([]renderer.DeviceConfigOptionRenderingDto, 0)

	for _, option := range options {
		switch option.Type {
		case int(proto_types.DeviceConfigOptions_INT32):
			intOptions = append(intOptions, renderer.DeviceConfigOptionRenderingDto{
				Name:  option.Name,
				Value: fmt.Sprintf("%v", option.Value),
				Type:  int(proto_types.DeviceConfigOptions_INT32),
			})
		case int(proto_types.DeviceConfigOptions_DOUBLE):
			doubleOptions = append(doubleOptions, renderer.DeviceConfigOptionRenderingDto{
				Name:  option.Name,
				Value: fmt.Sprintf("%v", option.Value),
				Type:  int(proto_types.DeviceConfigOptions_DOUBLE),
			})
		case int(proto_types.DeviceConfigOptions_STRING):
			stringOptions = append(stringOptions, renderer.DeviceConfigOptionRenderingDto{
				Name:  option.Name,
				Value: fmt.Sprintf("%v", option.Value),
				Type:  int(proto_types.DeviceConfigOptions_STRING),
			})
		case int(proto_types.DeviceConfigOptions_BOOL):
			flagOptions = append(flagOptions, renderer.DeviceConfigOptionRenderingDto{
				Name:  option.Name,
				Value: fmt.Sprintf("%v", option.Value),
				Type:  int(proto_types.DeviceConfigOptions_BOOL),
			})
		}
	}

	return renderer.DeviceConfigRenderingDto{
		DeviceId:      deviceId,
		IntOptions:    intOptions,
		StringOptions: stringOptions,
		DoubleOptions: doubleOptions,
		FlagOptions:   flagOptions,
	}
}

type saveDeviceConfigDto struct {
	Value string `json:"value"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

func postConfigOfDevice(c *gin.Context) {
	deviceId := c.Param("device_id")

	formData := c.PostForm("data")

	var options []saveDeviceConfigDto
	if err := json.Unmarshal([]byte(formData), &options); err != nil {
		fmt.Printf("Failed to parse request body: %s\n", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	deviceConfigOptions := make([]config.DeviceConfigOption, len(options))
	for i, data := range options {
		option, err := config.CreateDeviceConfigOption(data.Name, data.Type, data.Value)
		if err != nil {
			fmt.Printf("Failed to create device config option: %s\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		deviceConfigOptions[i] = option
	}

	jsonString, err := config.DeviceConfigOptionsToJsonString(deviceConfigOptions)
	if err != nil {
		fmt.Printf("Failed to marshal device config options: %s\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	sql.UpsertConfigJson(&sql.DeviceConfigEntity{
		DeviceId:   deviceId,
		ConfigJson: jsonString,
	})

	protoConfig := config.CreateProtoConfig(&config.DeviceConfig{
		DeviceId: deviceId,
		Options:  deviceConfigOptions,
	})
	go mqtt.PushConfigToDevice(protoConfig)

	c.Data(http.StatusOK, "text/html", []byte("<p>Successfully updated device config</p>"))
}

func postPushConfigToDevice(c *gin.Context) {
	deviceId := c.Param("device_id")

	optionJson, err := sql.SelectDeviceConfigJson(deviceId)
	if err != nil {
		fmt.Printf("Failed to get device config: %s\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	deviceConfigOptions, err := config.DeviceConfigOptionsFromJsonString(optionJson)
	if err != nil {
		fmt.Printf("Failed to unmarshal device config options: %s\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	protoConfig := config.CreateProtoConfig(&config.DeviceConfig{
		DeviceId: deviceId,
		Options:  deviceConfigOptions,
	})
	go mqtt.PushConfigToDevice(protoConfig)
}
