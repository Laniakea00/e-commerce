package main

import (
	"api-gateway/router"
)

func main() {
	r := router.SetupRouter()
	r.Run(":8080")
}
