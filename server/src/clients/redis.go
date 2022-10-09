package clients

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"
)

func NewRedisClient() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})

	fmt.Println(host + ":" + port)

	pong, err := client.Ping().Result()
	log.Println("Testing redis connection: ", pong, err)

	if err != nil {
		panic(err)
	}

	return client
}
