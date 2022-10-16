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

const (
	maxDiffDaysFromDiagnosticToConsiderAtRisk = 15
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
	mqttRepo         interfaces.BrokerRepository

	// Channels
	reportChan        chan dto.ReportJob
	notifChan         chan dto.NotificationJob
	potentialRiskChan chan dto.PotentialRiskJob
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

	mqttRepo = repositories.NewMqttRepository(mqttClient)
}

func initWorkers() {
	reportChan = make(chan dto.ReportJob)
	notifChan = make(chan dto.NotificationJob)
	potentialRiskChan = make(chan dto.PotentialRiskJob)

	go workers.NewContacTracerWorker(contactRepo, maxDiffDaysFromDiagnosticToConsiderAtRisk).Work(reportChan, notifChan)
	go workers.NewRiskNotifierWorker(notificationRepo, cacheRepo, mqttRepo, 5*time.Minute, maxDiffDaysFromDiagnosticToConsiderAtRisk).Work(notifChan)
}

func initServer() {
	grpcService := service.NewGrpcService(userRepo, reportRepo, cacheRepo, reportChan)
	s := server.NewGrpcCServer(grpcService)
	s.Serve()
}
