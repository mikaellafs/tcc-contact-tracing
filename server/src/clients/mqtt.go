package clients

import (
	"os"
	"time"

	"contacttracing/src/utils/mqtt"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	mqttKeepAlive   = 60 * time.Second
	mqttPingTimeout = 1 * time.Second
)

func NewMqttClient() pahomqtt.Client {
	host := os.Getenv("MQTT_BROKER_HOST")
	port := os.Getenv("MQTT_BROKER_PORT")

	clientId := os.Getenv("MQTT_BROKER_CLIENTID")
	user := os.Getenv("MQTT_BROKER_USERNAME")
	password := os.Getenv("MQTT_BROKER_PASSWORD")

	opts := pahomqtt.NewClientOptions().AddBroker(host + ":" + port).SetClientID("emqx_test_client")

	opts.SetKeepAlive(mqttKeepAlive)

	opts.SetClientID(clientId)
	opts.SetUsername(user)
	opts.SetPassword(password)

	opts.SetDefaultPublishHandler(mqtt.DefaultPublishHandler)
	opts.SetPingTimeout(mqttPingTimeout)

	client := pahomqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}
