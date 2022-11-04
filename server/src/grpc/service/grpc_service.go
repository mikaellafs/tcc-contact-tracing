package service

import (
	"context"
	"log"
	"net/http"

	"contacttracing/src/grpc/pb"
	"contacttracing/src/interfaces"
	"contacttracing/src/models/db"
	"contacttracing/src/models/dto"
	"contacttracing/src/utils"
	"contacttracing/src/workers"

	"github.com/google/uuid"
)

const (
	registerLog = "[Register]"
	reportLog   = "[Report]"
)

type GrpcService struct {
	pb.UnimplementedContactTracingServer
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
		Status: http.StatusOK,
	}

	// Check params
	if request.GetDeviceId() == "" || request.GetPk() == "" {
		result.Status = http.StatusBadRequest
		result.Message = "Missing deviceId or public key"
		log.Println(registerLog, result.Message)
		return result, nil
	}

	// Encrypt deviceId
	deviceId := utils.EncryptStr(request.GetDeviceId())

	// Save new user
	user := db.User{
		Id:       uuid.New().String(),
		DeviceId: deviceId,
		Pk:       request.GetPk(),
	}

	userSaved, err := s.userRepo.Create(context.TODO(), user)
	if err != nil {
		result.Status = http.StatusInternalServerError
		result.Message = err.Error()
		log.Println(registerLog, result.Message)
		return result, nil
	}

	log.Println(registerLog, userSaved)
	result.UserId = userSaved.Id

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
