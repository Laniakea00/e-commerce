package main

import (
	"inventory-service/config"
	"inventory-service/redis"
	"log"
)

func main() {
	redis.InitRedis()
	db := config.InitDB()
	if err := config.SetupRouter(db); err != nil {
		log.Fatalf("Failed to run inventory service: %v", err)
	}
}
