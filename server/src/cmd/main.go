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
	contactsTopic                             = "contact"
	reportExpiration                          = 15 * time.Hour * 24 // 15 days
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
	cleanNotifChannel chan dto.CleanNotificationJob
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
	cacheRepo = repositories.NewRedisRepository(redisClient, reportExpiration)

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
	reportChan = make(chan dto.ReportJob, 50)
	notifChan = make(chan dto.NotificationJob, 50)
	potentialRiskChan = make(chan dto.PotentialRiskJob, 50)
	cleanNotifChannel = make(chan dto.CleanNotificationJob, 50)

	go workers.NewContacTracerWorker(contactRepo, maxDiffDaysFromDiagnosticToConsiderAtRisk).Work(reportChan, notifChan, cleanNotifChannel)
	go workers.NewRiskNotifierWorker(notificationRepo, cacheRepo, mqttRepo, 5*time.Minute, maxDiffDaysFromDiagnosticToConsiderAtRisk).Work(notifChan)
	go workers.NewPotentialRiskTracerWorker(contactRepo, reportRepo, cacheRepo).Work(potentialRiskChan, notifChan)
	go workers.NewRiskNotificationCleanerWorker(cacheRepo, mqttRepo, reportExpiration).Work(cleanNotifChannel)
}

func initServer() {
	// Init mqtt broker handler
	contactsProcessor := workers.NewContactsProcessor(contactRepo, userRepo, cacheRepo, potentialRiskChan, maxDiffDaysFromDiagnosticToConsiderAtRisk)
	mqttRepo.SubscribeToReceiveContacts(contactsTopic, contactsProcessor.Process)

	// Init grpc server
	grpcService := service.NewGrpcService(userRepo, reportRepo, cacheRepo, reportChan)
	s := server.NewGrpcCServer(grpcService)
	s.Serve()
}
