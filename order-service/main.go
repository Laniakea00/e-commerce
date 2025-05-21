package main

import (
	"log"
	"order-service/config"
	"order-service/redis"
)

func main() {
	redis.InitRedis()
	db := config.InitDB()
	if err := config.SetupRouter(db); err != nil {

		log.Fatalf("failed to run order service: %v", err)
	}
}
