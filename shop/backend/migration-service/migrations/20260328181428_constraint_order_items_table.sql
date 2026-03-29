-- +goose Up
ALTER TABLE order_items
ADD CONSTRAINT fk_order_items_orders FOREIGN KEY (order_id) REFERENCES orders(id);

ALTER TABLE order_items
ADD CONSTRAINT fk_order_items_products FOREIGN KEY (product_id) REFERENCES products(id);

-- +goose Down
ALTER TABLE order_items
DROP CONSTRAINT fk_order_items_orders;

ALTER TABLE order_items
DROP CONSTRAINT fk_order_items_products;
