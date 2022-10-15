package workers

import (
	"context"
	"log"
	"time"

	"contacttracing/src/interfaces"
	"contacttracing/src/models/db"
	"contacttracing/src/models/dto"
)

const (
	riskNotifierWorkerLog      = "[Risk Notifier Worker]"
	maxAttemptsNotificationJob = 5
)

type RiskNotifierWorker struct {
	cacheRepository        interfaces.CacheRepository
	notificationRepo       interfaces.NotificationRepository
	brokerRepository       interfaces.BrokerRepository
	riskContactMinDuration time.Duration
}

func NewRiskNotifierWorker(repo interfaces.NotificationRepository,
	cacheRepository interfaces.CacheRepository,
	brokerRepository interfaces.BrokerRepository,
	minContactDuration time.Duration) *RiskNotifierWorker {
	worker := new(RiskNotifierWorker)
	worker.notificationRepo = repo
	worker.cacheRepository = cacheRepository
	worker.brokerRepository = brokerRepository
	worker.riskContactMinDuration = minContactDuration

	return worker
}

func (w *RiskNotifierWorker) Work(notifications chan dto.NotificationJob) {
	log.Println(riskNotifierWorkerLog, "Start work")
	for {
		// Wait for notification
		notificationJob := <-notifications

		// Check contact duration: discard if it was too short
		if notificationJob.Duration < w.riskContactMinDuration {
			log.Println(riskNotifierWorkerLog, "Contact last less than", w.riskContactMinDuration, "minutes. User is not going to get notified")
			continue
		}

		// Check if user has been notified
		cacheNotification := w.cacheRepository.GetNotificationFrom(notificationJob.ForUser, notificationJob.FromReport)
		if cacheNotification != nil {
			log.Println(riskNotifierWorkerLog, "User already has been notified about the contact with this user for this report:", cacheNotification)
			continue
		}

		// Notify via mqtt topic
		now := time.Now()
		err := w.brokerRepository.PublishNotification(notificationJob.ForUser, w.makeUserNotificationMessage(notificationJob, now))

		// Push report back to 'queue' if some error ocurred for new attempt
		if err != nil {
			time.AfterFunc(tryAgainAfter, func() {
				w.pushNotificationBack(notifications, notificationJob)
			})
			continue
		}

		// If everything went well, save notification to db
		w.notificationRepo.Create(context.TODO(), db.Notification{
			ForUser:    notificationJob.ForUser,
			FromReport: notificationJob.FromReport,
			Date:       now,
		})

		// And also in cache
		w.cacheRepository.SaveNotification(notificationJob.ForUser, notificationJob.FromReport, now)
	}
}

func (w *RiskNotifierWorker) pushNotificationBack(notifications chan<- dto.NotificationJob, job dto.NotificationJob) {
	// Check attempts
	if job.Attempts >= maxAttemptsNotificationJob {
		log.Println(contactTracerWorkerLog, "Cannot try to push notification job back to queue, attempts ran out")
		return
	}

	// New attempt
	log.Println(riskNotifierWorkerLog, "Push notification job back to channel for new attempt")
	job.Attempts += 1
	notifications <- job
}

func (w *RiskNotifierWorker) makeUserNotificationMessage(notification dto.NotificationJob, date time.Time) dto.NotificationMessage {
	return dto.NotificationMessage{
		Risk: true,
		Message: `Você esteve em contato com uma pessoa diagnosticada com covid-19 nos últimos 15 dias. 
				  Siga as recomendações de saúde. Notificado(a) em ` + date.Format(time.RFC822),
		Date:         date,
		AmountPeople: 1,
	}
}

func AddNotificationJob(userNotified string, fromReport int64, contactDuration time.Duration, notifChan chan<- dto.NotificationJob) {
	notificationJob := dto.NotificationJob{
		ForUser:    userNotified,
		FromReport: fromReport,
		Duration:   contactDuration,
		Attempts:   0,
	}

	log.Println(riskNotifierWorkerLog, "Add notification job:", notificationJob)
	notifChan <- notificationJob
}
