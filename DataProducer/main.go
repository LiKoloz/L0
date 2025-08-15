package main

import (
	. "DataProducer/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	ctx := context.Background()

	order1 := Order{
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDT:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []Item{
			{
				ChrtID:     9934930,
				TrackNum:   "WBILMTESTTRACK",
				Price:      453,
				RID:        "ab4219087a764ae0btest",
				Name:       "Mascaras",
				Sale:       30,
				Size:       "0",
				TotalPrice: 317,
				NmID:       2389212,
				Brand:      "Vivienne Sabo",
				Status:     202,
			},
		},
		Locale:          "en",
		InternalSig:     "",
		CustomerID:      "test",
		DeliveryService: "meest",
		ShardKey:        "9",
		SmID:            99,
		DateCreated:     time.Date(2021, 11, 26, 6, 22, 19, 0, time.UTC),
		OofShard:        "1",
	}

	// Заказ 2 (с двумя товарами)
	order2 := Order{
		OrderUID:    "order2_uid_12345",
		TrackNumber: "TRACK-789-XYZ",
		Entry:       "ENTRY2",
		Delivery: Delivery{
			Name:    "Иван Иванов",
			Phone:   "+79161234567",
			Zip:     "127001",
			City:    "Москва",
			Address: "ул. Тверская, д. 1",
			Region:  "Москва",
			Email:   "ivan@example.com",
		},
		Payment: Payment{
			Transaction:  "trans_78910",
			Currency:     "RUB",
			Provider:     "sberpay",
			Amount:       4500,
			PaymentDT:    1676543210,
			Bank:         "sber",
			DeliveryCost: 500,
			GoodsTotal:   4000,
			CustomFee:    0,
		},
		Items: []Item{
			{
				ChrtID:     111111,
				TrackNum:   "ITEM-111",
				Price:      2000,
				Name:       "Ноутбук",
				Sale:       0,
				Size:       "",
				TotalPrice: 2000,
				Brand:      "BrandX",
				Status:     200,
			},
			{
				ChrtID:     222222,
				TrackNum:   "ITEM-222",
				Price:      2500,
				Name:       "Смартфон",
				Sale:       20,
				Size:       "",
				TotalPrice: 2000,
				Brand:      "BrandY",
				Status:     200,
			},
		},
		Locale:          "ru",
		CustomerID:      "customer789",
		DeliveryService: "cdelivery",
		ShardKey:        "5",
		SmID:            88,
		DateCreated:     time.Now().UTC(),
		OofShard:        "2",
	}

	// Заказ 3 (минималистичный)
	order3 := Order{
		OrderUID:    "min_order_333",
		TrackNumber: "MIN-TRACK-001",
		Entry:       "MINENTRY",
		Delivery: Delivery{
			Name:    "Анна Сидорова",
			Phone:   "+380991112233",
			City:    "Киев",
			Address: "пр. Победы, 10",
		},
		Payment: Payment{
			Transaction: "trans_min_333",
			Currency:    "EUR",
			Provider:    "paypal",
			Amount:      99,
		},
		Items: []Item{
			{
				ChrtID: 333333,
				Name:   "Книга",
				Price:  15,
				Brand:  "Издательство",
				Status: 200,
			},
		},
		DateCreated: time.Date(2023, 5, 15, 12, 30, 0, 0, time.UTC),
	}

	orders := []Order{order1, order2, order3}
	time.Sleep(10 * time.Second)
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "test-topic",
	})
	defer writer.Close()
	for _, v := range orders {
		js, err := json.Marshal(v)
		if err != nil {
			panic("Ошибка при сериализации")
		}
		err = writer.WriteMessages(ctx, kafka.Message{
			Value: js,
		})
		if err != nil {
			fmt.Println("Ошибка: ", err)
			panic("Ошибка при отправке")
		}
	}

}
