package repositories

import (
	"encoding/json"
	"log"

	"contacttracing/src/interfaces"
	"contacttracing/src/models/dto"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	mqttRepositoryLog = "[MQTT Repository Log]"
)

type MqttRepository struct {
	client        mqtt.Client
	qos           byte
	userBaseTopic string
}

func NewMqttRepository(client mqtt.Client) *MqttRepository {
	repo := new(MqttRepository)
	repo.client = client

	repo.qos = 1
	repo.userBaseTopic = "user/"

	return repo
}

func (r *MqttRepository) PublishNotification(user string, notification dto.NotificationMessage) error {
	notificationJSON, _ := json.Marshal(notification)
	notificationStr := string(notificationJSON)
	log.Println(mqttRepositoryLog, "Publish notification to user", user, "-", notificationStr)

	token := r.client.Publish(r.userBaseTopic+user, r.qos, false, notificationStr)

	if token.Wait() && token.Error() != nil {
		log.Println(mqttRepositoryLog, token.Error().Error())
		return token.Error()
	}

	log.Println(mqttRepositoryLog, "Notification sent successfully")
	return nil
}

func (r *MqttRepository) SubscribeToReceiveContacts(topic string, processContacts interfaces.ProcessContactsHandler) {
	log.Println(mqttRepositoryLog, "Subscribe to topic", topic, "and receive contacts")
	handler := func(c mqtt.Client, m mqtt.Message) {
		var contactMessage dto.ContactMessage
		json.Unmarshal(m.Payload(), &contactMessage)

		processContacts(contactMessage)
	}

	token := r.client.Subscribe(topic, r.qos, handler)
	if token.Wait() && token.Error() != nil {

	}
}
