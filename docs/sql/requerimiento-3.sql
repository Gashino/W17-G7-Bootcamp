CREATE TABLE product_batches (
                                 id INT AUTO_INCREMENT PRIMARY KEY,
                                 batch_number VARCHAR(255) NOT NULL,
                                 current_quantity INT NOT NULL,
                                 current_temperature DECIMAL(10,2),
                                 due_date DATE,
                                 manufacturing_date DATE,
                                 manufacturing_hour TIME,
                                 minimum_temperature DECIMAL(10,2),
                                 initial_quantity INT NOT NULL,
                                 product_id INT NOT NULL,
                                 section_id INT NOT NULL,
                                 FOREIGN KEY (product_id) REFERENCES products(id),
                                 FOREIGN KEY (section_id) REFERENCES sections(id)
);

CREATE TABLE sections (
                          id INT AUTO_INCREMENT PRIMARY KEY,
                          section_number VARCHAR(255) NOT NULL,
                          current_temperature DECIMAL(10,2),
                          current_capacity INT,
                          minimum_temperature DECIMAL(10,2),
                          minimum_capacity INT,
                          product_type_id INT NOT NULL,
                          warehouse_id INT NOT NULL,
                          FOREIGN KEY (product_type_id) REFERENCES products_types(id),
                          FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
);

CREATE INDEX idx_product_batches_product_id ON product_batches(product_id);
CREATE INDEX idx_product_batches_section_id ON product_batches(section_id);
CREATE UNIQUE INDEX idx_product_batches_batch_number ON product_batches(batch_number);
CREATE INDEX idx_product_batches_batch_number_lookup ON product_batches(batch_number);
CREATE INDEX idx_sections_product_type_id ON sections(product_type_id);
CREATE INDEX idx_sections_warehouse_id ON sections(warehouse_id);