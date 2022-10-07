package main

import (
	"contacttracing/src/clients"
	"contacttracing/src/grpc/server"
	"contacttracing/src/services"
)

func main() {
	postgresDB := clients.NewPostgreSQLClient()
	defer postgresDB.Close()

	grpcService := services.NewGrpcService()
	s := server.NewGrpcCServer(grpcService)
	s.Serve()
}
