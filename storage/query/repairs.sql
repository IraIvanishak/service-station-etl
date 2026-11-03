-- name: ListOrders :many
SELECT 
    ro.repair_id,
    ro.repair_number,
    ro.start_date,
    ro.notes,
    ro.total_sum,
    rs.title AS repair_status,
    ps.title AS payment_status,
    concat(emp.first_name, ' ', emp.last_name)::text AS assigned_employee,
    cc.registration_number,
    cc.description AS vehicle_description,
    concat(cl.first_name, ' ', cl.last_name)::text AS client_name,
    car.model AS vehicle_model,
    car.year AS vehicle_year,
    srep.service_repairs_id,
    ser.name AS service_name,
    srep.quantity AS service_quantity,
    srep.completion_date,
    srep.service_total_price,
    inv.service_inventory_id,
    inv.amount AS inventory_amount,
    det.description AS inventory_item_description,
    det.price AS detail_item_price_per_unit,
    invdet.quantity AS inventory_item_quantity,
    ss.address  as address 
FROM 
    public.repair_order ro
LEFT JOIN public.repair_status rs ON ro.repair_status_id = rs.status_id
LEFT JOIN public.payment_status ps ON ro.payment_status_id = ps.status_id
LEFT JOIN public.employee emp ON ro.employee_id = emp.employee_id
LEFT JOIN public.car_client cc ON ro.car_client_id = cc.car_client_id
LEFT JOIN public.client cl ON cc.client_id = cl.client_id
LEFT JOIN public.car car ON cc.car_id = car.car_id
LEFT JOIN public.service_repairs srep ON ro.repair_id = srep.repair_id
LEFT JOIN public.service ser ON srep.service_id = ser.service_id
LEFT JOIN public.service_inventory inv ON srep.service_repairs_id = inv.service_repair_id
LEFT JOIN public.inventory invdet ON inv.inventory_id = invdet.inventory_id
LEFT JOIN public.detail det ON invdet.detail_id = det.detail_id
left join public.service_station ss on emp.service_station_id = ss.station_id 
ORDER BY 
    ro.repair_id, srep.service_repairs_id, inv.service_inventory_id;
