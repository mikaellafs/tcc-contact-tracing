package server

import (
	"contacttracing/src/grpc/pb"
	"contacttracing/src/services"
	"log"
	"net"

	"google.golang.org/grpc"
)

type ContactTracingGrpcServer struct {
	services.GrpcService
}

func NewGrpcCServer(s services.GrpcService) ContactTracingGrpcServer {
	return ContactTracingGrpcServer{s}
}

func (g ContactTracingGrpcServer) Serve() {
	lis, err := net.Listen("tcp", ":50052") // TODO: get host and port from .env
	if err != nil {
		log.Fatalf("Could not listen tpc port")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterContactTracingServer(grpcServer, &g)

	log.Println("Listening on port 50052")

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve: " + err.Error())
	}
}
