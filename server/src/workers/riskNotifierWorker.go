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
	riskContactMinDuration time.Duration
}

func NewRiskNotifierWorker(repo interfaces.NotificationRepository, cacheRepository interfaces.CacheRepository, minContactDuration time.Duration) *RiskNotifierWorker {
	worker := new(RiskNotifierWorker)
	worker.notificationRepo = repo
	worker.cacheRepository = cacheRepository
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

		// TODO: Notify via mqtt topic

		// Push report back to 'queue' if some error ocurred for new attempt
		time.AfterFunc(tryAgainAfter, func() {
			w.pushNotificationBack(notifications, notificationJob)
		})

		// If everything went well, save notification to db
		now := time.Now()
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
