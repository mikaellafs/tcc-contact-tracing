package repositories

import (
	"encoding/json"
	"log"

	"contacttracing/pkg/interfaces"
	"contacttracing/pkg/models/dto"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	mqttRepositoryLog = "[MQTT Repository]"
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
	repo.userBaseTopic = "notificacao/"

	return repo
}

func (r *MqttRepository) PublishNotification(user string, notification dto.NotificationMessage) error {
	notificationJSON, _ := json.Marshal(notification)
	notificationStr := string(notificationJSON)
	log.Println(mqttRepositoryLog, "Publish notification to user", user, "-", notificationStr)

	topic := r.userBaseTopic + user
	token := r.client.Publish(topic, r.qos, true, notificationStr)

	if token.Wait() && token.Error() != nil {
		log.Println(mqttRepositoryLog, token.Error().Error())
		return token.Error()
	}

	log.Println(mqttRepositoryLog, "Notification sent successfully to topic:", topic)
	return nil
}

func (r *MqttRepository) SubscribeToReceiveContacts(topic string, processContacts interfaces.ProcessContactsHandler) {
	log.Println(mqttRepositoryLog, "Subscribe to topic", topic, "and receive contacts")
	handler := func(c mqtt.Client, m mqtt.Message) {
		var contactMessage dto.ContactMessage
		err := json.Unmarshal(m.Payload(), &contactMessage)
		if err != nil {
			log.Println(mqttRepositoryLog, "Failed to unmarshal message from contacts topic:", err.Error())
			return
		}

		processContacts(contactMessage)
	}

	token := r.client.Subscribe(topic, r.qos, handler)
	if token.Wait() && token.Error() != nil {

	}
}
