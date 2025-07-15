-- =====================================================
-- MERCADO FRESCO DATABASE SETUP WITH SEED DATA
-- =====================================================
-- Script completo para crear la base de datos con todas las tablas
-- organizadas por dependencias y datos de seed iniciales
-- =====================================================

DROP DATABASE IF EXISTS frescos;
CREATE DATABASE frescos;
USE frescos;

-- Configuración inicial
SET FOREIGN_KEY_CHECKS = 0;

-- =====================================================
-- PASO 1: TABLAS BASE (SIN DEPENDENCIAS)
-- =====================================================

-- Tabla localities (base)
DROP TABLE IF EXISTS localities;
CREATE TABLE localities (
    id INT AUTO_INCREMENT PRIMARY KEY,
    locality_name VARCHAR(255) NOT NULL UNIQUE,
    province_name VARCHAR(255) NOT NULL,
    country_name VARCHAR(255) NOT NULL
);

-- Tabla warehouses (base)
DROP TABLE IF EXISTS warehouses;
CREATE TABLE warehouses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    warehouse_code VARCHAR(10) NOT NULL UNIQUE,
    address VARCHAR(100) NOT NULL,
    telephone VARCHAR(15) NOT NULL,
    minimum_capacity FLOAT NOT NULL,
    minimum_temperature FLOAT NOT NULL
);

-- Tabla buyers (base)
DROP TABLE IF EXISTS buyers;
CREATE TABLE buyers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    card_number_id VARCHAR(50) NOT NULL UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Tabla product_types (base)
DROP TABLE IF EXISTS product_types;
CREATE TABLE product_types (
    id INT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- =====================================================
-- PASO 2: TABLAS CON DEPENDENCIAS DE NIVEL 1
-- =====================================================

-- Tabla carries (depende de localities)
DROP TABLE IF EXISTS carries;
CREATE TABLE if not exists carries (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cid VARCHAR(12) NOT NULL UNIQUE,
    company_name VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    telephone VARCHAR(20) NOT NULL,
    locality_id INT NOT NULL,
   FOREIGN KEY (locality_id) REFERENCES localities(id)
);

-- Tabla sellers (depende de localities)
DROP TABLE IF EXISTS sellers;
CREATE TABLE sellers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cid VARCHAR(255) NOT NULL UNIQUE,
    company_name VARCHAR(255) NOT NULL,
    address VARCHAR(500) NOT NULL,
    telephone VARCHAR(50) NOT NULL,
    locality_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (locality_id) REFERENCES localities(id)
);

-- Tabla employees (depende de warehouses)
DROP TABLE IF EXISTS employees;
CREATE TABLE employees (
    id INT AUTO_INCREMENT PRIMARY KEY,
    card_number_id VARCHAR(64) NOT NULL UNIQUE,
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64) NOT NULL,
    warehouse_id INT NOT NULL,
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
);

-- =====================================================
-- PASO 3: TABLAS CON DEPENDENCIAS DE NIVEL 2
-- =====================================================

-- Tabla products (depende de product_types y sellers)
DROP TABLE IF EXISTS products;
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

-- Tabla sections (depende de product_types y warehouses)
DROP TABLE IF EXISTS sections;
CREATE TABLE sections (
    id INT AUTO_INCREMENT PRIMARY KEY,
    section_number VARCHAR(255) NOT NULL,
    current_temperature DECIMAL(10,2),
    current_capacity INT,
    minimum_temperature DECIMAL(10,2),
    minimum_capacity INT,
    product_type_id INT NOT NULL,
    warehouse_id INT NOT NULL,
    FOREIGN KEY (product_type_id) REFERENCES product_types(id),
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
);

-- =====================================================
-- PASO 4: TABLAS CON DEPENDENCIAS DE NIVEL 3
-- =====================================================

-- Tabla product_records (depende de products)
DROP TABLE IF EXISTS product_records;
CREATE TABLE product_records (
    id INT PRIMARY KEY AUTO_INCREMENT,
    last_update_date DATETIME NOT NULL,
    purchase_price DOUBLE NOT NULL,
    sale_price DOUBLE NOT NULL,
    product_id INT NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- Tabla product_batches (depende de products y sections)
DROP TABLE IF EXISTS product_batches;
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

-- =====================================================
-- PASO 5: TABLAS CON DEPENDENCIAS DE NIVEL 4
-- =====================================================

-- Tabla inbound_orders (depende de employees, product_batches, y warehouses)
DROP TABLE IF EXISTS inbound_orders;
CREATE TABLE inbound_orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_date DATE NOT NULL,
    order_number VARCHAR(64) NOT NULL,
    employee_id INT NOT NULL,
    product_batch_id INT NOT NULL,
    warehouse_id INT NOT NULL,
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (product_batch_id) REFERENCES product_batches(id),
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
);

