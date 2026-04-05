-- +goose Up
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    product_name text,
    category_id INT NOT NULL,
    price NUMERIC,
    image_url TEXT
);

INSERT INTO products (product_name, price, category_id)
VALUES
('Product A', 25.25, 1),
('Product B', 50.00, 1),
('Product C', 50.00, 1),
('Product D', 100.00, 1);

-- +goose Down
DROP TABLE products;