package services

import (
	. "L0_WB/models"
	. "L0_WB/repository"
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func GetDataFromKafka(mas *[5]Order, i *int) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "test-topic",
		GroupID: "my-groupID",
	})
	defer reader.Close()
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

		err = json.Unmarshal(msg.Value, &order)
		if err := order.Validate(); err != nil {
			fmt.Errorf("invalid order data: %v", err)
		} else {
			err = InsertOrder(order)
			if err != nil {
				fmt.Sprintf("Failed to insert order: %v", err)
			} else {
				if *i == 5 {
					*i = 0
				}
				mas[*i] = order
				*i++
				fmt.Println("Order inserted successfully!")
			}
		}
	}
}
