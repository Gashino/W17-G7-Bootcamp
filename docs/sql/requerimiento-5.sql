CREATE TABLE employees (
    id INT AUTO_INCREMENT PRIMARY KEY,
    card_number_id VARCHAR(64) NOT NULL UNIQUE,
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64) NOT NULL,
    warehouse_id INT NOT NULL,
    -- FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
);



CREATE TABLE inbound_orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_date DATE NOT NULL,
    order_number VARCHAR(64) NOT NULL,
    employee_id INT NOT NULL,
    product_batch_id INT NOT NULL,
    warehouse_id INT NOT NULL,
    -- FOREIGN KEY (employee_id) REFERENCES employees(id),
    -- FOREIGN KEY (product_batch_id) REFERENCES product_batches(id),
    -- FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
);