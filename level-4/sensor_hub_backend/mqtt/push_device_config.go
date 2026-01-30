package mqtt

import (
	"fmt"
	"sensor_hub_backend/lifecycle"
	"sensor_hub_backend/proto_types"

	"github.com/eclipse/paho.golang/paho"
	"google.golang.org/protobuf/proto"
)

func PushConfigToDevice(data *proto_types.DeviceConfigOptions) {
	stopContext := lifecycle.GetStopContext()

	encodedPayload, err := proto.Marshal(data)
	if err != nil {
		fmt.Printf("Failed to marshal device config: %s\n", err)
		return
	}

	_, err = connection.Publish(stopContext, &paho.Publish{
		QoS:     1,
		Topic:   topicUpdateDeviceConfig + "/" + data.DeviceId,
		Payload: encodedPayload,
	})

	if err != nil {
		fmt.Printf("Failed to publish device config: %s\n", err)
		return
	}
	fmt.Printf("Successfully published device config: %s\n", topicUpdateDeviceConfig+"/"+data.DeviceId)
}
