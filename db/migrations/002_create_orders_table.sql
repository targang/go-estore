-- +goose Up
CREATE TABLE orders
(
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_name  VARCHAR(255) NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    status         INT          NOT NULL,
    created_at     TIMESTAMPTZ      DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMPTZ      DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_order_timestamp() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE OR REPLACE TRIGGER trigger_update_order_timestamp
    BEFORE UPDATE
    ON orders
    FOR EACH ROW
    WHEN (OLD.status IS DISTINCT FROM NEW.status)
EXECUTE FUNCTION update_order_timestamp();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_update_order_timestamp ON orders;
DROP FUNCTION IF EXISTS update_order_timestamp();
DROP TABLE orders;
