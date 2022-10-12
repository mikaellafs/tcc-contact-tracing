package main

import (
	"context"
	"log"

	"contacttracing/src/clients"
	"contacttracing/src/grpc/server"
	"contacttracing/src/models/dto"
	"contacttracing/src/repositories"
	"contacttracing/src/services"
	"contacttracing/src/workers"

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

	contactRepo := repositories.NewPostgreSQLContactRepository(postgresDB)
	err = contactRepo.Migrate(context.TODO())

	var reportChan chan dto.ReportJob
	var notifChan chan int
	go workers.NewContacTracerWorker(contactRepo, 30).Work(reportChan, notifChan)

	grpcService := services.NewGrpcService(userRepo, reportRepo, cacheRepo, reportChan)
	s := server.NewGrpcCServer(grpcService)
	s.Serve()
}
