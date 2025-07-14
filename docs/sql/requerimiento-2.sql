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
INSERT INTO `warehouses` (`id`, `warehouse_code`, `address`, `telephone`, `minimun_capacity`, `minimun_temperature`) VALUES
(1, 'ASD123','221 Baker Street', '4555666', 100, 50);