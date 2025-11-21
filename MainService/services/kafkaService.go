package services

import (
	. "L0_WB/models"
	. "L0_WB/repository"
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func GetDataFromKafka(mas map[string]Order) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "test-topic",
		GroupID: "my-groupID",
	})
	defer reader.Close()
	initDb()
	fmt.Println("Start get data")
	for {
		fmt.Println("Start sycle get data")
		var order Order
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Ошибка при получении:", err)
			panic("Ошибка при получении")
		}

		fmt.Println("Получили: ", string(msg.Value))

		if err = json.Unmarshal(msg.Value, &order); err != nil {
			fmt.Errorf("invalid order data: %v", err)
		} else {
			if err := validate.Struct(order); err != nil {
				fmt.Errorf("invalid order data: %v", err)
			} else {
				err = InsertOrder(order)
				if err != nil {
					fmt.Sprintf("Failed to insert order: %v", err)
				} else {
					mas[order.OrderUID] = order
					fmt.Println("Order inserted successfully!")
				}
			}
		}
	}
}
