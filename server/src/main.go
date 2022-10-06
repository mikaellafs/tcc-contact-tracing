package main

import (
	"contacttracing/src/grpc/server"
	"contacttracing/src/services"
)

func main() {
	grpcService := services.NewGrpcService()
	s := server.NewGrpcCServer(grpcService)
	s.Serve()
}
