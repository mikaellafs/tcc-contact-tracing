package workers

import (
	"contacttracing/src/models/dto"
	"log"
)

const (
	potentialRiskTracerWorkerLog = "[Potential Risk Tracer Worker]"
)

func AddPotentialRiskJob(user string, reportId int64, channel chan<- dto.PotentialRiskJob) {
	job := dto.PotentialRiskJob{
		User:     user,
		ReportId: reportId,
	}

	log.Println(potentialRiskTracerWorkerLog, "Add potential risk job:", job)
	channel <- job
}
