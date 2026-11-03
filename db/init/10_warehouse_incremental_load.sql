CREATE OR REPLACE PROCEDURE storage.fill_olap_from_staging_incremental()
LANGUAGE plpgsql
AS $$
DECLARE
    meta_start_time TIMESTAMP := clock_timestamp();
    row_count BIGINT := 0;
    inserted_count BIGINT := 0;
    updated_count BIGINT := 0;
    db_table_count INT := 0;
BEGIN

    -- Inserting data into supplier_dim
    -- !!!!!!!!!!!!!
    INSERT INTO storage.supplier_dim (supplier_dim_id, bk_supplier_id, name, phone_number, address)
    SELECT 
        supplier_id AS supplier_dim_id, 
        supplier_id AS bk_supplier_id, 
        name, 
        phone_number, 
        address
    FROM stage.supplier
    ON CONFLICT (bk_supplier_id) DO UPDATE SET
        name = EXCLUDED.name,
        phone_number = EXCLUDED.phone_number,
        address = EXCLUDED.address;


    GET DIAGNOSTICS inserted_count = ROW_COUNT;
    row_count := row_count + inserted_count;
    db_table_count := db_table_count + 1;
    

    -- Inserting data into service_dim
    --- !!!!!!!!!!!!
    INSERT INTO storage.service_dim (service_dim_id, bk_service_id, name, description, price, estimated_duration)
    SELECT 
        service_id AS service_dim_id,
        service_id AS bk_service_id, 	
        name, 	
        description, 	
        price, 
        estimated_duration
    FROM stage.service
    ON CONFLICT (bk_service_id) DO UPDATE SET
        service_dim_id = EXCLUDED.service_dim_id,
        name = EXCLUDED.name,
        description = EXCLUDED.description,
        price = EXCLUDED.price,
        estimated_duration = EXCLUDED.estimated_duration;

    GET DIAGNOSTICS inserted_count = ROW_COUNT;
    row_count := row_count + inserted_count;
    db_table_count := db_table_count + 1;

    -- Inserting data into inventory_dim 
    -- !!!!!!!!!!!!
    INSERT INTO storage.inventory_dim (inventory_dim_id, bk_detail_id, description, article, price, dimensions, weight)
    SELECT 
        d.detail_id AS inventory_dim_id, 
        d.detail_id AS bk_detail_id, 
        d.description, 
        d.article, 
        d.price, 
        d.dimensions, 
        d.weight
    -- One row per detail. Joining stage.inventory duplicated details stocked
    -- at several stations and dropped details that are not in stock.
    FROM stage.detail AS d
    ON CONFLICT (bk_detail_id) DO UPDATE SET
        description = EXCLUDED.description,
        article = EXCLUDED.article,
        price = EXCLUDED.price,
        dimensions = EXCLUDED.dimensions,
        weight = EXCLUDED.weight;


    GET DIAGNOSTICS inserted_count = ROW_COUNT;
    row_count := row_count + inserted_count;
    db_table_count := db_table_count + 1;


    -- Inserting data into date_dim
    -- !
    INSERT INTO storage.date_dim ("year", "month", "day")
    SELECT DISTINCT
        EXTRACT(YEAR FROM date_field) AS "year",
        EXTRACT(MONTH FROM date_field) AS "month",
        EXTRACT(DAY FROM date_field) AS "day"
    FROM (
        SELECT start_date AS date_field FROM stage.repair
        UNION
        SELECT completion_date AS date_field FROM stage.service_repair
        UNION
        SELECT supply_date AS date_field FROM stage.supply
        UNION
        SELECT (supply_date + INTERVAL '1 MONTH' - INTERVAL '1 DAY') AS date_field FROM stage.supply
    ) subquery
    WHERE NOT EXISTS (
        SELECT 1 FROM storage.date_dim
        WHERE "year" = EXTRACT(YEAR FROM date_field)
        AND "month" = EXTRACT(MONTH FROM date_field)
        AND "day" = EXTRACT(DAY FROM date_field)
    );

    GET DIAGNOSTICS inserted_count = ROW_COUNT;
    row_count := row_count + inserted_count;
    db_table_count := db_table_count + 1;


    -- Inserting data into repair_dim
    -- !!!!!!!!!!!!
    INSERT INTO storage.repair_dim (repair_dim_id, bk_repair_id, notes, repair_number, start_repair_date_dim_id)
    SELECT
        r.repair_id AS repair_dim_id,
        r.repair_id AS bk_repair_id,
        r.notes,
        r.repair_number,
        dd.date_id AS start_repair_date_dim_id
    FROM stage.repair r
    LEFT JOIN storage.date_dim dd ON r.start_date::date = (dd."year" || '-' || dd."month" || '-' || dd."day")::date
    ON CONFLICT (bk_repair_id) DO UPDATE SET
        repair_dim_id = EXCLUDED.repair_dim_id,
        notes = EXCLUDED.notes,
        repair_number = EXCLUDED.repair_number,
        start_repair_date_dim_id = EXCLUDED.start_repair_date_dim_id;

    GET DIAGNOSTICS inserted_count = ROW_COUNT;
    row_count := row_count + inserted_count;
    db_table_count := db_table_count + 1;


    -- Inserting data into service_repair_dim
    -- !!!!!!!!!!!!
    INSERT INTO storage.service_repair_dim (service_repair_dim_id, service_id, repair_id, quantity, service_total_price, complete_date_dim_id, bk_service_repair)
    SELECT
        sr.service_repairs_id,
        sd.service_dim_id,  
        rd.repair_dim_id,   
        sr.quantity,
        sr.service_total_price,
        dd.date_id,
        sr.service_repairs_id  as bk_service_repair
    FROM stage.service_repair sr
    JOIN storage.service_dim sd ON sr.service_id = sd.bk_service_id
    JOIN storage.repair_dim rd ON sr.repair_id = rd.bk_repair_id
    LEFT JOIN storage.date_dim dd ON sr.completion_date = (dd."year" || '-' || dd."month" || '-' || dd."day")::date
    ON CONFLICT (bk_service_repair) DO UPDATE SET
        service_repair_dim_id = EXCLUDED.service_repair_dim_id,
        service_id = EXCLUDED.service_id,
        repair_id = EXCLUDED.repair_id,
        quantity = EXCLUDED.quantity,
        service_total_price = EXCLUDED.service_total_price,
        complete_date_dim_id = EXCLUDED.complete_date_dim_id;

    GET DIAGNOSTICS inserted_count = ROW_COUNT;
    row_count := row_count + inserted_count;
    db_table_count := db_table_count + 1;


    -- Inserting data into service_inventory_fact
    -- !!!!!!!!!!!!
      WITH count_estimated_duration AS (
        SELECT 
            sr.repair_id, 
            SUM(s.estimated_duration * sr.quantity) AS total_estimated_duration
        FROM stage.service_repair sr
        JOIN stage.service s ON sr.service_id = s.service_id
        GROUP BY sr.repair_id
    )
    INSERT INTO storage.service_inventory_fact (
        service_inventory_fact_id,
        service_repair_dim_id, 
        inventory_dim_id, 
        parts_cost, 
        total_cost, 
        parts_cost_percentage, 
        actual_service_duration, 
        estimated_service_duration, 
        service_duration_deviation, 
        bk_service_inventory_id, 
        amount
    )
    SELECT 
        si.service_inventory_id AS service_inventory_fact_id,
        sr.service_repairs_id AS service_repair_dim_id,
        inv.inventory_dim_id AS inventory_dim_id,
        (d.price * si.amount) AS parts_cost,
        (sr.service_total_price + (d.price * si.amount)) AS total_cost,
        (d.price * si.amount) / NULLIF((sr.service_total_price + (d.price * si.amount)), 0) * 100 AS parts_cost_percentage,
        make_interval(days =>(sr.completion_date - r.start_date)) AS actual_service_duration,
        ed.total_estimated_duration AS estimated_service_duration,
        make_interval(days => (sr.completion_date - r.start_date)) - ed.total_estimated_duration  AS service_duration_deviation,
        si.service_inventory_id AS bk_service_inventory_id,
        si.amount
    FROM stage.service_inventory si
    JOIN stage.service_repair sr ON sr.service_repairs_id = si.service_repair_id
    JOIN stage.inventory i ON si.inventory_id = i.inventory_id
    JOIN stage.detail d ON i.detail_id = d.detail_id
    JOIN storage.inventory_dim inv ON inv.bk_detail_id = d.detail_id
    JOIN stage.service s ON sr.service_id = s.service_id
    JOIN stage.repair r ON sr.repair_id = r.repair_id
    JOIN count_estimated_duration ed ON sr.repair_id = ed.repair_id
    ON CONFLICT (bk_service_inventory_id) DO UPDATE SET
        service_repair_dim_id = EXCLUDED.service_repair_dim_id,
        inventory_dim_id = EXCLUDED.inventory_dim_id,
        parts_cost = EXCLUDED.parts_cost,
        total_cost = EXCLUDED.total_cost,
        parts_cost_percentage = EXCLUDED.parts_cost_percentage,
        actual_service_duration = EXCLUDED.actual_service_duration,
        estimated_service_duration = EXCLUDED.estimated_service_duration,
        service_duration_deviation = EXCLUDED.service_duration_deviation,
        amount = EXCLUDED.amount;


    GET DIAGNOSTICS inserted_count = ROW_COUNT;
    row_count := row_count + inserted_count;
    db_table_count := db_table_count + 1;


    -- Inserting data into supply_details_fact
    WITH supply_data AS (
    SELECT
        inv.inventory_dim_id AS inventory_dim_id,
        s.supply_id,
        sup.supplier_id,
        s.supply_date,
        sd.quantity,
        sd.price_per_unit
    FROM stage.supply_details sd
    JOIN stage.supply s ON sd.supply_id = s.supply_id
    JOIN stage.supplier sup ON s.supplier_id = sup.supplier_id
    JOIN stage.detail d ON sd.detail_id = d.detail_id
    JOIN storage.inventory_dim inv ON d.detail_id = inv.bk_detail_id
    ),
    -- One row per supply: these CTEs are joined on supply_id alone, so with one
    -- row per supply line a supply with k lines matched k*k times and its
    -- quantities and costs were multiplied by k^2.
    start_dates AS (
        SELECT DISTINCT
            sd.supply_id AS supply_id, 
            dd.date_id AS start_date_dim_id
        FROM supply_data sd
        JOIN storage.date_dim dd ON EXTRACT(YEAR FROM sd.supply_date) = dd.year 
        AND EXTRACT(MONTH FROM sd.supply_date) = dd.month 
        AND EXTRACT(DAY FROM sd.supply_date) = dd.day
    ),
    end_dates AS (
        SELECT DISTINCT
            sd.supply_id AS supply_id, 
            dd.date_id AS end_date_dim_id
        FROM supply_data sd
        JOIN storage.date_dim dd ON EXTRACT(YEAR FROM (sd.supply_date + INTERVAL '1 MONTH' - INTERVAL '1 DAY')) = dd.year 
        AND EXTRACT(MONTH FROM (sd.supply_date + INTERVAL '1 MONTH' - INTERVAL '1 DAY')) = dd.month 
        AND EXTRACT(DAY FROM (sd.supply_date + INTERVAL '1 MONTH' - INTERVAL '1 DAY')) = dd.day
    )
    INSERT INTO storage.supply_details_fact (
        detail_dim_id,
        detail_quantity,
        detail_total_cost,
        supplier_id,
        start_date_dim_id,
        end_date_dim_id
    )
    SELECT 
        sd.inventory_dim_id AS detail_dim_id,
        SUM(sd.quantity) AS detail_quantity,
        SUM(sd.quantity * sd.price_per_unit) AS detail_total_cost,
        sd.supplier_id,
        st.start_date_dim_id,
        ed.end_date_dim_id 
    FROM 
        supply_data sd
    JOIN start_dates st ON sd.supply_id = st.supply_id
    JOIN end_dates ed ON sd.supply_id = ed.supply_id
    GROUP BY 
        sd.inventory_dim_id, 
        sd.supplier_id,
        st.start_date_dim_id,
        ed.end_date_dim_id
    ON CONFLICT (detail_dim_id, supplier_id, start_date_dim_id, end_date_dim_id) DO UPDATE SET
        detail_quantity = EXCLUDED.detail_quantity,
        detail_total_cost = EXCLUDED.detail_total_cost;


    GET DIAGNOSTICS inserted_count = ROW_COUNT;
    row_count := row_count + inserted_count;
    db_table_count := db_table_count + 1;


    -- Record load activity in the metadata database
    INSERT INTO meta.dataloadhistory(load_datetime, load_time, load_rows, affected_table_count, source_table_count)
    VALUES (
        clock_timestamp(), 
        make_interval(secs => EXTRACT(EPOCH FROM clock_timestamp() - meta_start_time)),
        row_count, 
        db_table_count, 
        db_table_count
    );
    
    -- Retrieve the last inserted ID for use in further operations if necessary
    UPDATE meta.dwtable
    SET data_load_history_id = (SELECT MAX(data_load_history_id) FROM meta.dataloadhistory)
    WHERE dwtable.data_load_history_id IS NULL;

    -- Optionally, log the process completion
    RAISE NOTICE 'Data load completed with, affecting % rows across % tables.', row_count, db_table_count;
END;
$$;
