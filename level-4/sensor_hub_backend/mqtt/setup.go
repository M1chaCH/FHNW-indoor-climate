package mqtt

import (
	"fmt"
	"net/url"
	"os"
	"sensor_hub_backend/lifecycle"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

func RunMqttClient() {
	urlString := os.Getenv("MQTT_BROKER_URL")
	user := os.Getenv("MQTT_BROKER_USER")
	password := os.Getenv("MQTT_BROKER_PASSWORD")

	u, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}

	clientConfig := autopaho.ClientConfig{
		ServerUrls:                    []*url.URL{u},
		ConnectUsername:               user,
		ConnectPassword:               []byte(password),
		KeepAlive:                     20,
		CleanStartOnInitialConnection: false,
		SessionExpiryInterval:         60,
		OnConnectionUp:                onConnectionUp,
		OnConnectError:                onConnectionError,
		OnConnectionDown:              onConnectionDown,
		ClientConfig: paho.ClientConfig{
			ClientID:          "sensor_hub_backend",
			OnPublishReceived: []func(paho.PublishReceived) (bool, error){onPublishReceived},
			OnClientError:     onClientError,
		},
	}

	stopContext := lifecycle.GetStopContext()

	c, err := autopaho.NewConnection(stopContext, clientConfig)
	if err != nil {
		panic(err)
	}

	if err = c.AwaitConnection(stopContext); err != nil {
		panic(err)
	}

	<-stopContext.Done()
	fmt.Println("Stopping MQTT service...")
	<-c.Done()
}
