package config

import (
	"encoding/json"
	"sensor_hub_backend/proto_types"
	"strconv"
)

type DeviceConfig struct {
	DeviceId string               `json:"deviceId"`
	Options  []DeviceConfigOption `json:"options"`
}

type DeviceConfigOption struct {
	Value interface{} `json:"value"`
	Type  int         `json:"type"`
	Name  string      `json:"name"`
}

func CreateProtoConfig(dc *DeviceConfig) *proto_types.DeviceConfigOptions {
	options := make([]*proto_types.DeviceConfigOptions_ConfigOption, len(dc.Options))
	for i, option := range dc.Options {
		options[i] = &proto_types.DeviceConfigOptions_ConfigOption{
			ConfigName: option.Name,
			ConfigType: proto_types.DeviceConfigOptions_ConfigOptionValueType(option.Type),
		}

		switch options[i].ConfigType {
		case proto_types.DeviceConfigOptions_STRING:
			v := option.Value.(string)
			options[i].StringValue = &v
		case proto_types.DeviceConfigOptions_DOUBLE:
			v := option.Value.(float64)
			options[i].DoubleValue = &v
		case proto_types.DeviceConfigOptions_INT32:
			v, ok := option.Value.(int)
			if !ok {
				v = int(option.Value.(float64))
			}
			intValue := int32(v)
			options[i].IntValue = &intValue
		case proto_types.DeviceConfigOptions_BOOL:
			v := option.Value.(bool)
			options[i].FlagValue = &v
		}
	}

	return &proto_types.DeviceConfigOptions{
		DeviceId: dc.DeviceId,
		Options:  options,
	}
}

func ParseProtoConfig(protoConfig *proto_types.DeviceConfigOptions) DeviceConfig {
	parsedOptions := make([]DeviceConfigOption, len(protoConfig.GetOptions()))

	for i, option := range protoConfig.GetOptions() {
		parsedOptions[i] = DeviceConfigOption{
			Value: extractValue(option),
			Type:  int(option.GetConfigType()),
			Name:  option.GetConfigName(),
		}
	}

	return DeviceConfig{
		DeviceId: protoConfig.GetDeviceId(),
		Options:  parsedOptions,
	}
}

func extractValue(option *proto_types.DeviceConfigOptions_ConfigOption) interface{} {
	switch option.GetConfigType() {
	case proto_types.DeviceConfigOptions_STRING:
		return option.GetStringValue()
	case proto_types.DeviceConfigOptions_BOOL:
		return option.GetFlagValue()
	case proto_types.DeviceConfigOptions_INT32:
		return option.GetIntValue()
	case proto_types.DeviceConfigOptions_DOUBLE:
		return option.GetDoubleValue()
	default:
		return nil
	}
}

func DeviceConfigOptionsToJsonString(options []DeviceConfigOption) (string, error) {
	bytes, err := json.Marshal(options)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func DeviceConfigOptionsFromJsonString(jsonString string) ([]DeviceConfigOption, error) {
	var options []DeviceConfigOption
	err := json.Unmarshal([]byte(jsonString), &options)
	return options, err
}

func CreateDeviceConfigOption(name string, valueType string, value string) (DeviceConfigOption, error) {
	parsedValueType, err := strconv.Atoi(valueType)
	if err != nil {
		return DeviceConfigOption{}, err
	}

	var parsedValue interface{}
	switch parsedValueType {
	case int(proto_types.DeviceConfigOptions_INT32):
		parsedValue, err = strconv.Atoi(value)
	case int(proto_types.DeviceConfigOptions_DOUBLE):
		parsedValue, err = strconv.ParseFloat(value, 64)
	case int(proto_types.DeviceConfigOptions_BOOL):
		parsedValue, err = strconv.ParseBool(value)
	case int(proto_types.DeviceConfigOptions_STRING):
		parsedValue = value
	}

	if err != nil {
		return DeviceConfigOption{}, err
	}

	return DeviceConfigOption{
		Value: parsedValue,
		Type:  parsedValueType,
		Name:  name,
	}, nil
}
