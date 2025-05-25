package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"user-service/config"
	"user-service/redis"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2113", nil) // порт может быть разный
	}()
	redis.InitRedis()
	db := config.InitDB()
	if err := config.SetupRouter(db); err != nil {
		log.Fatalf("failed to run user service: %v", err)
	}
}
