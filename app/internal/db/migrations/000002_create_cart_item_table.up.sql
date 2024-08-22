CREATE TABLE IF NOT EXISTS cart_item (
    id BIGSERIAL PRIMARY KEY,
    cart_id BIGINT NOT NULL,
    name VARCHAR(128) NOT NULL,
    quantity BIGINT NOT NULL,
    FOREIGN KEY (cart_id) REFERENCES cart (id)
);