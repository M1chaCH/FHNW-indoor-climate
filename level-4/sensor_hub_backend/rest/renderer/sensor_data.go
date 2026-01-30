package renderer

import (
	"bytes"
	"text/template"
)

var sensorDataTemplate *template.Template

type SensorDataTemplateDto struct {
	Id           string
	Name         string
	Ip           string
	Buffered     string
	Measurements []*SensorMeasurementTemplateDto
}

type SensorMeasurementTemplateDto struct {
	SensorType      string
	SensorValueType string
	SensorValueName string
	ReadTimestamp   string
	Value           string
}

func RenderSensorDataHtml(data *SensorDataTemplateDto) (string, error) {
	if sensorDataTemplate == nil {
		var file = "rest/templates/sensor_data.html"
		var err error
		sensorDataTemplate, err = template.New("sensor_data.html").ParseFiles(file)

		if err != nil {
			return "", err
		}
	}

	buf := new(bytes.Buffer)
	err := sensorDataTemplate.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
