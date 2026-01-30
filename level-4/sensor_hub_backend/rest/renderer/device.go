package renderer

import (
	"bytes"
	"sensor_hub_backend/sql"
	"text/template"
)

type DeviceListDto struct {
	BufferSize int              `json:"buffer_size"`
	Device     sql.DeviceEntity `json:"device"`
}

var deviceTemplate *template.Template

func RenderDeviceHtml(data []*DeviceListDto) (string, error) {
	if deviceTemplate == nil {
		var file = "rest/templates/device.html"
		var err error
		deviceTemplate, err = template.New("device.html").ParseFiles(file)

		if err != nil {
			return "", err
		}
	}

	buf := new(bytes.Buffer)
	err := deviceTemplate.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
