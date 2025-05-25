package main

import (
	"github.com/Laniakea00/e-commerce/api-gateway/config"
	"github.com/Laniakea00/e-commerce/api-gateway/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil) // порт может быть разный
	}()
	clients := config.InitClients()
	r := router.SetupRouter(clients)
	r.Run(":8080")
}
