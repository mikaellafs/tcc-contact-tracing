package repositories

import (
	"contacttracing/src/models/dto"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	cacheRepositoryLog = "[Cache Repository]"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	repo := new(RedisRepository)
	repo.client = client

	return repo
}

func (r *RedisRepository) SaveReport(userId string, date time.Time) {
	log.Println(cacheRepositoryLog, "Save report from user ", userId, ", diagnosticated at ", date)
	key := r.makeReportKey(userId, date.Format(time.RFC3339))
	r.client.Set(key, date, 0)
}

func (r *RedisRepository) GetReportsFrom(userId string) []time.Time {
	log.Println(cacheRepositoryLog, "Get reports from user ", userId)

	iter := r.client.Scan(0, "prefix:report:"+userId, 0).Iterator()

	var dates []time.Time
	for iter.Next() {
		t, err := time.Parse(time.RFC3339, iter.Val())
		if err == nil {
			dates = append(dates, t)
		}
	}

	return dates
}

func (r *RedisRepository) SaveNotification(userId string, fromReport int64, date time.Time) {
	log.Println(cacheRepositoryLog, "Save notification for user", userId, "and report id", fromReport)
	key := r.makeNotificationKey(userId, fromReport)
	log.Println(cacheRepositoryLog, date.Format(time.RFC3339))
	r.client.Set(key, date.Format(time.RFC3339), 0)
}

func (r *RedisRepository) GetNotificationFrom(userId string, reportId int64) *dto.Notification {
	log.Println(cacheRepositoryLog, "Get notification from user", userId, "and report id", reportId)

	key := r.makeNotificationKey(userId, reportId)
	value, err := r.client.Get(key).Result()

	if err != nil {
		log.Println(cacheRepositoryLog, err.Error())
		return nil
	}

	log.Println(cacheRepositoryLog, "Value found:", value, "for key:", key)
	dateNotification, err := time.Parse(time.RFC3339, value)
	if err != nil {
		log.Println(cacheRepositoryLog, err.Error())
		return nil
	}

	return &dto.Notification{
		ReportId:     reportId,
		ForUser:      userId,
		DateNotified: dateNotification,
	}
}

func (r *RedisRepository) makeReportKey(userId, reportDate string) string {
	return "report:" + userId + "/" + reportDate
}

func (r *RedisRepository) makeNotificationKey(userId string, reportId int64) string {
	return "notificaton:" + userId + "/" + strconv.FormatInt(reportId, 10)
}
