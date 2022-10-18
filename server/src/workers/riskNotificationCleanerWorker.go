package workers

import (
	"contacttracing/src/interfaces"
	"contacttracing/src/models/dto"
	"log"
	"time"
)

const (
	riskNotificationCleanerWorkerLog = "[Risk Notification Cleaner]"
	scheduleCleanNotificationsTime   = 24 * time.Hour
)

type RiskNotificationCleanerWorker struct {
	cacheRepo      interfaces.CacheRepository
	brokerRepo     interfaces.BrokerRepository
	expirationDays time.Duration
}

func NewRiskNotificationCleanerWorker(
	cacheRepo interfaces.CacheRepository,
	brokerRepo interfaces.BrokerRepository,
	expirationDays time.Duration) *RiskNotificationCleanerWorker {

	worker := new(RiskNotificationCleanerWorker)
	worker.cacheRepo = cacheRepo
	worker.brokerRepo = brokerRepo
	worker.expirationDays = expirationDays

	return worker
}

func (w *RiskNotificationCleanerWorker) Work() {
	log.Println(riskNotificationCleanerWorkerLog, "Start work")
	w.scheduleClean()
}

func (w *RiskNotificationCleanerWorker) scheduleClean() {
	log.Println(riskNotificationCleanerWorkerLog, "Schedule cleaning for", time.Now().Add(scheduleCleanNotificationsTime))
	time.AfterFunc(scheduleCleanNotificationsTime, func() {
		w.clean()
		w.scheduleClean()
	})
}

func (w *RiskNotificationCleanerWorker) clean() {
	// Remove expired notifications and get users
	users := w.cacheRepo.RemoveNotificationsAfter(w.expirationDays, time.Now())
	log.Println(riskNotificationCleanerWorkerLog, len(users), "expired notifications...")

	for _, user := range users {
		// Check if there's still a notification for each user
		userNotifications := w.cacheRepo.GetNotificationKeysFrom(user)
		if len(userNotifications) > 0 {
			log.Println(riskNotificationCleanerWorkerLog, "User", user, "had contact with other people infected, still at risk")
			continue
		}

		// If there's none, notify user about not being at risk anymore
		log.Println(riskNotificationCleanerWorkerLog, "User", user, "is not at risk anymore")
		w.brokerRepo.PublishNotification(user, w.makeNotificationMessage())
	}
}

func (w *RiskNotificationCleanerWorker) makeNotificationMessage() dto.NotificationMessage {
	return dto.NotificationMessage{
		Risk: false,
	}
}
