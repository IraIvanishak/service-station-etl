CREATE OR REPLACE PROCEDURE stage.fill_staging_from_oltp_incremental()
LANGUAGE plpgsql
AS $$
DECLARE
    last_load_time TIMESTAMP;
begin

    -- Отримати останній час завантаження з таблиці метаданих
    SELECT load_datetime INTO last_load_time 
    FROM meta.dataloadhistory 
    ORDER BY load_datetime DESC 
    LIMIT 1;

    -- Вивести отриманий час завантаження
    RAISE NOTICE 'Last load time: %', last_load_time;
    -- Оновлення таблиці detail
	 TRUNCATE TABLE stage.detail;
    INSERT INTO stage.detail (detail_id, description, manufacturer_id, detail_type_id, article, price, dimensions, weight)
    SELECT DISTINCT d.detail_id, d.description, d.manufacturer_id, d.detail_type_id, d.article, d.price, d.dimensions, d.weight
    FROM public.detail d
    LEFT JOIN public.supply_details sd ON d.detail_id = sd.detail_id
    WHERE d.updated_at > last_load_time OR d.created_at > last_load_time
       OR sd.updated_at > last_load_time OR sd.created_at > last_load_time;

    -- Оновлення таблиці repair
  	 TRUNCATE TABLE stage.repair;
    INSERT INTO stage.repair (repair_id, employee_id, start_date, repair_status_id, notes, repair_number, payment_status_id, car_client_id, total_sum)
    SELECT DISTINCT r.repair_id, r.employee_id, r.start_date, r.repair_status_id, r.notes, r.repair_number, r.payment_status_id, r.car_client_id, r.total_sum
    FROM public.repair_order  r
    JOIN public.service_repairs sr ON r.repair_id = sr.repair_id
    WHERE r.updated_at > last_load_time OR r.created_at > last_load_time
       OR sr.updated_at > last_load_time OR sr.created_at > last_load_time;

    -- Оновлення таблиці service
	 TRUNCATE TABLE stage.service;
    INSERT INTO stage.service (service_id, name, description, price, estimated_duration)
    SELECT DISTINCT s.service_id, s.name, s.description, s.price, s.estimated_duration
    FROM public.service s
    LEFT JOIN public.service_repairs sr ON s.service_id = sr.service_id
    WHERE s.updated_at > last_load_time OR s.created_at > last_load_time
       OR sr.updated_at > last_load_time OR sr.created_at > last_load_time;
   
     -- Оновлення таблиці inventory
 	 TRUNCATE TABLE stage.inventory;
    INSERT INTO stage.inventory (inventory_id, detail_id, quantity, service_station_id)
    SELECT DISTINCT i.inventory_id, i.detail_id, i.quantity, i.service_station_id
    FROM public.inventory i
    LEFT JOIN public.detail d ON i.detail_id = d.detail_id
    LEFT JOIN public.service_station ss ON i.service_station_id = ss.station_id
    WHERE i.updated_at > last_load_time OR i.created_at > last_load_time
       OR d.updated_at > last_load_time OR d.created_at > last_load_time
       OR ss.updated_at > last_load_time OR ss.created_at > last_load_time;

    -- Оновлення таблиці service_inventory
    TRUNCATE TABLE stage.service_inventory;
    INSERT INTO stage.service_inventory (service_inventory_id, service_repair_id, inventory_id, amount)
    SELECT DISTINCT si.service_inventory_id, si.service_repair_id, si.inventory_id, si.amount
    FROM public.service_inventory si
    JOIN public.inventory i ON si.inventory_id = i.inventory_id
    JOIN public.detail d ON i.detail_id = d.detail_id
    JOIN public.service_repairs sr ON si.service_repair_id = sr.service_repairs_id
    JOIN public.repair_order r ON sr.repair_id = r.repair_id
    WHERE si.updated_at > last_load_time OR si.created_at > last_load_time
       OR i.updated_at > last_load_time OR i.created_at > last_load_time
       OR d.updated_at > last_load_time OR d.created_at > last_load_time
       OR sr.updated_at > last_load_time OR sr.created_at > last_load_time
       OR r.updated_at > last_load_time OR r.created_at > last_load_time;

    -- Оновлення таблиці service_repair
    TRUNCATE TABLE stage.service_repair;
    INSERT INTO stage.service_repair (service_repairs_id, service_id, repair_id, quantity, completion_date, service_total_price)
    SELECT DISTINCT sr.service_repairs_id, sr.service_id, sr.repair_id, sr.quantity, sr.completion_date, sr.service_total_price
    FROM public.service_repairs sr
    JOIN public.repair_order r ON sr.repair_id = r.repair_id
    JOIN public.service s ON sr.service_id = s.service_id
    WHERE r.updated_at > last_load_time OR r.created_at > last_load_time
       OR sr.updated_at > last_load_time OR sr.created_at > last_load_time
       OR s.updated_at > last_load_time OR s.created_at > last_load_time;

    -- Оновлення таблиці supply
    TRUNCATE TABLE stage.supply;
    INSERT INTO stage.supply (supply_id, supply_date, supply_number, total_sum, status_id, supplier_id, employee_id)
    SELECT DISTINCT s.supply_id, s.supply_date, s.supply_number, s.total_sum, s.status_id, s.supplier_id, s.employee_id
    FROM public.supply s
    LEFT JOIN public.supply_details sd ON s.supply_id = sd.supply_id
    LEFT JOIN public.detail d ON sd.detail_id = d.detail_id
    LEFT JOIN public.inventory i ON sd.detail_id = i.detail_id
    LEFT JOIN public.supplier sp ON s.supplier_id = sp.supplier_id
    WHERE s.updated_at > last_load_time OR s.created_at > last_load_time
       OR sd.updated_at > last_load_time OR sd.created_at > last_load_time
       OR i.updated_at > last_load_time OR i.created_at > last_load_time
       OR d.updated_at > last_load_time OR d.created_at > last_load_time
       OR sp.updated_at > last_load_time OR sp.created_at > last_load_time;

    -- Оновлення таблиці supplier
    TRUNCATE TABLE stage.supplier;
    INSERT INTO stage.supplier (supplier_id, name, city_id, phone_number, address)
    SELECT DISTINCT supplier_id, name, city_id, phone_number, address
    FROM public.supplier
    WHERE updated_at > last_load_time OR created_at > last_load_time;

    -- Оновлення таблиці supply_details
    TRUNCATE TABLE stage.supply_details;
    INSERT INTO stage.supply_details (supply_details_id, detail_id, quantity, price_per_unit, supply_id)
    SELECT DISTINCT sd.supply_details_id, sd.detail_id, sd.quantity, sd.price_per_unit, sd.supply_id
    FROM public.supply_details sd
    JOIN public.supply s ON sd.supply_id = s.supply_id
    JOIN public.detail d ON sd.detail_id = d.detail_id
    WHERE sd.updated_at > last_load_time OR sd.created_at > last_load_time
       OR s.updated_at > last_load_time OR s.created_at > last_load_time
       OR d.updated_at > last_load_time OR d.created_at > last_load_time;

END;
$$;
