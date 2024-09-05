CREATE TABLE IF NOT EXISTS cart_item (
    id BIGSERIAL PRIMARY KEY,
    cart_id BIGINT NOT NULL,
    product VARCHAR(128) NOT NULL,
    quantity BIGINT NOT NULL,
    UNIQUE (cart_id, product),
    FOREIGN KEY (cart_id) REFERENCES cart (id)
);