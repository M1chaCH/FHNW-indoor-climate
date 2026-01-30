package renderer

import (
	"bytes"
	"text/template"
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

var deviceConfigTemplate *template.Template

func RenderDeviceConfigHtml(data *DeviceConfigRenderingDto) (string, error) {
	if deviceConfigTemplate == nil {
		var file = "rest/templates/device_config.html"
		var err error
		deviceConfigTemplate, err = template.New("device_config.html").ParseFiles(file)

		if err != nil {
			return "", err
		}
	}

	buf := new(bytes.Buffer)
	err := deviceConfigTemplate.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
