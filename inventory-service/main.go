package main

import (
	"inventory-service/config"
)

func main() {
	db := config.InitDB()
	r := config.SetupRouter(db)
	r.Run(":8081")
}