-- Tabla purchase_orders (depende de buyers y product_records)
DROP TABLE IF EXISTS purchase_orders;
CREATE TABLE purchase_orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_number VARCHAR(50) NOT NULL UNIQUE,
    order_date DATE NOT NULL,
    tracking_code VARCHAR(100) NOT NULL,
    buyer_id INT NOT NULL,
    product_record_id INT NOT NULL,
    FOREIGN KEY (buyer_id) REFERENCES buyers(id) ON DELETE CASCADE,
    FOREIGN KEY (product_record_id) REFERENCES product_records(id)
);

-- =====================================================
-- ÍNDICES PARA MEJORAR PERFORMANCE
-- =====================================================

-- Índices para buyers
CREATE INDEX idx_buyers_card_number_id ON buyers(card_number_id);

-- Índices para product_batches
CREATE INDEX idx_product_batches_product_id ON product_batches(product_id);
CREATE INDEX idx_product_batches_section_id ON product_batches(section_id);
CREATE UNIQUE INDEX idx_product_batches_batch_number ON product_batches(batch_number);

-- Índices para sections
CREATE INDEX idx_sections_product_type_id ON sections(product_type_id);
CREATE INDEX idx_sections_warehouse_id ON sections(warehouse_id);

-- Índices para purchase_orders
CREATE INDEX idx_purchase_orders_order_number ON purchase_orders(order_number);
CREATE INDEX idx_purchase_orders_buyer_id ON purchase_orders(buyer_id);
CREATE INDEX idx_purchase_orders_product_record_id ON purchase_orders(product_record_id);

-- =====================================================
-- DATOS DE SEED - TABLAS BASE
-- =====================================================

-- Seed para localities
INSERT INTO localities (locality_name, province_name, country_name) VALUES
('Ciudad Autónoma de Buenos Aires', 'CABA', 'Argentina'),
('Córdoba', 'Córdoba', 'Argentina'),
('Rosario', 'Santa Fe', 'Argentina'),
('Mendoza', 'Mendoza', 'Argentina'),
('La Plata', 'Buenos Aires', 'Argentina'),
('San Miguel de Tucumán', 'Tucumán', 'Argentina'),
('Mar del Plata', 'Buenos Aires', 'Argentina'),
('Salta', 'Salta', 'Argentina'),
('Santa Fe', 'Santa Fe', 'Argentina'),
('San Juan', 'San Juan', 'Argentina');

-- Seed para warehouses (primeros 10 del archivo JSON)
INSERT INTO warehouses (id, warehouse_code, address, telephone, minimum_capacity, minimum_temperature) VALUES
(1, 'WH001', '7 Calle Principal, Ciudad 1, País', '+1234567001', 3694, -4),
(2, 'WH002', '14 Calle Principal, Ciudad 2, País', '+1234567002', 4930, -24),
(3, 'WH003', '21 Calle Principal, Ciudad 3, País', '+1234567003', 2722, -14),
(4, 'WH004', '28 Calle Principal, Ciudad 4, País', '+1234567004', 2963, 0),
(5, 'WH005', '35 Calle Principal, Ciudad 5, País', '+1234567005', 3487, -9),
(6, 'WH006', '42 Calle Principal, Ciudad 6, País', '+1234567006', 3965, -14),
(7, 'WH007', '49 Calle Principal, Ciudad 7, País', '+1234567007', 4234, -7),
(8, 'WH008', '56 Calle Principal, Ciudad 8, País', '+1234567008', 3876, -12),
(9, 'WH009', '63 Calle Principal, Ciudad 9, País', '+1234567009', 4567, -18),
(10, 'WH010', '70 Calle Principal, Ciudad 10, País', '+1234567010', 3234, -6);

-- Seed para product_types
INSERT INTO product_types (id, name) VALUES
(101, 'Electronics'),
(102, 'Clothing'),
(103, 'Home & Kitchen'),
(104, 'Beauty & Personal Care'),
(105, 'Sports & Outdoors'),
(106, 'Toys & Games'),
(107, 'Books'),
(108, 'Health & Household'),
(109, 'Grocery'),
(110, 'Automotive'),
(111, 'Tools & Home Improvement'),
(112, 'Pet Supplies'),
(113, 'Office Products'),
(114, 'Baby Products'),
(115, 'Jewelry'),
(116, 'Musical Instruments'),
(117, 'Industrial & Scientific'),
(118, 'Arts & Crafts'),
(119, 'Luggage & Travel Gear'),
(120, 'Software');

