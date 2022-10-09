package main

import (
	"contacttracing/src/clients"
	"contacttracing/src/grpc/server"
	"contacttracing/src/repositories"
	"contacttracing/src/services"
	"context"
	"log"
)

func main() {
	postgresDB := clients.NewPostgreSQLClient()
	defer postgresDB.Close()

	userRepo := repositories.NewPostGreSQLUserRepository(postgresDB)
	err := userRepo.Migrate(context.TODO())
	if err != nil {
		log.Println(err.Error())
	}

	grpcService := services.NewGrpcService(userRepo)
	s := server.NewGrpcCServer(grpcService)
	s.Serve()
}
