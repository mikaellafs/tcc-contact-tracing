package repositories

import (
	"contacttracing/src/models/dto"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	cacheRepositoryLog = "Cache Repository:"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	repo := new(RedisRepository)
	repo.client = client

	return repo
}

// TODO: change date to time.Time and convert to string inside the function
func (r *RedisRepository) SaveReport(userId, date string) {
	log.Println(cacheRepositoryLog, "Save report from user ", userId, ", diagnosticated at ", date)
	key := r.makeReportKey(userId, date)
	r.client.Set(key, date, 0)
}

func (r *RedisRepository) GetReportsFrom(userId string) []time.Time {
	log.Println(cacheRepositoryLog, "Get reports from user ", userId)

	iter := r.client.Scan(0, "prefix:"+userId, 0).Iterator()

	var dates []time.Time
	for iter.Next() {
		t, err := time.Parse(iter.Val(), time.RFC3339)
		if err == nil {
			dates = append(dates, t)
		}
	}

	return dates
}

// TODO: change date to time.Time and convert to string inside the function
func (r *RedisRepository) SaveNotification(userId string, fromReport int64, date string) {
	log.Println(cacheRepositoryLog, "Save notification for user", userId, "and report id", fromReport)
	key := r.makeNotificationKey(userId, fromReport)
	r.client.Set(key, date, 0)
}

func (r *RedisRepository) GetNotificationFrom(userId string, reportId int64) *dto.Notification {
	log.Println(cacheRepositoryLog, "Get notification from user", userId, "and report id", reportId)

	key := r.makeNotificationKey(userId, reportId)
	result := r.client.Get(key)

	if result == nil {
		return nil
	}

	dateNotification, err := time.Parse(result.Val(), time.RFC3339)
	if err != nil {
		return nil
	}

	return &dto.Notification{
		ReportId:     reportId,
		ForUser:      userId,
		DateNotified: dateNotification,
	}
}

func (r *RedisRepository) makeReportKey(userId, reportDate string) string {
	return userId + "/" + reportDate
}

func (r *RedisRepository) makeNotificationKey(userId string, reportId int64) string {
	return userId + "/" + strconv.FormatInt(reportId, 10)
}
