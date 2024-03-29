package workers

import (
	"context"
	"log"
	"strconv"
	"time"

	"contacttracing/pkg/interfaces"
	"contacttracing/pkg/models/db"
	"contacttracing/pkg/models/dto"
	"contacttracing/pkg/repositories"
	"contacttracing/pkg/utils"
)

const (
	riskNotifierWorkerLog      = "[Risk Notifier Worker]"
	maxAttemptsNotificationJob = 5
)

type RiskNotifierWorker struct {
	cacheRepository  interfaces.CacheRepository
	notificationRepo interfaces.NotificationRepository
	brokerRepository interfaces.BrokerRepository
	daysTraced       int
}

func NewRiskNotifierWorker(repo interfaces.NotificationRepository,
	cacheRepository interfaces.CacheRepository,
	brokerRepository interfaces.BrokerRepository,
	daysTraced int) *RiskNotifierWorker {

	worker := new(RiskNotifierWorker)
	worker.notificationRepo = repo
	worker.cacheRepository = cacheRepository
	worker.brokerRepository = brokerRepository
	worker.daysTraced = daysTraced

	return worker
}

func (w *RiskNotifierWorker) Work(notifications chan dto.NotificationJob) {
	log.Println(riskNotifierWorkerLog, "Start work")
	for {
		// Wait for notification
		notificationJob := <-notifications

		// Check if user has been notified
		cacheNotification := w.cacheRepository.GetNotificationFrom(notificationJob.ForUser, notificationJob.FromReport)
		if cacheNotification != nil {
			log.Println(riskNotifierWorkerLog, "User already has been notified about the contact with this user for this report:", cacheNotification)
			continue
		}

		// Notify via mqtt topic if user still at risk
		timeDiffStillAtRisk := time.Duration(w.daysTraced) * time.Hour * 24
		isAtRisk := utils.VerifyUserAtRisk(notificationJob.DateLastContact, time.Now(), 0, timeDiffStillAtRisk)

		var err error
		now := time.Now()

		if !isAtRisk {
			log.Println(riskNotifierWorkerLog, "User is not at risk anymore, don't need to notify them right now")
		} else {
			// Save in cache
			w.cacheRepository.SaveNotification(notificationJob.ForUser, notificationJob.FromReport, now)
			err = w.brokerRepository.PublishNotification(notificationJob.ForUser, w.makeUserNotificationMessage(notificationJob, now))
		}

		// Push report back to 'queue' if some error ocurred for new attempt
		if err != nil {
			w.scheduleToPushBack(notifications, notificationJob)
			continue
		}

		// If everything went well, save notification to db
		// Even if user has not been notified, it's important to save that a risk in the past was identied
		_, err = w.notificationRepo.Create(context.TODO(), db.Notification{
			ForUser:    notificationJob.ForUser,
			FromReport: notificationJob.FromReport,
			Date:       now,
		})
		if err != nil && err != repositories.ErrDuplicate {
			w.scheduleToPushBack(notifications, notificationJob)
			continue
		}
	}
}

func (w *RiskNotifierWorker) scheduleToPushBack(notifications chan<- dto.NotificationJob, job dto.NotificationJob) {
	time.AfterFunc(tryAgainAfter, func() {
		w.pushNotificationBack(notifications, job)
	})
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
		Message: "Você esteve em contato com uma pessoa diagnosticada com covid-19 nos últimos " +
			strconv.Itoa(w.daysTraced) + " dias." +
			" Siga as recomendações de saúde. Notificado(a) em " + date.In(time.FixedZone("UTC-3", -3*60*60)).Format(time.RFC822),
		Date:         date,
		AmountPeople: 1,
	}
}

func AddNotificationJob(lastContact time.Time, userNotified string, fromReport int64, contactDuration time.Duration, notifChan chan<- dto.NotificationJob) {
	notificationJob := dto.NotificationJob{
		DateLastContact: lastContact,
		ForUser:         userNotified,
		FromReport:      fromReport,
		Duration:        contactDuration,
		Attempts:        0,
	}

	log.Println(riskNotifierWorkerLog, "Add notification job:", notificationJob)
	notifChan <- notificationJob
}