-- Seed para buyers (primeros 20 del archivo JSON)
INSERT INTO buyers (id, card_number_id, first_name, last_name) VALUES
(1, '402323', 'John', 'Doe'),
(2, '402324', 'Jane', 'Smith'),
(3, '402325', 'Bob', 'Johnson'),
(4, '402326', 'Buyer4', 'Last4'),
(5, '402327', 'Buyer5', 'Last5'),
(6, '402328', 'Buyer6', 'Last6'),
(7, '402329', 'Buyer7', 'Last7'),
(8, '402330', 'Buyer8', 'Last8'),
(9, '402331', 'Buyer9', 'Last9'),
(10, '402332', 'Buyer10', 'Last10'),
(11, '402333', 'Buyer11', 'Last11'),
(12, '402334', 'Buyer12', 'Last12'),
(13, '402335', 'Buyer13', 'Last13'),
(14, '402336', 'Buyer14', 'Last14'),
(15, '402337', 'Buyer15', 'Last15'),
(16, '402338', 'Buyer16', 'Last16'),
(17, '402339', 'Buyer17', 'Last17'),
(18, '402340', 'Buyer18', 'Last18'),
(19, '402341', 'Buyer19', 'Last19'),
(20, '402342', 'Buyer20', 'Last20');

-- =====================================================
-- DATOS DE SEED - TABLAS CON DEPENDENCIAS NIVEL 1
-- =====================================================

-- Seed para carries
INSERT INTO carries (cid, company_name, address, telephone, locality_id) VALUES
    ('C001', 'Company Alpha', '123 Alpha St', '123-456-7890', 1),
    ('C002', 'Beta Solutions', '456 Beta Ave', '234-567-8901', 2),
    ('C003', 'Gamma Enterprises', '789 Gamma Blvd', '345-678-9012', 1),
    ('C004', 'Delta Corp', '321 Delta Rd', '456-789-0123', 2),
    ('C005', 'Epsilon LLC', '654 Epsilon Pl', '567-890-1234', 3);

-- Seed para sellers (del archivo JSON)
INSERT INTO sellers (id, cid, company_name, address, telephone, locality_id) VALUES
(1, '20304050601', 'Alkemy', 'Monroe 860', '47470000', 1),
(2, '20304050602', 'Globant S.A.', 'Av. Córdoba 1111', '48000001', 1),
(3, '20304050603', 'Mercado Libre', 'Av. Caseros 3039', '47470002', 1),
(4, '20304050604', 'Despegar.com', 'Av. Libertador 7208', '47470003', 1),
(5, '20304050605', 'PedidosYa', 'Av. del Libertador 400', '47470004', 1),
(6, '20304050606', 'Rappi', 'Calle Falsa 123', '47470005', 2),
(7, '20304050607', 'Naranja X', 'Av. Colón 144', '47470006', 2),
(8, '20304050608', 'Ualá', 'San Martín 500', '47470007', 3),
(9, '20304050609', 'Tiendanube', 'Av. Belgrano 675', '47470008', 3),
(10, '20304050610', 'Nubank', 'Diagonal Norte 350', '47470009', 4);

-- Seed para employees (primeros 20 del archivo JSON)
INSERT INTO employees (id, card_number_id, first_name, last_name, warehouse_id) VALUES
(1, '12345678', 'Juan', 'Pérez', 1),
(2, '23456789', 'María', 'González', 2),
(3, '34567890', 'Carlos', 'López', 1),
(4, '45678901', 'Ana', 'Martínez', 3),
(5, '56789012', 'Luis', 'Rodríguez', 2),
(6, '67890123', 'Laura', 'Hernández', 1),
(7, '78901234', 'Diego', 'García', 4),
(8, '89012345', 'Carmen', 'Ruiz', 3),
(9, '90123456', 'Miguel', 'Jiménez', 2),
(10, '01234567', 'Isabel', 'Moreno', 1),
(11, '12340678', 'Roberto', 'Álvarez', 5),
(12, '23450789', 'Elena', 'Muñoz', 3),
(13, '34560890', 'Francisco', 'Romero', 2),
(14, '45670901', 'Pilar', 'Sanz', 4),
(15, '56780912', 'Javier', 'Torres', 1),
(16, '67890234', 'Rosa', 'Flores', 3),
(17, '78901345', 'Antonio', 'Vega', 2),
(18, '89012456', 'Cristina', 'Ramos', 5),
(19, '90123567', 'Fernando', 'Castro', 1),
(20, '01234678', 'Mercedes', 'Ortega', 4);

