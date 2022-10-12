package services

import (
	"context"
	"log"
	"net/http"
	"time"

	"contacttracing/src/grpc/pb"
	"contacttracing/src/interfaces"
	"contacttracing/src/models/db"
	"contacttracing/src/models/dto"
	"contacttracing/src/utils"
	"contacttracing/src/workers"
)

const (
	registerLog = "Register: "
	reportLog   = "Report: "
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
	isValid, err := utils.ValidateMessage(request.GetRegister(), request.GetRegister().GetPk(), request.GetSignature())
	if err != nil {
		result.Status = http.StatusBadRequest
		result.Message = "Failed to validate message: " + err.Error()
		log.Println(registerLog, result.Message)
		return result, nil
	}

	if !isValid {
		result.Status = http.StatusForbidden
		result.Message = "Signature is not valid for this message"
		log.Println(registerLog, result.Message)
		return result, nil
	}

	// Save new user
	user := db.User{
		UserId:   request.GetRegister().GetUserId(),
		Pk:       request.GetRegister().GetPk(),
		Password: request.GetRegister().GetPassword(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(150)*time.Millisecond)
	defer cancel()

	userSaved, err := s.userRepo.Create(ctx, user)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(150)*time.Millisecond)
	defer cancel()

	// Get user pk
	user, err := s.userRepo.GetByUserId(ctx, request.GetReport().GetUserId())
	if err != nil {
		result.Status = http.StatusBadRequest
		result.Message = "Failed to get user from db: " + err.Error()
		log.Println(reportLog, result.Message)
		return result, nil
	}

	// Validate message
	isValid, err := utils.ValidateMessage(request.GetReport(), user.Pk, request.GetSignature())
	if err != nil {
		result.Status = http.StatusBadRequest
		result.Message = "Failed to validate message: " + err.Error()
		log.Println(reportLog, result.Message)
		return result, nil
	}

	if !isValid {
		result.Status = http.StatusForbidden
		result.Message = "Signature is not valid for this message"
		log.Println(reportLog, result.Message)
		return result, nil
	}

	// Create report in DB
	report, err := s.reportRepo.Create(ctx, db.Report{
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
	s.cache.SaveReport(report.UserId, report.DateDiagnostic.Format(time.RFC3339))

	// Add job to trace contacts
	go workers.AddReportJob(report.UserId, report.DateDiagnostic, s.tracingJobChan)

	result.Message = "Reported infection. Contacts are going to get notified."
	return result, nil
}
