-- +goose Up
CREATE TABLE order_item
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id   UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity   INT  NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES product (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE order_item;