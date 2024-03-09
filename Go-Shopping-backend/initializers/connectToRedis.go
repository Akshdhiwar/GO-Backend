package initializers

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

// RedisClient is the exported Redis client variable
var RedisClient *redis.Client

func ConnectToRedis() {
	var client *redis.Client
	if os.Getenv("ENVIRONMENT") == "LOCAL" {
		// Create a new Redis client
		client = redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379", // Redis server address
			Password: "",               // No password
			DB:       0,                // Use the default database
		})
		log.Println("getting redis from local")
	} else {
		redisURI := os.Getenv("RAILS_REDIS_URL")
		addr, err := redis.ParseURL(redisURI)
		if err != nil {
			panic(err)
		}
		client = redis.NewClient(addr)
		log.Println("getting redis from aiven")
	}

	log.Printf("From redis file")

	// Ping the Redis server to test the connection
	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	log.Printf("Connected to Redis: %s", pong)

	RedisClient = client
}
