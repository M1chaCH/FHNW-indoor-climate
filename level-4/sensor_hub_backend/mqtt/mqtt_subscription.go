package mqtt

import (
	"sensor_hub_backend/lifecycle"
	"sensor_hub_backend/logs"
	"sensor_hub_backend/mqtt/device_config"
	"sensor_hub_backend/mqtt/sensor/sensor_data"
	"strings"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

const (
	topicSensorData         = "sensor/data"
	topicDeviceConfig       = "device/config"
	topicUpdateDeviceConfig = "device/config/update"
)

func initSubscriptions(cm *autopaho.ConnectionManager) error {
	stopContext := lifecycle.GetStopContext()

	_, err := cm.Subscribe(stopContext, &paho.Subscribe{
		Subscriptions: []paho.SubscribeOptions{
			{
				Topic:             topicSensorData + "/+",
				QoS:               0,
				RetainAsPublished: false,
			},
			{
				Topic:             topicDeviceConfig + "/+",
				QoS:               0,
				RetainAsPublished: false,
			},
		},
	})

	return err
}

func onClientError(err error) {
	logs.LogErr("MQTT Client error occurred", err)
}

func onPublishReceived(p paho.PublishReceived) (bool, error) {
	logs.LogInfo("Received MQTT published message: %s %d bytes\n", p.Packet.Topic, len(p.Packet.Payload))

	handleReceivedMessageErrors(p.Errs)

	if strings.HasPrefix(p.Packet.Topic, topicSensorData) {
		sensor_data.HandleSensorDataReceived(p.Packet)
	} else if strings.HasPrefix(p.Packet.Topic, topicDeviceConfig) {
		device_config.HandleDeviceConfigReceived(p.Packet)
	} else {
		logs.LogWarn("No handler found for topic: %s\n", p.Packet.Topic)
	}

	return true, nil
}

func handleReceivedMessageErrors(errs []error) {
	if len(errs) == 0 {
		return
	}

	logs.LogErrCustom("Received MQTT published message with %d errors:\n", len(errs))
	for _, err := range errs {
		logs.LogInfo(err.Error())
	}
}
