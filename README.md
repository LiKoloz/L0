## Архитектура
Состоит из:
1. Сервиса-отправителя данных
2. Оснвного сервиса
3. Кафки
4. Базы данных postresql

Архитектура основного сервиса построенна по паттерну **MVC** 

Также используется паттерт **Репозиторий** для доступа к данным

Для общения используется **Kafka**

Для контейнеризации и оркестарации - **Docker & Docker Compose**

Для запуска проекта достаточно выполнить команду:
```
    docker compose up
```

## Зависимости
Проект содержит 1   зависимость:
1. segmentio/kafka-go - для взаимодействия с kafka

Проект писался с оглядкой на _минимальное_ количество зависимостей

## Особенности
1. Не получилось использовать таблицу от другого пользователя (даже при полном передаче прав на таблицу БД не хочет принимать запрос), поэтому используется супер-пользователь
```go
    	db, err := sql.Open("postgres","postgres://postgres:postgres@postgres:5432/user_db?sslmode=disable")
```
2. Данные из БД собираются постепенно, что сведетельствует о N+1 проблеме, но такой подход был выбран из-за того, чтобы не писать огромный Join и разнести логику
```go
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
```