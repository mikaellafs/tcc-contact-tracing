package service

import (
	"context"
	"log"
	"net/http"

	"contacttracing/src/grpc/pb"
	"contacttracing/src/interfaces"
	"contacttracing/src/models/db"
	"contacttracing/src/models/dto"
	"contacttracing/src/workers"
)

const (
	registerLog = "[Register]"
	reportLog   = "[Report]"
)

type GrpcService struct {
	userRepo       interfaces.UserRepository
	reportRepo     interfaces.ReportRepository
	cache          interfaces.CacheRepository
	tracingJobChan chan<- dto.ReportJob
}

func NewGrpcService(
	userRepo interfaces.UserRepository,
	reportRepo interfaces.ReportRepository,
	cache interfaces.CacheRepository,
	tracingJobChan chan<- dto.ReportJob) GrpcService {

	return GrpcService{userRepo: userRepo, reportRepo: reportRepo, cache: cache, tracingJobChan: tracingJobChan}
}

func (s GrpcService) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResult, error) {
	result := &pb.RegisterResult{
		Status:   http.StatusOK,
		ServerPk: "ok",
	}

	// Validate message
	r, err := validateGrpcMessage(request.GetRegister(), request.GetRegister().GetPk(), request.GetSignature())
	if err != nil {
		result.Status = r.Status
		result.Message = r.Message
		log.Println(registerLog, err.Error())
		return result, nil
	}

	// Save new user
	user := db.User{
		UserId:   request.GetRegister().GetUserId(),
		Pk:       request.GetRegister().GetPk(),
		Password: request.GetRegister().GetPassword(),
	}

	userSaved, err := s.userRepo.Create(context.TODO(), user)
	if err != nil {
		result.Status = http.StatusInternalServerError
		result.Message = err.Error()
		log.Println(registerLog, result.Message)
		return result, nil
	}

	log.Println(registerLog, userSaved)

	return result, nil
}

func (s GrpcService) ReportInfection(ctx context.Context, request *pb.ReportRequest) (*pb.ReportResult, error) {
	result := &pb.ReportResult{
		Status: http.StatusOK,
	}

	// Get user pk
	user, err := s.userRepo.GetByUserId(context.TODO(), request.GetReport().GetUserId())
	if err != nil {
		result.Status = http.StatusBadRequest
		result.Message = "Failed to get user from db: " + err.Error()
		log.Println(reportLog, result.Message)
		return result, nil
	}

	// Validate message
	r, err := validateGrpcMessage(parseGrpcReport(request.GetReport()), user.Pk, request.GetSignature())
	if err != nil {
		result.Status = r.Status
		result.Message = r.Message
		log.Println(reportLog, err.Error())
		return result, nil
	}

	// Create report in DB
	report, err := s.reportRepo.Create(context.TODO(), db.Report{
		UserId:         request.GetReport().GetUserId(),
		DateStart:      request.GetReport().GetDateStartSymptoms().AsTime(),
		DateDiagnostic: request.GetReport().GetDateDiagnostic().AsTime(),
		DateReport:     request.GetReport().GetDateReport().AsTime(),
	})
	if err != nil {
		result.Status = http.StatusInternalServerError
		result.Message = "Could not create report in DB: " + err.Error()
		log.Println(reportLog, result.Message)
		return result, nil
	}

	// Save report at risk cache
	s.cache.SaveReport(report.UserId, report.ID, report.DateDiagnostic)

	// Add job to trace contacts
	go workers.AddReportJob(report.ID, report.UserId, report.DateDiagnostic, s.tracingJobChan)

	result.Message = "Reported infection. Contacts are going to get notified."
	return result, nil
}
