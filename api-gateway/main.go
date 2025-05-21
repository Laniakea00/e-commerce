package main

import (
	"api-gateway/config"
	"api-gateway/router"
)

func main() {
	clients := config.InitClients()
	r := router.SetupRouter(clients)
	r.Run(":8080")
}
