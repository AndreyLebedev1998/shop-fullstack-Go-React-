-- +goose Up
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    product_name text,
    category_id INT NOT NULL,
    price NUMERIC
);

INSERT INTO products (product_name, price, category_id)
VALUES ('Burn', 250.00, 1);

-- +goose Down
DROP TABLE products;