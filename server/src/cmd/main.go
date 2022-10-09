package main

import (
	"contacttracing/src/clients"
	"contacttracing/src/grpc/server"
	"contacttracing/src/repositories"
	"contacttracing/src/services"
	"context"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}

	redisClient := clients.NewRedisClient()
	cacheRepo := repositories.NewRedisRepository(redisClient)

	postgresDB := clients.NewPostgreSQLClient()
	defer postgresDB.Close()
	log.Println("DB connection succeed")

	userRepo := repositories.NewPostGreSQLUserRepository(postgresDB)
	err = userRepo.Migrate(context.TODO())
	if err != nil {
		log.Println(err.Error())
		return
	}

	reportRepo := repositories.NewPostGreSQLReportRepository(postgresDB)
	err = reportRepo.Migrate(context.TODO())
	if err != nil {
		log.Println(err.Error())
		return
	}

	grpcService := services.NewGrpcService(userRepo, reportRepo, cacheRepo)
	s := server.NewGrpcCServer(grpcService)
	s.Serve()
}
