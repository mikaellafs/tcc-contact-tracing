package workers

import (
	"contacttracing/src/interfaces"
	"contacttracing/src/models/db"
	"contacttracing/src/models/dto"
	"context"
	"log"
	"time"
)

const (
	contactTracerWorkerLog = "Contact Tracer Worker: "
	max_attempts           = 5
)

type ContactTracerWorker struct {
	contactRepo interfaces.ContactRepository
	daysToTrace int
}

func NewContacTracerWorker(repo interfaces.ContactRepository, days int) *ContactTracerWorker {
	return &ContactTracerWorker{contactRepo: repo, daysToTrace: days}
}

func (w *ContactTracerWorker) NewJob(userId string, date time.Time) dto.ReportJob {
	return dto.ReportJob{
		UserId:   userId,
		Date:     date,
		Attempts: 0,
	}
}

func (w *ContactTracerWorker) Work() (reports chan dto.ReportJob, notifications chan<- int) {
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
			w.pushReportBack(reports, report)
			continue
		}

		// Create notification for each contact
		for _, contact := range contacts {
			w.sendNotification(notifications, contact)
		}
	}
}

func (w *ContactTracerWorker) pushReportBack(reports chan<- dto.ReportJob, report dto.ReportJob) {
	// Check attempts
	if report.Attempts >= max_attempts {
		log.Println(contactTracerWorkerLog, "Cannot try to push report back to queue, attempts ran out")
		return
	}

	// New attempt
	log.Println(contactTracerWorkerLog, "Push report back to channel for new attempt")
	report.Attempts += 1
	reports <- report
}

// TODO: implement send notifications
func (w *ContactTracerWorker) sendNotification(notifications chan<- int, contact db.Contact) {

}
