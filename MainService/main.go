package main

import (
	. "L0_WB/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

var mas [5]Order
var i int = 0

func main() {
	ch := make(chan struct{})
	go func() {

		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"kafka:9092"},
			Topic:   "test-topic",
			GroupID: "my-groupID",
		})
		defer reader.Close()
		for {
			var order Order
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				fmt.Println("Ошибка при получении:", err)
				panic("Ошибка при получении")
			}

			fmt.Println("Получили: ", string(msg.Value))

			err = json.Unmarshal(msg.Value, &order)

			db, err := sql.Open("postgres",
				"postgres://postgres:postgres@postgres:5432/user_db?sslmode=disable")

			if err != nil {
				panic(fmt.Sprintf("Failed to connect to database: %v", err))
			}
			defer db.Close()

			// Проверка соединения
			err = db.Ping()
			if err != nil {
				panic(fmt.Sprintf("Failed to ping database: %v", err))
			}

			// Вставка заказа
			err = insertOrder(db, order)
			if err != nil {
				panic(fmt.Sprintf("Failed to insert order: %v", err))
			}
			if i == 5 {
				i = 0
			}
			mas[i] = order
			i++
			fmt.Println("Order inserted successfully!")
		}

	}()
	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		
	})
}

func insertOrder(db *sql.DB, order Order) error {

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid = $1)",
		order.OrderUID).Scan(&exists)

	if exists {
		fmt.Printf("Order %s already exists, skipping\n", order.OrderUID)
		return nil
	}
	order.Payment.Transaction = order.OrderUID

	// Начало транзакции
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(`
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature, customer_id,
			delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSig,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		return fmt.Errorf("insert into orders failed: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO delivery (
			order_uid, name, phone, zip, city, address, region, email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)
	if err != nil {
		return fmt.Errorf("insert into delivery failed: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO payment (
			transaction, request_id, currency, provider, amount, payment_dt, bank,
			delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDT,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err != nil {
		return fmt.Errorf("insert into payment failed: %w", err)
	}

	for _, item := range order.Items {
		_, err = tx.Exec(`
			INSERT INTO items (
				order_uid, chrt_id, track_number, price, rid, name, sale, size, 
				total_price, nm_id, brand, status
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			order.OrderUID,
			item.ChrtID,
			item.TrackNum,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			return fmt.Errorf("insert into items failed: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	fmt.Println("Успешная вставка!")
	return nil
}
