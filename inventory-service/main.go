package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"inventory-service/config"
	"inventory-service/redis"
	"log"
	"net/http"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2115", nil) // порт может быть разный
	}()
	redis.InitRedis()
	db := config.InitDB()
	if err := config.SetupRouter(db); err != nil {
		log.Fatalf("Failed to run inventory service: %v", err)
	}
}
