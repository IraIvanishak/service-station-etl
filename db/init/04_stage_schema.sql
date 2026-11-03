CREATE SCHEMA IF NOT EXISTS stage;

-- Таблиця inventory
CREATE TABLE stage.inventory (
    inventory_id SERIAL PRIMARY KEY,
    detail_id INT,
    quantity INT,
    service_station_id INT
);

-- Таблиця service_inventory
CREATE TABLE stage.service_inventory (
    service_inventory_id SERIAL PRIMARY KEY,
    service_repair_id INT,
    inventory_id INT,
    amount INT
);

-- Таблиця service_repair
CREATE TABLE stage.service_repair (
    service_repairs_id SERIAL PRIMARY KEY,
    service_id INT,
    repair_id INT,
    quantity INT,
    completion_date DATE,
    service_total_price NUMERIC
);

-- Таблиця supplier
CREATE TABLE stage.supplier (
    supplier_id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    city_id INT,
    phone_number VARCHAR(255),
    address VARCHAR(255)
);

-- Таблиця supply_details
CREATE TABLE stage.supply_details (
    supply_details_id SERIAL PRIMARY KEY,
    detail_id INT,
    quantity INT,
    price_per_unit NUMERIC,
    supply_id INT
);

-- Таблиця detail
CREATE TABLE stage.detail (
    detail_id SERIAL PRIMARY KEY,
    description TEXT,
    manufacturer_id INT,
    detail_type_id INT,
    article VARCHAR(255),
    price NUMERIC,
    dimensions VARCHAR(255),
    weight NUMERIC
);

-- Таблиця repair
CREATE TABLE stage.repair (
    repair_id SERIAL PRIMARY KEY,
    employee_id INT,
    start_date DATE,
    repair_status_id INT,
    notes TEXT,
    repair_number VARCHAR(255),
    payment_status_id INT,
    car_client_id INT,
    total_sum NUMERIC
);

-- Таблиця service
CREATE TABLE stage.service (
    service_id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    price NUMERIC,
    estimated_duration INTERVAL
);

CREATE TABLE stage.supply(
    supply_id SERIAL PRIMARY KEY,
    supply_number VARCHAR(255),
    supply_date DATE,
    total_sum NUMERIC,
    status_id INT,
    supplier_id INT,
    employee_id INT
);

