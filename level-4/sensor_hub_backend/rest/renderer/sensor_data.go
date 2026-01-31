package renderer

import (
	"sensor_hub_backend/rest/templates"
)

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
	return templates.RenderTemplate("sensor_data.html", data)
}
