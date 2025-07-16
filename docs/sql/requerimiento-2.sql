DROP TABLE IF EXISTS `warehouses`;
-- Crear tabla warehouses donde se almacenan datos de los almacenes de products
CREATE TABLE if not exists warehouses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    warehouse_code varchar(10) NOT NULL UNIQUE,
    address VARCHAR(100) NOT NULL,
    telephone VARCHAR(15) NOT NULL,
    minimun_capacity FLOAT NOT NULL,
    minimun_temperature FLOAT NOT NULL
);

-- Volcado de datos para la tabla `warehouses`
INSERT INTO warehouses (warehouse_code, address, telephone, minimun_capacity, minimun_temperature) VALUES
    ('WH001', '123 Main St', '123-456-7890', 1000.0, -5.0),
    ('WH002', '456 Elm St', '234-567-8901', 1500.0, -10.0),
    ('WH003', '789 Maple Ave', '345-678-9012', 2000.0, -8.0),
    ('WH004', '135 Oak St', '456-789-0123', 2500.0, -6.0),
    ('WH005', '246 Pine Rd', '567-890-1234', 3000.0, -7.0),
    ('WH006', '357 Cedar Blvd', '678-901-2345', 1200.0, -4.0),
    ('WH007', '468 Spruce Ln', '789-012-3456', 2200.0, -12.0),
    ('WH008', '579 Birch Pl', '890-123-4567', 1800.0, -9.0),
    ('WH009', '681 Aspen Dr', '901-234-5678', 2300.0, -11.0),
    ('WH010', '792 Walnut Ct', '012-345-6789', 2600.0, -6.0);


-- Creacion de carries
drop table if exists carries;
CREATE TABLE if not exists carries (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cid VARCHAR(12) NOT NULL UNIQUE,
    company_name VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    telephone VARCHAR(20) NOT NULL,
    locality_id INT NOT NULL,
   FOREIGN KEY (locality_id) REFERENCES localities(id)
);

INSERT INTO carries (cid, company_name, address, telephone, locality_id) VALUES
    ('C001', 'Company Alpha', '123 Alpha St', '123-456-7890', 1),
    ('C002', 'Beta Solutions', '456 Beta Ave', '234-567-8901', 2),
    ('C003', 'Gamma Enterprises', '789 Gamma Blvd', '345-678-9012', 1),
    ('C004', 'Delta Corp', '321 Delta Rd', '456-789-0123', 2),
    ('C005', 'Epsilon LLC', '654 Epsilon Pl', '567-890-1234', 3);