package config

import (
	"consumer-service/domain"
	"consumer-service/usecase"
	"encoding/json"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

func SetupRouter() {
	// Подключение к NATS
	nc, err := nats.Connect("nats://natsuser:natspass@localhost:4222", nats.Timeout(10*time.Second))
	if err != nil {
		log.Fatalf("❌ Failed to connect to NATS: %v", err)
	}
	log.Println("✅ Connected to NATS")

	// Подключение к Inventory-service
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect to ProductService: %v", err)
	}
	log.Println("✅ Connected to ProductService via gRPC")

	// Создаем usecase и запускаем подписку
	reducer := usecase.NewStockReducer(conn)

	// Подписка на события "order.created"
	_, err = nc.Subscribe("order.created", func(msg *nats.Msg) {
		var order domain.Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("❌ Failed to unmarshal order: %v", err)
			return
		}

		log.Printf("📦 Received order ID %d with %d items", order.ID, len(order.Items))

		if err := reducer.ProcessOrder(order); err != nil {
			log.Printf("❌ Failed to process order: %v", err)
		} else {
			log.Printf("✅ Processed order %d", order.ID)
		}
	})

	if err != nil {
		log.Fatalf("❌ Failed to subscribe to NATS topic: %v", err)
	}
	log.Println("✅ Subscribed to order.created event")

	// Ожидаем завершения работы
	select {}
}
