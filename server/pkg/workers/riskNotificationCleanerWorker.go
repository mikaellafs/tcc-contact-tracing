package workers

import (
	"contacttracing/pkg/interfaces"
	"contacttracing/pkg/models/dto"
	"log"
	"time"
)

const (
	riskNotificationCleanerWorkerLog = "[Risk Notification Cleaner]"
	scheduleCleanNotificationsTime   = 24 * time.Hour
	maxAttemptsCleanNotificationJob  = 3
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

func (w *RiskNotificationCleanerWorker) Work(jobs chan dto.CleanNotificationJob) {
	log.Println(riskNotificationCleanerWorkerLog, "Start work")
	w.scheduleClean()

	for {
		job := <-jobs
		log.Println(riskNotificationCleanerWorkerLog, "Job received to clean", job.UserId, "notification")

		err := w.brokerRepo.PublishNotification(job.UserId, w.makeNotificationMessage())
		if err != nil {
			log.Println(riskNotificationCleanerWorkerLog, "Failed to publish notification:", err.Error())
			w.scheduleToPushBack(jobs, job)
		}
	}
}

func (w *RiskNotificationCleanerWorker) scheduleToPushBack(jobs chan<- dto.CleanNotificationJob, job dto.CleanNotificationJob) {
	time.AfterFunc(tryAgainAfter, func() {
		w.pushCleanNotificationJobBack(jobs, job)
	})
}

func (w *RiskNotificationCleanerWorker) pushCleanNotificationJobBack(jobs chan<- dto.CleanNotificationJob, job dto.CleanNotificationJob) {
	// Check attempts
	if job.Attempts >= maxAttemptsCleanNotificationJob {
		log.Println(riskNotificationCleanerWorkerLog, "Cannot try to push clean notification job back to queue, attempts ran out")
		return
	}

	// New attempt
	log.Println(riskNotificationCleanerWorkerLog, "Push clean notification job back to channel for new attempt")
	job.Attempts += 1
	jobs <- job
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

func AddCleanNotificationJob(userId string, channel chan<- dto.CleanNotificationJob) {
	log.Println(riskNotificationCleanerWorkerLog, "Add clean notification job for user", userId)

	job := dto.CleanNotificationJob{
		UserId:   userId,
		Attempts: 0,
	}

	channel <- job
}
