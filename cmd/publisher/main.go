package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"order-stream/internal/entity"
	"order-stream/pkg/nats_streaming/server"
)

func main() {
	// Создание нового сервера NATS Streaming
	natsServer, err := server.New(
		"localhost:4222",
		"test-cluster",
		"publisher",
		make(map[string]server.MsgHandler), // Пустой маршрутизатор, так как мы публикуем, а не подписываемся
		nil,                                // Без логгера, так как это простой пример
		server.Timeout(5*time.Second),
		server.ConnWaitTime(2*time.Second),
		server.ConnAttempts(3),
	)
	if err != nil {
		log.Fatalf("Failed to create NATS Streaming server: %v", err)
	}
	defer natsServer.Shutdown() // Закрытие ресурсов после завершения работы

	// Чтение заказа от пользователя
	fmt.Println("Enter order details (JSON format) and press Enter to publish:")
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read order: %v", err)
		}

		// Парсинг заказа
		var order entity.Order
		err = json.Unmarshal([]byte(input), &order)
		if err != nil {
			log.Fatalf("Failed to parse order: %v", err)
		}

		// Публикация заказа в NATS Streaming
		data, err := json.Marshal(order)
		if err != nil {
			log.Fatalf("Failed to marshal order: %v", err)
		}

		err = natsServer.Publish("orders", data)
		if err != nil {
			log.Fatalf("Failed to publish order: %v", err)
		}

		fmt.Println("Order published successfully!")
	}
}
