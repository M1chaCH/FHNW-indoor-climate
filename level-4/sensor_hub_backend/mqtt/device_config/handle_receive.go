package device_config

import (
	"fmt"
	"sensor_hub_backend/config"
	"sensor_hub_backend/proto_types"
	"sensor_hub_backend/sql"

	"github.com/eclipse/paho.golang/paho"
	"google.golang.org/protobuf/proto"
)

func HandleDeviceConfigReceived(p *paho.Publish) {
	protoDeviceConfig := &proto_types.DeviceConfigOptions{}

	if err := proto.Unmarshal(p.Payload, protoDeviceConfig); err != nil {
		fmt.Printf("Failed to unmarshal received device config: %s", err)
		return
	}

	deviceConfig := config.ParseProtoConfig(protoDeviceConfig)

	jsonOptions, err := config.DeviceConfigOptionsToJsonString(deviceConfig.Options)

	if err != nil {
		fmt.Printf("Failed to marshal device config options: %s", err)
		return
	}

	deviceConfigEntity := sql.DeviceConfigEntity{
		DeviceId:   deviceConfig.DeviceId,
		ConfigJson: jsonOptions,
	}

	sql.UpsertConfigJson(&deviceConfigEntity)
}
