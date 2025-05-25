package main

import (
	"github.com/Laniakea00/e-commerce/api-gateway/config"
	"github.com/Laniakea00/e-commerce/api-gateway/router"
)

func main() {
	clients := config.InitClients()
	r := router.SetupRouter(clients)
	r.Run(":8080")
}
