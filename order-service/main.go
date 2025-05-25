package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"order-service/config"
	"order-service/redis"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2114", nil) // порт может быть разный
	}()
	redis.InitRedis()
	db := config.InitDB()
	if err := config.SetupRouter(db); err != nil {

		log.Fatalf("failed to run order service: %v", err)
	}
}
