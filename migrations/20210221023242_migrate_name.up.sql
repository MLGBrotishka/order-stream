CREATE TABLE IF NOT EXISTS orders(
    order_uid TEXT PRIMARY KEY,
    order_data JSONB
);