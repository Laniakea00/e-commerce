package main

import (
	"consumer-service/config"
)

func main() {
	config.SetupRouter()
	select {} // блокируем, чтобы процесс не завершался
}
