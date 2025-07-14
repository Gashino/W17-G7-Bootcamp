CREATE TABLE product_types (
       id INT PRIMARY KEY,
       name VARCHAR(255) NOT NULL
);

CREATE TABLE products (
      id INT PRIMARY KEY AUTO_INCREMENT,
      product_code VARCHAR(50) UNIQUE,
      description TEXT,
      net_weight DOUBLE,
      expiration_rate INT,
      recommended_freezing_temperature DOUBLE,
      freezing_rate INT,
      product_type_id INT,
      seller_id INT,
      width DOUBLE,
      height DOUBLE,
      length DOUBLE,
      FOREIGN KEY (product_type_id) REFERENCES product_types(id),
      FOREIGN KEY (seller_id) REFERENCES sellers(id)
);
