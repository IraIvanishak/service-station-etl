-- name: ListSupply :many
SELECT 
    s.supply_id,
    s.supply_number,
    s.supply_date,
    sd.detail_id,
    sd.quantity,
    sd.price_per_unit,
    su.name AS supplier_name,
    e.first_name AS employee_first_name,
    e.last_name AS employee_last_name,
    ss.title AS supply_status,
    si.invoice_number as invoice_number,
    ss2.address as service_station,
    d.article as article,
    d.description as description
    
FROM 
    supply s
left JOIN 
    supply_details sd ON s.supply_id = sd.supply_id
left JOIN 
    supplier su ON s.supplier_id = su.supplier_id
left JOIN 
    employee e ON s.employee_id = e.employee_id
left JOIN 
    supply_status ss ON s.status_id = ss.status_id
left join 
	supply_invoice si on si.supply_details_id = sd.supply_details_id
left join 
	inventory i on si.inventory_id = i.inventory_id 
left join 
	service_station ss2 on ss2.station_id = i.service_station_id 
left join 
	detail d on d.detail_id = sd.detail_id;


-- name: NewSupply :one
INSERT INTO supply (
    supply_number,
    supply_date,
    supplier_id,
    employee_id,
    total_sum,
    status_id
)
VALUES (
    @supply_number,
    @supply_date,
    @supplier_id,
    @employee_id,
    @total_sum,
    @status_id
)
RETURNING supply_id;

-- name: GetSupplyDetailPrice :one
SELECT 
    price
FROM
    detail
WHERE
    detail_id = @detail_id;

-- name: NewSupplyDetail :one
INSERT INTO supply_details
(supply_id, detail_id, quantity, price_per_unit)
Values
(@supply_id, @detail_id, @quantity, @price_per_unit)
Returning supply_details_id;

-- name: StoreInventory :one
UPDATE inventory
SET quantity = quantity + @quantity
WHERE detail_id = @detail_id
RETURNING inventory_id;


-- name: NewSupplyInvoice :one
INSERT INTO supply_invoice
(supply_details_id, invoice_number, inventory_id)
Values
(@supply_details_id, @invoice_number, @inventory_id)
Returning supply_invoice_id;


-- name: GetSupplyDetailsByID :many
SELECT 
    sd.supply_details_id,
    sd.detail_id,
    sd.quantity,
    sd.price_per_unit
FROM
    supply_details sd
WHERE
    sd.supply_id = @supply_id;


-- name: GetInventoryIDbyDetailsID :one
SELECT 
    inventory_id
FROM
    inventory
WHERE
    detail_id = @detail_id
ORDER BY inventory_id DESC
LIMIT 1;

-- name: SetSupplyStatus :one
UPDATE supply
SET status_id = @status_id
WHERE supply_id = @supply_id
RETURNING supply_id;

-- name: DeleteSupply :exec
DELETE FROM supply
WHERE supply_id = @supply_id;


-- name: DeleteSupplyDetailsBySupplyId :exec
UPDATE supply_details
SET supply_id = NULL
WHERE supply_id = @supply_id;
    