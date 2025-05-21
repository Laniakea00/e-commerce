package main

import (
	"log"
	"user-service/config"
	"user-service/redis"
)

func main() {
	redis.InitRedis()
	db := config.InitDB()
	if err := config.SetupRouter(db); err != nil {
		log.Fatalf("failed to run user service: %v", err)
	}
}
