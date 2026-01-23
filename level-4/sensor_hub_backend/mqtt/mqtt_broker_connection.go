package mqtt

import (
	"fmt"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

func onConnectionUp(cm *autopaho.ConnectionManager, _ *paho.Connack) {
	fmt.Println("Connected to MQTT broker")

	err := initSubscriptions(cm)
	if err != nil {
		fmt.Printf("Failed to setup subscriptions... %s\n", err)
	}
}

func onConnectionError(err error) {
	fmt.Printf("error while attempting to connect to mqtt broker: %s\n", err)
}

func onConnectionDown() bool {
	fmt.Println("Disconnected from MQTT broker, reconnecting...")
	return true // true: Library will attempt to reconnect
}
