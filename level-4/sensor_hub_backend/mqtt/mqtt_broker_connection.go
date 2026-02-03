package mqtt

import (
	"sensor_hub_backend/logs"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

func onConnectionUp(cm *autopaho.ConnectionManager, _ *paho.Connack) {
	logs.LogInfo("Connected to MQTT broker")

	err := initSubscriptions(cm)
	if err != nil {
		logs.LogErr("Failed to setup subscriptions", err)
	}
}

func onConnectionError(err error) {
	logs.LogErr("error while attempting to connect to mqtt broker", err)
}

func onConnectionDown() bool {
	logs.LogWarn("Disconnected from MQTT broker, reconnecting...")
	return true // true: Library will attempt to reconnect
}
