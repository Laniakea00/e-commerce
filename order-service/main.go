package main

import (
	"order-service/config"
)

func main() {
	db := config.InitDB()
	r := config.SetupRouter(db)
	r.Run(":8082")
}
