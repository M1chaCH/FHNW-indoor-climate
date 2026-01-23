package mqtt

import (
	"fmt"
	"sensor_hub_backend/lifecycle"
	"sensor_hub_backend/mqtt/sensor/sensor_data"
	"strings"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

const (
	topicSensorData = "sensor/data"
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
		},
	})

	return err
}

func onClientError(err error) {
	fmt.Printf("MQTT Client error occurred: %s\n", err)
}

func onPublishReceived(p paho.PublishReceived) (bool, error) {
	fmt.Printf("Received MQTT published message: %s %d bytes\n", p.Packet.Topic, len(p.Packet.Payload))

	handleReceivedMessageErrors(p.Errs)

	if strings.HasPrefix(p.Packet.Topic, topicSensorData) {
		sensor_data.HandleSensorDataReceived(p.Packet)
	} else {
		fmt.Printf("No handler found for topic: %s\n", p.Packet.Topic)
	}

	return true, nil
}

func handleReceivedMessageErrors(errs []error) {
	if len(errs) == 0 {
		return
	}

	fmt.Printf("Received MQTT published message with %d errors:\n", len(errs))
	for _, err := range errs {
		fmt.Printf("  %s\n", err)
	}
}
