package services

import (
	"contacttracing/src/grpc/pb"
	"context"
)

type GrpcService struct{}

func NewGrpcService() GrpcService {
	return GrpcService{}
}

func (s GrpcService) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResult, error) {
	//TODO
	return nil, nil
}

func (s GrpcService) ReportInfection(ctx context.Context, request *pb.ReportRequest) (*pb.ReportResult, error) {
	//TODO
	return nil, nil
}
