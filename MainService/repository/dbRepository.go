package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	. "L0_WB/models"

	"github.com/joho/godotenv"
)

var db *sql.DB = nil
var err error

func initDb() {
	if err := godotenv.Load(); err != nil {
		fmt.Print("No .env file found")
	}
	if dbUser, e := os.LookupEnv("DB_USER"); e != nil {
		panic("Can't get DB_USER from .env")
	}
	if dbPassword, e := os.LookupEnv("DB_Password"); e != nil {
		panic("Can't get DB_PASSWORD from .env")
	}
	if dbPort, e := os.LookupEnv("DB_Password"); e != nil {
		panic("Can't get DB_PASSWORD from .env")
	}
	db, err = sql.Open("postgres",
		"postgres://postgres:"+dbUser+"@"+dbPassword+":"+dbPort+"/user_db?sslmode=disable")
}

func InsertOrder(order Order) error {
	if db == nil {
		panic("db = nil")
	}
	if err != nil {
		fmt.Println("Failed to connect to database: ", err)
		return errors.New("Failed to connect to database")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to ping  database: ", err)
		return errors.New("Failed to ping  database")
	}

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid = $1)",
		order.OrderUID).Scan(&exists)

	if exists {
		fmt.Printf("Order %s already exists, skipping\n", order.OrderUID)
		return nil
	}
	order.Payment.Transaction = order.OrderUID

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
			item.TrackNumber,
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

func GetOrder(idStr string) (Order, error) {
	if db == nil {
		panic("db = nil")
	}
	if err != nil {
		fmt.Println("Failed to connect to database: ", err)
		return Order{}, errors.New("Failed to connect to database")
	}
	defer db.Close()
	var order = Order{}
	err = db.QueryRow("SELECT * FROM orders WHERE order_uid = $1", idStr).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSig, &order.CustomerID, &order.DeliveryService,
		&order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard)
	if err != nil {
		return order, errors.New("Cannot get order from db")
	}

	// Информация о доставке
	err = db.QueryRow("SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1", idStr).Scan(
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
		&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Delivery not found for order_uid: %s", idStr)
		} else {
			return order, fmt.Errorf("Cannot get delivery from db: %v", err)
		}
	}

	// Информация о платеже
	err = db.QueryRow("SELECT * FROM payment WHERE transaction = $1", idStr).Scan(
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDT,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal,
		&order.Payment.CustomFee)
	if err != nil && err != sql.ErrNoRows {
		return order, errors.New("Cannot get payment from db")
	}

	// Товары
	rows, err := db.Query("SELECT chrt_id, price, rid, name, sale, size, total_price, nm_id, brand, status  FROM items WHERE order_uid = $1", idStr)
	if err != nil {
		return order, errors.New("Cannot get items from db: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var item Item
		err = rows.Scan(&item.ChrtID,
			&item.Price, &item.RID, &item.Name, &item.Sale, &item.Size,
			&item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return order, errors.New("Cannot scan items from db" + err.Error())
		}
		order.Items = append(order.Items, item)
	}
	fmt.Println(order)
	return order, nil
}

// Функция для восстановления данных из БД
func Get5Orderd(m map[string]Order) error {
	if db == nil {
		panic("db = nil")
	}
	if err != nil {
		fmt.Printf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Запрос для получения 5 order_uid
	query := "SELECT order_uid FROM orders ORDER BY date_created DESC LIMIT 5"

	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	var orderUIDs []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			fmt.Printf("Failed to scan row: %v", err)
		}
		orderUIDs = append(orderUIDs, uid)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Error iterating rows: %v", err)
	}

	if len(orderUIDs) == 0 {
		fmt.Printf("No orders found")
	}
	i := 0
	for u := range orderUIDs {
		o, err := GetOrder(string(u))
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
		m[string(u)] = o
		i++
	}
	return nil
}
