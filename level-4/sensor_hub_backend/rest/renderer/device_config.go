package renderer

import (
	"sensor_hub_backend/rest/templates"
)

type DeviceConfigOptionRenderingDto struct {
	Name  string
	Value string
	Type  int
}

type DeviceConfigRenderingDto struct {
	DeviceId      string
	IntOptions    []DeviceConfigOptionRenderingDto
	StringOptions []DeviceConfigOptionRenderingDto
	DoubleOptions []DeviceConfigOptionRenderingDto
	FlagOptions   []DeviceConfigOptionRenderingDto
}

func RenderDeviceConfigHtml(data *DeviceConfigRenderingDto) (string, error) {
	return templates.RenderTemplate("device_config.html", data)
}
