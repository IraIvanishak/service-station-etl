-- name: ListServiceStations :many
SELECT * FROM service_station
JOIN city ON service_station.city_id = city.city_id
JOIN country ON city.country_id = country.country_id;


-- name: AddServiceStation :exec
INSERT INTO service_station (address, city_id) VALUES (
@address, @city_id
);

-- name: UpdateServiceStation :exec
UPDATE service_station SET address = @address, city_id = @city_id WHERE station_id = @station_id;

-- name: DeleteServiceStation :execrows
DELETE FROM service_station WHERE station_id = @station_id;


-- name: ListCities :many
SELECT * FROM city;


-- name: ListEmployees :many
SELECT 
    emp.employee_id,
    emp.first_name,
    emp.last_name,
    er.name AS role,
    ss.address AS service_station,
    ss.station_id AS service_station_id
FROM
    public.employee emp
LEFT JOIN public.employee_role er ON emp.role_id = er.role_id
LEFT JOIN public.service_station ss ON emp.service_station_id = ss.station_id

ORDER BY
    emp.employee_id;


-- name: ListClients :many
SELECT * FROM client;



-- name: ListCars :many
SELECT 
    car.car_id,
    car.model,
    car.year,
    car.body_type,
    car.brand_id,
    br.name AS brand_name
FROM
    public.car car
LEFT JOIN public.brand br ON car.brand_id = br.brand_id
ORDER BY
    car.car_id;

-- name: StoreOrder :one
INSERT INTO repair_order (employee_id, start_date, repair_status_id, notes, payment_status_id, car_client_id, total_sum) VALUES (
    @employee_id, @start_date, @repair_status_id, @notes, @payment_status_id, @car_client_id, @total_sum
) RETURNING repair_id;

-- name: SetRepairNumber :exec
UPDATE repair_order SET repair_number = @repair_number WHERE repair_id = @repair_id;

-- name: StoreCarClient :one
INSERT INTO car_client (car_id, client_id, description, registration_number, vin, color) VALUES (
    @car_id, @client_id, @description, @registration_number, @vin, @color
) RETURNING car_client_id;


-- name: DeleteRepairParts :exec
DELETE FROM service_inventory
WHERE service_repair_id IN (SELECT service_repairs_id FROM service_repairs WHERE repair_id = @repair_id);

-- name: DeleteRepairServices :exec
DELETE FROM service_repairs WHERE repair_id = @repair_id;

-- name: DeleteOrder :execrows
DELETE FROM repair_order WHERE repair_id = @repair_id;

-- name: StoreServiceRepair :one
-- Priced in numeric from the service's catalog price; no row if the service does not exist.
INSERT INTO service_repairs (service_id, repair_id, quantity, completion_date, service_total_price)
SELECT s.service_id, sqlc.arg(repair_id)::int4, sqlc.arg(quantity)::int4, sqlc.arg(completion_date)::date,
       s.price * sqlc.arg(quantity)::int4
FROM service s
WHERE s.service_id = sqlc.arg(service_id)::int4
RETURNING service_repairs_id;

-- name: ListDetails :many 
SELECT * FROM detail;



-- name: ListServices :many
SELECT * FROM service;


-- name: ListSuppliers :many
SELECT * FROM supplier;


-- name: MetaData :many
SELECT * from meta.dataloadhistory;

-- name: ListCarClients :many
SELECT 
    cc.car_client_id,
    cc.car_id,
    cc.client_id,
    cc.description,
    cc.registration_number,
    cc.vin,
    cc.color,
    c.model,
    c.year,
    c.body_type,
    c.brand_id,
    b.name AS brand_name
FROM
    public.car_client cc
LEFT JOIN public.car c ON cc.car_id = c.car_id
LEFT JOIN public.brand b ON c.brand_id = b.brand_id
ORDER BY
    cc.car_client_id;


-- name: CreateClient :exec
INSERT INTO client (first_name, last_name, phone_number) VALUES (
    @first_name, @last_name, @phone_number
);


-- name: ListEmployeesRoles :many
SELECT * FROM employee_role;


-- name: GetServiceRepairByRepairIdAndServiceName :one
SELECT 
    sr.service_repairs_id
FROM
    public.service_repairs sr
LEFT JOIN public.service s ON sr.service_id = s.service_id
WHERE
    sr.repair_id = @repair_id
    AND s.name = @service_name
ORDER BY
    sr.service_repairs_id;


-- name: StoreServiceInventory :exec
INSERT INTO service_inventory (service_repair_id, inventory_id, amount) VALUES (
    @service_repair_id, @inventory_id, @amount
);
