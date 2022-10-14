package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"contacttracing/src/clients"
	"contacttracing/src/grpc/server"
	"contacttracing/src/interfaces"
	"contacttracing/src/models/dto"
	"contacttracing/src/repositories"
	"contacttracing/src/services"
	"contacttracing/src/workers"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

var (
	// Clients
	redisClient *redis.Client
	postgresDB  *sql.DB

	// Repositories
	cacheRepo        interfaces.CacheRepository
	userRepo         interfaces.UserRepository
	reportRepo       interfaces.ReportRepository
	contactRepo      interfaces.ContactRepository
	notificationRepo interfaces.NotificationRepository

	// Channels
	reportChan chan dto.ReportJob
	notifChan  chan dto.NotificationJob
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}

	initClients()
	defer closeClients()

	initRepositories()

	initWorkers()

	grpcService := services.NewGrpcService(userRepo, reportRepo, cacheRepo, reportChan)
	s := server.NewGrpcCServer(grpcService)
	s.Serve()
}

func initClients() {
	redisClient = clients.NewRedisClient()
	postgresDB = clients.NewPostgreSQLClient()
	log.Println("DB connection succeed")
}

func closeClients() {
	postgresDB.Close()
	redisClient.Close()
}

func initRepositories() {
	cacheRepo = repositories.NewRedisRepository(redisClient)

	userRepo = repositories.NewPostGreSQLUserRepository(postgresDB)
	err := userRepo.Migrate(context.TODO())
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}

	reportRepo = repositories.NewPostGreSQLReportRepository(postgresDB)
	err = reportRepo.Migrate(context.TODO())
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}

	contactRepo = repositories.NewPostgreSQLContactRepository(postgresDB)
	err = contactRepo.Migrate(context.TODO())
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}

	notificationRepo = repositories.NewPostgreSQLNotificationRepository(postgresDB)
	err = notificationRepo.Migrate(context.TODO())
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
}

func initWorkers() {
	reportChan = make(chan dto.ReportJob)
	notifChan = make(chan dto.NotificationJob)

	go workers.NewContacTracerWorker(contactRepo, 30).Work(reportChan, notifChan)
	go workers.NewRiskNotifierWorker(notificationRepo, cacheRepo, 5*time.Minute).Work(notifChan)
}
