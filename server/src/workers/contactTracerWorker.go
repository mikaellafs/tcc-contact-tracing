package workers

import (
	"contacttracing/src/interfaces"
	"contacttracing/src/models/dto"
	"context"
	"log"
	"time"
)

const (
	contactTracerWorkerLog = "Contact Tracer Worker:"
	maxAttemptsReportJob   = 5
	tryAgainAfter          = 5 * time.Second
)

type ContactTracerWorker struct {
	contactRepo interfaces.ContactRepository
	daysToTrace int
}

func NewContacTracerWorker(repo interfaces.ContactRepository, days int) *ContactTracerWorker {
	return &ContactTracerWorker{contactRepo: repo, daysToTrace: days}
}

func (w *ContactTracerWorker) Work(reports chan dto.ReportJob, notifications chan<- dto.NotificationJob) {
	log.Println(contactTracerWorkerLog, "Start work")
	for {
		// Wait for report
		report := <-reports
		log.Println(contactTracerWorkerLog, "Report received: ", report)

		// Get contacts
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(150)*time.Millisecond)
		defer cancel()

		contacts, err := w.contactRepo.GetContactsWithin(ctx, w.daysToTrace, report.Date, report.UserId)

		// Push report back to 'queue' if some error ocurred for new attempt
		if err != nil {
			log.Println(contactTracerWorkerLog, "Failed to get contacts from db: ", err.Error())

			time.AfterFunc(tryAgainAfter, func() {
				w.pushReportBack(reports, report)
			})

			continue
		}

		// Create notification job for each contact
		go w.notifyContacts(contacts, report, notifications)
	}
}

func (w *ContactTracerWorker) pushReportBack(reports chan<- dto.ReportJob, report dto.ReportJob) {
	// Check attempts
	if report.Attempts >= maxAttemptsReportJob {
		log.Println(contactTracerWorkerLog, "Cannot try to push report back to queue, attempts ran out")
		return
	}

	// New attempt
	log.Println(contactTracerWorkerLog, "Push report back to channel for new attempt")
	report.Attempts += 1
	reports <- report
}

func (w *ContactTracerWorker) notifyContacts(contacts []dto.Contact, report dto.ReportJob, channel chan<- dto.NotificationJob) {
	for _, contact := range contacts {
		log.Println(contactTracerWorkerLog, "Contato com:", contact.AnotherUser, "por", contact.Duration.Minutes(), "minutos")
		AddNotificationJob(contact.AnotherUser, contact.User, report.DBID, contact.Duration, channel)
	}
}

func AddReportJob(dbID int64, userId string, date time.Time, channel chan<- dto.ReportJob) {
	reportJob := dto.ReportJob{
		DBID:     dbID,
		UserId:   userId,
		Date:     date,
		Attempts: 0,
	}

	log.Println(contactTracerWorkerLog, "Add report job: ", reportJob)
	channel <- reportJob
}
