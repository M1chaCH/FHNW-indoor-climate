package device_config

import (
	"sensor_hub_backend/config"
	"sensor_hub_backend/logs"
	"sensor_hub_backend/proto_types"
	"sensor_hub_backend/sql"

	"github.com/eclipse/paho.golang/paho"
	"google.golang.org/protobuf/proto"
)

func HandleDeviceConfigReceived(p *paho.Publish) {
	protoDeviceConfig := &proto_types.DeviceConfigOptions{}

	if err := proto.Unmarshal(p.Payload, protoDeviceConfig); err != nil {
		logs.LogErr("Failed to unmarshal received device config", err)
		return
	}

	deviceConfig := config.ParseProtoConfig(protoDeviceConfig)

	jsonOptions, err := config.DeviceConfigOptionsToJsonString(deviceConfig.Options)

	if err != nil {
		logs.LogErr("Failed to marshal device config options", err)
		return
	}

	deviceConfigEntity := sql.DeviceConfigEntity{
		DeviceId:   deviceConfig.DeviceId,
		ConfigJson: jsonOptions,
	}

	sql.UpsertConfigJson(&deviceConfigEntity)
}
