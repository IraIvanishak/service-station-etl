CREATE OR REPLACE PROCEDURE stage.fill_staging_from_oltp_incremental()
LANGUAGE plpgsql
AS $$
DECLARE
    last_load_time TIMESTAMP;
begin

    -- Отримати останній час завантаження з таблиці метаданих
    -- Before the first load there is no history: take everything. (Comparing
    -- with NULL matched no rows, so the first incremental load was empty.)
    SELECT COALESCE(max(load_datetime), '-infinity') INTO last_load_time
    FROM meta.dataloadhistory;

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
    -- LEFT JOIN: a repair without services yet is still a changed repair.
    LEFT JOIN public.service_repairs sr ON r.repair_id = sr.repair_id
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

-- fill_staging_from_oltp_incremental() stages only rows that changed. The
-- warehouse load joins stage tables to look up prices, dates and durations,
-- and recomputes per-repair and per-supply-group aggregates, so it also needs
-- the unchanged rows those changes depend on. Without them it silently
-- dropped facts (a new parts usage of an unchanged detail), overwrote supply
-- totals with the sum of the changed lines only, and kept stale estimated
-- durations for the other services of a changed repair.
--
-- Called by storage.fill_olap_from_staging_incremental() before it loads, so
-- the change counts reported from stage.* before that call stay "changed rows".
-- Deletions are still not propagated; run a full load for those.
CREATE OR REPLACE PROCEDURE stage.stage_referenced_rows()
LANGUAGE plpgsql
AS $$
BEGIN
    -- The service line behind every staged parts usage.
    INSERT INTO stage.service_repair (service_repairs_id, service_id, repair_id, quantity, completion_date, service_total_price)
    SELECT sr.service_repairs_id, sr.service_id, sr.repair_id, sr.quantity, sr.completion_date, sr.service_total_price
    FROM public.service_repairs sr
    WHERE sr.service_repairs_id IN (SELECT service_repair_id FROM stage.service_inventory)
    ON CONFLICT (service_repairs_id) DO NOTHING;

    -- Every service line of an affected repair: estimated duration is per repair.
    INSERT INTO stage.service_repair (service_repairs_id, service_id, repair_id, quantity, completion_date, service_total_price)
    SELECT sr.service_repairs_id, sr.service_id, sr.repair_id, sr.quantity, sr.completion_date, sr.service_total_price
    FROM public.service_repairs sr
    WHERE sr.repair_id IN (SELECT repair_id FROM stage.service_repair UNION SELECT repair_id FROM stage.repair)
    ON CONFLICT (service_repairs_id) DO NOTHING;

    -- Every parts usage of those service lines, so their facts are recomputed.
    INSERT INTO stage.service_inventory (service_inventory_id, service_repair_id, inventory_id, amount)
    SELECT si.service_inventory_id, si.service_repair_id, si.inventory_id, si.amount
    FROM public.service_inventory si
    WHERE si.service_repair_id IN (SELECT service_repairs_id FROM stage.service_repair)
    ON CONFLICT (service_inventory_id) DO NOTHING;

    -- The repair and service behind every staged service line.
    INSERT INTO stage.repair (repair_id, employee_id, start_date, repair_status_id, notes, repair_number, payment_status_id, car_client_id, total_sum)
    SELECT r.repair_id, r.employee_id, r.start_date, r.repair_status_id, r.notes, r.repair_number, r.payment_status_id, r.car_client_id, r.total_sum
    FROM public.repair_order r
    WHERE r.repair_id IN (SELECT repair_id FROM stage.service_repair)
    ON CONFLICT (repair_id) DO NOTHING;

    INSERT INTO stage.service (service_id, name, description, price, estimated_duration)
    SELECT s.service_id, s.name, s.description, s.price, s.estimated_duration
    FROM public.service s
    WHERE s.service_id IN (SELECT service_id FROM stage.service_repair)
    ON CONFLICT (service_id) DO NOTHING;

    INSERT INTO stage.inventory (inventory_id, detail_id, quantity, service_station_id)
    SELECT i.inventory_id, i.detail_id, i.quantity, i.service_station_id
    FROM public.inventory i
    WHERE i.inventory_id IN (SELECT inventory_id FROM stage.service_inventory)
    ON CONFLICT (inventory_id) DO NOTHING;

    -- Every line of an affected supply group (detail, supplier, supply date):
    -- the fact row holds the group total and is replaced on upsert.
    INSERT INTO stage.supply_details (supply_details_id, detail_id, quantity, price_per_unit, supply_id)
    SELECT sd.supply_details_id, sd.detail_id, sd.quantity, sd.price_per_unit, sd.supply_id
    FROM public.supply_details sd
    JOIN public.supply s ON s.supply_id = sd.supply_id
    WHERE (sd.detail_id, s.supplier_id, s.supply_date) IN (
        SELECT sd2.detail_id, s2.supplier_id, s2.supply_date
        FROM stage.supply_details sd2
        JOIN public.supply s2 ON s2.supply_id = sd2.supply_id)
    ON CONFLICT (supply_details_id) DO NOTHING;

    INSERT INTO stage.supply (supply_id, supply_number, supply_date, total_sum, status_id, supplier_id, employee_id)
    SELECT s.supply_id, s.supply_number, s.supply_date, s.total_sum, s.status_id, s.supplier_id, s.employee_id
    FROM public.supply s
    WHERE s.supply_id IN (SELECT supply_id FROM stage.supply_details)
    ON CONFLICT (supply_id) DO NOTHING;

    INSERT INTO stage.supplier (supplier_id, name, city_id, phone_number, address)
    SELECT sp.supplier_id, sp.name, sp.city_id, sp.phone_number, sp.address
    FROM public.supplier sp
    WHERE sp.supplier_id IN (SELECT supplier_id FROM stage.supply)
    ON CONFLICT (supplier_id) DO NOTHING;

    -- Details used by staged stock rows and supply lines.
    INSERT INTO stage.detail (detail_id, description, manufacturer_id, detail_type_id, article, price, dimensions, weight)
    SELECT d.detail_id, d.description, d.manufacturer_id, d.detail_type_id, d.article, d.price, d.dimensions, d.weight
    FROM public.detail d
    WHERE d.detail_id IN (SELECT detail_id FROM stage.inventory UNION SELECT detail_id FROM stage.supply_details)
    ON CONFLICT (detail_id) DO NOTHING;
END;
$$;
