package workers

import (
	"contacttracing/src/models/dto"
	"log"
	"time"
)

const (
	riskNotifierWorkerLog = "Risk Notifier Worker:"
)

func AddNotificationJob(userNotified, userInfected string, fromReport int64, contactDuration time.Duration, notifChan chan<- dto.NotificationJob) {
	notificationJob := dto.NotificationJob{
		ForUser:         userNotified,
		FromContactWith: userInfected,
		FromReport:      fromReport,
		Duration:        contactDuration,
		Attempts:        0,
	}

	log.Println(riskNotifierWorkerLog, "Add notification job:", notificationJob)
	notifChan <- notificationJob
}
