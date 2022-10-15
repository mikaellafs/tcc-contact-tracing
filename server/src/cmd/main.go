package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"contacttracing/src/clients"
	"contacttracing/src/grpc/server"
	"contacttracing/src/grpc/service"
	"contacttracing/src/interfaces"
	"contacttracing/src/models/dto"
	"contacttracing/src/repositories"
	"contacttracing/src/workers"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

var (
	// Clients
	redisClient *redis.Client
	postgresDB  *sql.DB
	mqttClient  pahomqtt.Client

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

	initServer()
}

func initClients() {
	redisClient = clients.NewRedisClient()
	log.Println("Redis connection succeed")

	postgresDB = clients.NewPostgreSQLClient()
	log.Println("DB connection succeed")

	mqttClient = clients.NewMqttClient()
	log.Println("MQTT connection to broker succeed")
}

func closeClients() {
	postgresDB.Close()
	redisClient.Close()
	mqttClient.Disconnect(10)
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

func initServer() {
	grpcService := service.NewGrpcService(userRepo, reportRepo, cacheRepo, reportChan)
	s := server.NewGrpcCServer(grpcService)
	s.Serve()
}
