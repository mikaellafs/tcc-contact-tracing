package repositories

import (
	"log"
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

func (r *RedisRepository) SaveReport(userId, date string) {
	log.Println(cacheRepositoryLog, "Save report from user ", userId, ", diagnosticated at ", date)
	key := userId + "/" + date
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
