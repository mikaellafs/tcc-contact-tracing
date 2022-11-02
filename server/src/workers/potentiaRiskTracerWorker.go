package workers

import (
	"contacttracing/src/interfaces"
	"contacttracing/src/models/dto"
	"contacttracing/src/utils"
	"context"
	"log"
	"time"
)

const (
	potentialRiskTracerWorkerLog = "[Potential Risk Tracer Worker]"
	maxAttemptsPotentialRiskJob  = 10
	scheduleRiskJobTime          = 1 * time.Hour
)

type PotentialRiskTracerWorker struct {
	contactRepo            interfaces.ContactRepository
	reportRepo             interfaces.ReportRepository
	cacheRepo              interfaces.CacheRepository
	riskContactMinDuration time.Duration
}

func NewPotentialRiskTracerWorker(
	contactRepo interfaces.ContactRepository,
	reportRepo interfaces.ReportRepository,
	cacheRepo interfaces.CacheRepository,
	riskContactMinDuration time.Duration) *PotentialRiskTracerWorker {

	worker := new(PotentialRiskTracerWorker)
	worker.contactRepo = contactRepo
	worker.reportRepo = reportRepo
	worker.cacheRepo = cacheRepo
	worker.riskContactMinDuration = riskContactMinDuration

	return worker
}

func (w *PotentialRiskTracerWorker) Work(potentialRisks chan dto.PotentialRiskJob, notifications chan<- dto.NotificationJob) {
	log.Println(potentialRiskTracerWorkerLog, "Start work")
	for {
		// Wait for potential risk
		job := <-potentialRisks
		log.Println(potentialRiskTracerWorkerLog, "Potential risk received:", job)

		// Get report info
		report, err := w.reportRepo.GetById(context.TODO(), job.ReportId)

		// Push job back to 'queue' if some error ocurred for new attempt
		if err != nil {
			log.Println(potentialRiskTracerWorkerLog, "Failed to get report from db: ", err.Error())
			w.scheduleToPushBack(potentialRisks, job)
			continue
		}

		// Trace contacts between user and infected user
		contacts, err := w.contactRepo.GetContactsBetweenUsers(context.TODO(), report.UserId, job.User, report.DateDiagnostic, time.Now())
		if err != nil {
			log.Println(potentialRiskTracerWorkerLog, "Failed to get contacts from db: ", err.Error())
			w.scheduleToPushBack(potentialRisks, job)
			continue
		}

		// Get longest contact
		longestContact := utils.GetLongestContact(contacts)

		// Check contact duration to notify
		if longestContact.Duration >= w.riskContactMinDuration {
			log.Println(potentialRiskTracerWorkerLog, "Risky contact with duration of", w.riskContactMinDuration, "minutes.")

			// Notify about contact
			go AddNotificationJob(longestContact.DateLastContact, job.User, job.ReportId, longestContact.Duration, notifications)
		}

		// Remove potential risk job from cache
		w.cacheRepo.RemovePotentialRiskJob(job.User, job.ReportId)
	}
}

func (w *PotentialRiskTracerWorker) scheduleToPushBack(potentialRisks chan<- dto.PotentialRiskJob, job dto.PotentialRiskJob) {
	time.AfterFunc(tryAgainAfter, func() {
		w.pushPotentialRiskBack(potentialRisks, job)
	})
}

func (w *PotentialRiskTracerWorker) pushPotentialRiskBack(potentialRisks chan<- dto.PotentialRiskJob, potentialRiskJob dto.PotentialRiskJob) {
	// Check attempts
	if potentialRiskJob.Attempts >= maxAttemptsPotentialRiskJob {
		log.Println(potentialRiskTracerWorkerLog, "Cannot try to push potential risk job back to queue, attempts ran out")
		w.cacheRepo.RemovePotentialRiskJob(potentialRiskJob.User, potentialRiskJob.ReportId)
		return
	}

	// New attempt
	log.Println(potentialRiskTracerWorkerLog, "Push potential risk job back to channel for new attempt")
	potentialRiskJob.Attempts += 1
	potentialRisks <- potentialRiskJob
}

func AddPotentialRiskJob(user string, reportId int64, channel chan<- dto.PotentialRiskJob, cache interfaces.CacheRepository) {
	job := dto.PotentialRiskJob{
		User:     user,
		ReportId: reportId,
		Attempts: 0,
	}

	// Check if there's already a job scheduled
	isThereAJob := cache.GetPotentialRiskJob(user, reportId)
	if isThereAJob {
		log.Println(potentialRiskTracerWorkerLog, "There's already a potential risk job for user:", user)
		return
	}

	log.Println(potentialRiskTracerWorkerLog, "Schedule new potential risk job")
	cache.SavePotentialRiskJob(user, reportId)

	// Scheudle job
	time.AfterFunc(scheduleRiskJobTime, func() {
		log.Println(potentialRiskTracerWorkerLog, "Add potential risk job:", job)
		channel <- job
	})
}