-- =====================================================
-- DATOS DE SEED - TABLAS CON DEPENDENCIAS NIVEL 2
-- =====================================================

-- Seed para products (primeros 18 del archivo JSON)
INSERT INTO products (id, product_code, description, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, width, height, length, product_type_id, seller_id) VALUES
(1, 'APL123', 'Fresh organic apples', 1.2, 5, -1.0, 3, 3.0, 3.0, 4.0, 101, 1),
(2, 'BRD456', 'Whole grain bread', 0.8, 7, -2.0, 2, 4.0, 5.0, 10.0, 102, 2),
(3, 'MLK789', 'Almond milk, unsweetened', 0.5, 10, -3.0, 1, 8.0, 20.0, 8.0, 103, 3),
(4, 'EGG101', 'Free-range chicken eggs', 0.6, 12, 0.0, 4, 3.5, 2.0, 5.0, 104, 4),
(5, 'SPN202', 'Organic spinach leaves', 0.2, 3, -1.5, 2, 5.0, 10.0, 6.0, 105, 5),
(6, 'YGT303', 'Greek yogurt, plain', 0.9, 15, -4.0, 1, 7.0, 12.0, 7.0, 106, 6),
(7, 'CHC404', 'Dark chocolate bar', 0.25, 18, -5.0, 0, 5.0, 1.0, 15.0, 107, 7),
(8, 'PST505', 'Whole wheat pasta', 0.4, 24, -6.0, 0, 7.0, 8.0, 20.0, 108, 8),
(9, 'CHS606', 'Cheddar cheese block', 0.7, 30, -2.5, 5, 5.0, 6.0, 10.0, 109, 9),
(10, 'BRY707', 'Frozen mixed berries', 1.0, 60, -18.0, 6, 10.0, 15.0, 10.0, 110, 10),
(11, 'QNA808', 'Quinoa grain', 0.45, 36, -5.0, 0, 8.0, 10.0, 8.0, 111, 1),
(12, 'HNY909', 'Organic honey', 0.35, 720, -10.0, 0, 8.0, 14.0, 8.0, 112, 2),
(13, 'PBT010', 'Peanut butter, smooth', 0.5, 180, -12.0, 0, 8.0, 12.0, 8.0, 113, 3),
(14, 'BNS111', 'Canned black beans', 0.4, 365, -8.0, 0, 7.0, 11.0, 7.0, 114, 4),
(15, 'JAM212', 'Raspberry jam', 0.3, 540, -9.0, 0, 7.0, 10.0, 7.0, 115, 5),
(16, 'SLS313', 'Spicy salsa sauce', 0.25, 120, -7.0, 0, 6.0, 9.0, 6.0, 116, 6),
(17, 'DRT414', 'Dried apricots', 0.2, 720, -15.0, 0, 5.0, 8.0, 5.0, 117, 7),
(18, 'CWT515', 'Coconut water', 0.5, 90, -11.0, 0, 5.0, 12.0, 5.0, 118, 8);

-- Seed para sections (del archivo JSON)
INSERT INTO sections (id, section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, product_type_id, warehouse_id) VALUES
(1, '1', 5.00, 2.00, 50, 10, 103, 1),
(2, '2', -18.00, -20.00, 75, 20, 105, 1),
(3, '3', 25.00, 15.00, 30, 5, 102, 2),
(4, '4', 10.00, 5.00, 120, 30, 107, 2),
(5, '5', -5.00, -10.00, 90, 15, 101, 3);

-- =====================================================
-- DATOS DE SEED - TABLAS CON DEPENDENCIAS NIVEL 3
-- =====================================================

-- Seed para product_records
INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id) VALUES
('2023-01-15 10:00:00', 15.50, 23.25, 1),
('2023-01-16 11:00:00', 8.00, 12.00, 2),
('2023-01-17 12:00:00', 22.75, 34.00, 3),
('2023-01-18 13:00:00', 18.25, 27.50, 4),
('2023-01-19 14:00:00', 12.00, 18.00, 5),
('2023-01-20 15:00:00', 25.50, 38.25, 6),
('2023-01-21 16:00:00', 35.00, 52.50, 7),
('2023-01-22 17:00:00', 14.75, 22.00, 8),
('2023-01-23 18:00:00', 28.00, 42.00, 9),
('2023-01-24 19:00:00', 45.50, 68.25, 10),
('2023-01-25 20:00:00', 19.25, 28.75, 11),
('2023-01-26 21:00:00', 42.00, 63.00, 12),
('2023-01-27 22:00:00', 24.50, 36.75, 13),
('2023-01-28 23:00:00', 16.75, 25.00, 14),
('2023-01-29 09:00:00', 31.25, 46.75, 15);

