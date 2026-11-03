CREATE OR REPLACE PROCEDURE stage.transfer_data_to_stage()
LANGUAGE plpgsql
AS $$
BEGIN
    -- Truncate and insert data into stage.service_inventory from public.service_inventory
    TRUNCATE TABLE stage.service_inventory;
    INSERT INTO stage.service_inventory (service_inventory_id, service_repair_id, inventory_id, amount)
    SELECT service_inventory_id, service_repair_id, inventory_id, amount
    FROM public.service_inventory;

    -- Truncate and insert data into stage.service_repair from public.service_repairs
    TRUNCATE TABLE stage.service_repair;
    INSERT INTO stage.service_repair (service_repairs_id, service_id, repair_id, quantity, completion_date, service_total_price)
    SELECT service_repairs_id, service_id, repair_id, quantity, completion_date, service_total_price
    FROM public.service_repairs;

    -- Truncate and insert data into stage.service from public.service
    TRUNCATE TABLE stage.service;
    INSERT INTO stage.service (service_id, name, description, price, estimated_duration)
    SELECT service_id, name, description, price, estimated_duration
    FROM public.service;

    -- Truncate and insert data into stage.inventory from public.inventory
    TRUNCATE TABLE stage.inventory;
    INSERT INTO stage.inventory (inventory_id, detail_id, quantity, service_station_id)
    SELECT inventory_id, detail_id, quantity, service_station_id
    FROM public.inventory;

    -- Truncate and insert data into stage.detail from public.detail
    TRUNCATE TABLE stage.detail;
    INSERT INTO stage.detail (detail_id, description, manufacturer_id, detail_type_id, article, price, dimensions, weight)
    SELECT detail_id, description, manufacturer_id, detail_type_id, article, price, dimensions, weight
    FROM public.detail;

    -- Truncate and insert data into stage.repair from public.repair_order
    TRUNCATE TABLE stage.repair;
    INSERT INTO stage.repair (repair_id, employee_id, start_date, repair_status_id, notes, repair_number, payment_status_id, car_client_id, total_sum)
    SELECT repair_id, employee_id, start_date, repair_status_id, notes, repair_number, payment_status_id, car_client_id, total_sum
    FROM public.repair_order;

    -- Truncate and insert data into stage.supply_details from public.supply_details
    TRUNCATE TABLE stage.supply_details;
    INSERT INTO stage.supply_details (supply_details_id, detail_id, quantity, price_per_unit, supply_id)
    SELECT supply_details_id, detail_id, quantity, price_per_unit, supply_id
    FROM public.supply_details;

    -- Truncate and insert data into stage.supplier from public.supplier
    TRUNCATE TABLE stage.supplier;
    INSERT INTO stage.supplier (supplier_id, name, city_id, phone_number, address)
    SELECT supplier_id, name, city_id, phone_number, address
    FROM public.supplier;

    -- Truncate and insert data into stage.supply from public.supply
    TRUNCATE TABLE stage.supply;
    INSERT INTO stage.supply (supply_id, supply_number, total_sum, supply_date, status_id, supplier_id, employee_id)
    SELECT supply_id, supply_number, total_sum, supply_date, status_id, supplier_id, employee_id
    FROM public.supply;

    -- Optionally log the completion of the transfer
    RAISE NOTICE 'Data transfer to stage schema completed successfully.';
END;
$$;
