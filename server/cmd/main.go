package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"contacttracing/pkg/clients"
	"contacttracing/pkg/grpc/server"
	"contacttracing/pkg/grpc/service"
	"contacttracing/pkg/interfaces"
	"contacttracing/pkg/models/dto"
	"contacttracing/pkg/repositories"
	"contacttracing/pkg/workers"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

const (
	maxDiffDaysFromDiagnosticToConsiderAtRisk = 15
	contactsTopic                             = "contato"
	reportExpiration                          = 15 * time.Hour * 24 // 15 days
	minContactDuration                        = 15 * time.Minute
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

	f := logsToFile()
	defer f.Close()

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

	go workers.NewContacTracerWorker(contactRepo, maxDiffDaysFromDiagnosticToConsiderAtRisk, minContactDuration).Work(reportChan, notifChan, cleanNotifChannel)
	go workers.NewRiskNotifierWorker(notificationRepo, cacheRepo, mqttRepo, maxDiffDaysFromDiagnosticToConsiderAtRisk).Work(notifChan)
	go workers.NewPotentialRiskTracerWorker(contactRepo, reportRepo, cacheRepo, minContactDuration).Work(potentialRiskChan, notifChan)
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

func logsToFile() *os.File {
	file := os.Getenv("LOGFILE_PATH")
	fmt.Println("File to save logs:", os.Getenv("LOGFILE_PATH"))

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)

	return f
}