-- Seed para product_batches
INSERT INTO product_batches (batch_number, current_quantity, current_temperature, due_date, manufacturing_date, manufacturing_hour, minimum_temperature, initial_quantity, product_id, section_id) VALUES
('BATCH001', 100, 5.0, '2023-06-15', '2023-01-15', '08:00:00', 2.0, 150, 1, 1),
('BATCH002', 75, -18.0, '2023-12-31', '2023-01-16', '09:00:00', -20.0, 100, 2, 2),
('BATCH003', 50, 25.0, '2023-05-20', '2023-01-17', '10:00:00', 15.0, 75, 3, 3),
('BATCH004', 125, 10.0, '2023-07-10', '2023-01-18', '11:00:00', 5.0, 200, 4, 4),
('BATCH005', 80, -5.0, '2024-01-30', '2023-01-19', '12:00:00', -10.0, 120, 5, 5),
('BATCH006', 60, 5.0, '2023-08-15', '2023-01-20', '13:00:00', 2.0, 90, 6, 1),
('BATCH007', 90, -18.0, '2023-11-25', '2023-01-21', '14:00:00', -20.0, 110, 7, 2),
('BATCH008', 40, 25.0, '2023-09-30', '2023-01-22', '15:00:00', 15.0, 60, 8, 3),
('BATCH009', 110, 10.0, '2023-10-05', '2023-01-23', '16:00:00', 5.0, 150, 9, 4),
('BATCH010', 70, -5.0, '2024-03-15', '2023-01-24', '17:00:00', -10.0, 100, 10, 5);

-- =====================================================
-- DATOS DE SEED - TABLAS CON DEPENDENCIAS NIVEL 4
-- =====================================================

-- Seed para inbound_orders
INSERT INTO inbound_orders (order_date, order_number, employee_id, product_batch_id, warehouse_id) VALUES
('2023-02-01', 'IO001', 1, 1, 1),
('2023-02-02', 'IO002', 2, 2, 2),
('2023-02-03', 'IO003', 3, 3, 3),
('2023-02-04', 'IO004', 4, 4, 2),
('2023-02-05', 'IO005', 5, 5, 3),
('2023-02-06', 'IO006', 6, 6, 1),
('2023-02-07', 'IO007', 7, 7, 2),
('2023-02-08', 'IO008', 8, 8, 3),
('2023-02-09', 'IO009', 9, 9, 2),
('2023-02-10', 'IO010', 10, 10, 3);

-- Seed para purchase_orders
INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id, product_record_id) VALUES
('PO001', '2023-01-15', 'TRK001ABC', 1, 1),
('PO002', '2023-01-16', 'TRK002DEF', 2, 2),
('PO003', '2023-01-17', 'TRK003GHI', 3, 3),
('PO004', '2023-01-18', 'TRK004JKL', 4, 4),
('PO005', '2023-01-19', 'TRK005MNO', 5, 5),
('PO006', '2023-01-20', 'TRK006PQR', 6, 6),
('PO007', '2023-01-21', 'TRK007STU', 7, 7),
('PO008', '2023-01-22', 'TRK008VWX', 8, 8),
('PO009', '2023-01-23', 'TRK009YZA', 9, 9),
('PO010', '2023-01-24', 'TRK010BCD', 10, 10);

-- =====================================================
-- CONFIGURACIÓN FINAL
-- =====================================================

-- Reactivar verificación de foreign keys
SET FOREIGN_KEY_CHECKS = 1;

-- Ajustar AUTO_INCREMENT para evitar conflictos
ALTER TABLE buyers AUTO_INCREMENT = 21;
ALTER TABLE employees AUTO_INCREMENT = 21;
ALTER TABLE products AUTO_INCREMENT = 19;
ALTER TABLE sections AUTO_INCREMENT = 6;
ALTER TABLE product_records AUTO_INCREMENT = 16;
ALTER TABLE product_batches AUTO_INCREMENT = 11;
ALTER TABLE inbound_orders AUTO_INCREMENT = 11;
ALTER TABLE purchase_orders AUTO_INCREMENT = 11;

-- =====================================================
-- SCRIPT COMPLETADO EXITOSAMENTE
-- =====================================================
-- Todas las tablas han sido creadas respetando dependencias
-- Todos los datos de seed han sido insertados
-- Todos los índices han sido creados
-- ===================================================== 