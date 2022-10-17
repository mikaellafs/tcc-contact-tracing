package workers

import (
	"context"
	"errors"
	"log"
	"time"

	"contacttracing/src/interfaces"
	"contacttracing/src/models/db"
	"contacttracing/src/models/dto"
	"contacttracing/src/utils"
)

const (
	contactsProcessorLog = "[Contacts Processor]"
)

type ContactsProcessor struct {
	contactRepository           interfaces.ContactRepository
	userRepository              interfaces.UserRepository
	cacheRepository             interfaces.CacheRepository
	potentialRiskChan           chan<- dto.PotentialRiskJob
	maxTimeDiffToConsiderAtRisk time.Duration
}

func NewContactsProcessor(
	contactRepository interfaces.ContactRepository,
	userRepository interfaces.UserRepository,
	cacheRepository interfaces.CacheRepository,
	potentialRiskChan chan<- dto.PotentialRiskJob,
	maxDiffDaysFromDiagnosticToConsiderAtRisk time.Duration) *ContactsProcessor {

	processor := new(ContactsProcessor)
	processor.contactRepository = contactRepository
	processor.userRepository = userRepository
	processor.potentialRiskChan = potentialRiskChan
	processor.maxTimeDiffToConsiderAtRisk = maxDiffDaysFromDiagnosticToConsiderAtRisk * time.Hour * 24

	return processor
}

func (p *ContactsProcessor) Process(contact dto.ContactMessage) {
	log.Println(contactsProcessorLog, "Processing contact from user", contact.User, ", info:", string(contact.Contact))

	// Validate message
	err := p.validateMessage(contact)
	if err != nil {
		log.Println(contactsProcessorLog, err.Error())
		return
	}

	// Parse contact to ContactMessage
	contactFromMsg := contact.ParseContact()
	if contactFromMsg == nil {
		log.Println(contactsProcessorLog, "Could not process message: invalid contact format")
		return
	}

	// Save contact in db
	p.saveContact(contact.User, contactFromMsg)

	// Verify if user contacted has reported covid in the last 15 days
	reports := p.cacheRepository.GetReportsFrom(contactFromMsg.User)

	for _, report := range reports {
		p.checkAndProcessUserRisk(contactFromMsg, report)
	}
}

func (p *ContactsProcessor) validateMessage(contact dto.ContactMessage) error {
	// Get user PK
	user, err := p.userRepository.GetByUserId(context.TODO(), contact.User)
	if err != nil {
		return errors.New("Error when getting user, cannot get public key: " + err.Error())
	}

	// Validate message
	isValid, err := utils.ValidateMessage(string(contact.Contact), user.Pk, contact.Signature)
	if err != nil {
		return errors.New("Error when validating message: " + err.Error())
	}

	if !isValid {
		return errors.New("Message is invalid")
	}

	return nil
}

func (p *ContactsProcessor) saveContact(userId string, contact *dto.ContactFromMessage) {
	ctx := context.Background()

	_, err := p.contactRepository.Create(ctx, db.Contact{
		User:                  userId,
		AnotherUser:           contact.User,
		FirstContactTimestamp: time.UnixMilli(contact.FirstContactTimestamp),
		LastContactTimestamp:  time.UnixMilli(contact.LastContactTimestamp),
		Distance:              contact.Distance,
		RSSI:                  contact.RSSI,
		BatteryLevel:          contact.BatteryLevel,
	})

	// If some error ocurred...
	if err != nil {
		// TODO: try again later?
	}
}

func (p *ContactsProcessor) checkAndProcessUserRisk(contact *dto.ContactFromMessage, report dto.Report) {
	isAtRisk := utils.VerifyUserAtRisk(time.UnixMilli(contact.LastContactTimestamp), report.DateDiagnostic, contact.Distance, p.maxTimeDiffToConsiderAtRisk)

	if !isAtRisk {
		log.Println(contactsProcessorLog, "User", contact.User, "is NOT at risk for now")
		return
	}

	log.Println(contactsProcessorLog, "User", contact.User, "is in contact with infected user", contact.User)
	AddPotentialRiskJob(contact.User, report.ID, p.potentialRiskChan, p.cacheRepository)
}
