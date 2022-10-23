package mqtt

import (
	"log"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
)

func DefaultPublishHandler(client pahomqtt.Client, msg pahomqtt.Message) {
	log.Printf("[Default Publish Handler]: Message received. \nTOPIC: %s\nMESSAGE: %s", msg.Topic(), msg.Payload())
}
