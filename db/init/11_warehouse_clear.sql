CREATE OR REPLACE PROCEDURE storage.clearstorage()
 LANGUAGE plpgsql
AS $procedure$
BEGIN
    -- Truncating existing data
    TRUNCATE TABLE storage.date_dim, storage.inventory_dim, storage.service_dim, 
           storage.supplier_dim, storage.repair_dim, storage.service_repair_dim,
           storage.supply_details_fact, storage.service_inventory_fact RESTART IDENTITY CASCADE;
END;
$procedure$
;
