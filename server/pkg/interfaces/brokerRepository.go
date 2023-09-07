package interfaces

import (
	"contacttracing/pkg/models/dto"
)

type ProcessContactsHandler func(dto.ContactMessage)

type BrokerRepository interface {
	PublishNotification(user string, notification dto.NotificationMessage) error
	SubscribeToReceiveContacts(topic string, processContacts ProcessContactsHandler)
}
