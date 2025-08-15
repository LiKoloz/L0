CREATE TABLE orders (
    order_uid VARCHAR(50) PRIMARY KEY,
    track_number VARCHAR(50) NOT NULL,
    entry VARCHAR(10) NOT NULL,
    locale VARCHAR(10) NOT NULL,
    internal_signature VARCHAR(100),
    customer_id VARCHAR(50) NOT NULL,
    delivery_service VARCHAR(50) NOT NULL,
    shardkey VARCHAR(10) NOT NULL,
    sm_id INTEGER NOT NULL,
    date_created TIMESTAMPTZ NOT NULL,
    oof_shard VARCHAR(10) NOT NULL
);

CREATE TABLE delivery (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(50) REFERENCES orders(order_uid) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    zip VARCHAR(20) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,
    region VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL
);

CREATE TABLE payment (
    transaction VARCHAR(50) PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    request_id VARCHAR(100),
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(50) NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER NOT NULL
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(50) REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id INTEGER NOT NULL,
    track_number VARCHAR(50) NOT NULL,
    price INTEGER NOT NULL,
    rid VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    sale INTEGER NOT NULL,
    size VARCHAR(10) NOT NULL,
    total_price INTEGER NOT NULL,
    nm_id INTEGER NOT NULL,
    brand VARCHAR(100) NOT NULL,
    status INTEGER NOT NULL
);

-- Выдача прав на конкретную таблицу
GRANT ALL PRIVILEGES ON TABLE orders TO demo_user;

-- Выдача прав на все таблицы в схеме public
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO demo_user;

-- Смена владельца таблицы
ALTER TABLE delivery OWNER TO demo_user;
ALTER TABLE payment OWNER TO demo_user;
ALTER TABLE items OWNER TO demo_user;