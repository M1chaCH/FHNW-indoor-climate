package renderer

import (
	"sensor_hub_backend/rest/templates"
	"sensor_hub_backend/sql"
)

type DeviceListDto struct {
	Devices []*DeviceDto `json:"devices"`
}

type DeviceDto struct {
	BufferSize int              `json:"buffer_size"`
	Device     sql.DeviceEntity `json:"device"`
}

func RenderDeviceHtml(data *DeviceListDto) (string, error) {
	return templates.RenderTemplate("device.html", data)
}
