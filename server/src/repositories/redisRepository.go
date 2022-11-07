package repositories

import (
	"contacttracing/src/models/dto"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

const (
	cacheRepositoryLog = "[Cache Repository]"
)

type RedisRepository struct {
	client           *redis.Client
	reportExpiration time.Duration
}

func NewRedisRepository(client *redis.Client, reportExpiration time.Duration) *RedisRepository {
	repo := new(RedisRepository)
	repo.client = client
	repo.reportExpiration = reportExpiration

	return repo
}

func (r *RedisRepository) SaveReport(userId string, reportId int64, date time.Time) {
	log.Println(cacheRepositoryLog, "Save report from user ", userId, ", diagnosticated at ", date)

	key := r.makeReportKey(userId, reportId)
	r.client.Set(key, date.Format(time.RFC3339), r.reportExpiration)
}

func (r *RedisRepository) GetReportsFrom(userId string) []dto.Report {
	log.Println(cacheRepositoryLog, "Get reports from user ", userId)

	prefix := "report:" + userId
	iter := r.client.Scan(0, prefix+"*", 0).Iterator()

	reports := r.parseScanReportsResults(iter)

	return reports
}

func (r *RedisRepository) SaveNotification(userId string, fromReport int64, date time.Time) {
	log.Println(cacheRepositoryLog, "Save notification for user", userId, "and report id", fromReport)

	key := r.makeNotificationKey(userId, fromReport)
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

func (r *RedisRepository) GetNotificationKeysFrom(userId string) []string {
	log.Println(cacheRepositoryLog, "Get notifications keys from user", userId)

	prefix := "notification:"
	keys, _, _ := r.client.Scan(0, prefix+"*", 0).Result()

	return keys
}

func (r *RedisRepository) RemoveNotificationsAfter(days time.Duration, maxDate time.Time) []string {
	log.Println(cacheRepositoryLog, "Remove notifications after", days, "days")

	var expiredNotificationsUserIds []string
	prefix := "notification:"
	iter := r.client.Scan(0, prefix+"*", 0).Iterator()

	for iter.Next() {
		key, user, date, err := r.extractNotification(iter.Val())
		if err != nil {
			continue
		}

		// Check if it is expired
		if date.Add(days).After(maxDate) {
			log.Println(cacheRepositoryLog, "Expired notification for user", user)

			r.client.Del(key)
			expiredNotificationsUserIds = append(expiredNotificationsUserIds, user)
		}
	}

	return expiredNotificationsUserIds
}

func (r *RedisRepository) SavePotentialRiskJob(userId string, reportId int64) {
	log.Println(cacheRepositoryLog, "Save potential risk job for user", userId, "and report", reportId)

	key := r.makePotentialRiskJobKey(userId, reportId)
	r.client.Set(key, "ok", 0)
}

func (r *RedisRepository) GetPotentialRiskJob(userId string, reportId int64) bool {
	log.Println(cacheRepositoryLog, "Get potential risk job for user", userId, "and report", reportId)

	key := r.makePotentialRiskJobKey(userId, reportId)
	_, err := r.client.Get(key).Result()
	if err != nil {
		log.Println(cacheRepositoryLog, err.Error())
		return false
	}

	return true
}

func (r *RedisRepository) RemovePotentialRiskJob(userId string, reportId int64) {
	log.Println(cacheRepositoryLog, "Remove potential risk job for user", userId, "and report", reportId)

	key := r.makePotentialRiskJobKey(userId, reportId)
	r.client.Del(key)
}

func (r *RedisRepository) makeReportKey(userId string, reportId int64) string {
	return "report:" + userId + "#" + strconv.FormatInt(reportId, 10)
}

func (r *RedisRepository) makeNotificationKey(userId string, reportId int64) string {
	return "notificaton:" + userId + "#" + strconv.FormatInt(reportId, 10)
}

func (r *RedisRepository) makePotentialRiskJobKey(userId string, reportId int64) string {
	return "potentialriskjob:" + userId + "#" + strconv.FormatInt(reportId, 10)
}

func (r *RedisRepository) parseScanReportsResults(iter *redis.ScanIterator) []dto.Report {
	var reports []dto.Report
	for iter.Next() {
		log.Println(cacheRepositoryLog, iter.Val())

		reportId, date, err := r.extractReport(iter.Val())

		if err == nil {
			reports = append(reports, dto.Report{
				ID:             reportId,
				DateDiagnostic: date,
			})
		}
	}

	return reports
}

func (r *RedisRepository) extractReport(key string) (int64, time.Time, error) {
	value, err := r.client.Get(key).Result()
	if err != nil {
		return 0, time.Now(), err
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return 0, t, err
	}

	reportIdStr := strings.Split(key, "#")[1]
	rId, err := strconv.ParseInt(reportIdStr, 10, 64)
	if err != nil {
		return 0, t, err
	}

	return rId, t, nil
}

func (r *RedisRepository) extractNotification(key string) (string, string, time.Time, error) {
	value, err := r.client.Get(key).Result()
	if err != nil {
		return "", "", time.Now(), err
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return "", "", t, err
	}

	userId := strings.Split(key, "#")[0]
	userId = strings.ReplaceAll(userId, "notification:", "")

	return key, userId, t, nil
}
