package services

import (
	"contacttracing/src/grpc/pb"
	"contacttracing/src/interfaces"
	"contacttracing/src/models/db"
	"contacttracing/src/utils"
	"context"
	"log"
	"net/http"
	"time"
)

type GrpcService struct {
	interfaces.UserRepository
}

func NewGrpcService(userRepo interfaces.UserRepository) GrpcService {
	return GrpcService{userRepo}
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
		log.Println(result.Message)
		return result, nil
	}

	if !isValid {
		result.Status = http.StatusForbidden
		result.Message = "Signature is not valid for this message"
		log.Println(result.Message)
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

	userSaved, err := s.Create(ctx, user)
	if err != nil {
		result.Status = http.StatusInternalServerError
		result.Message = err.Error()
		log.Println(result.Message)
		return result, nil
	}

	log.Println(userSaved)

	return result, nil
}

func (s GrpcService) ReportInfection(ctx context.Context, request *pb.ReportRequest) (*pb.ReportResult, error) {
	//TODO
	return nil, nil
}
