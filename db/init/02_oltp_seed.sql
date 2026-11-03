-- Demo data for the OLTP schema.
--
-- The seed files are DBeaver exports without primary keys, so ids come from
-- the serial sequences and the files must run in foreign-key order.
--
-- In the original database some tables did not start at id 1 (rows had been
-- deleted and re-inserted), and the exported foreign keys still point at those
-- ids. first_id() moves the sequence so the re-inserted rows get them back.
-- The offsets were recovered from the foreign-key ranges in the export and
-- cross-checked against the old warehouse snapshot.

CREATE FUNCTION pg_temp.first_id(tbl text, col text, first_value bigint) RETURNS void
LANGUAGE sql AS $$ SELECT setval(pg_get_serial_sequence(tbl, col), first_value - 1, true) $$;

-- Lookups
\ir seed/country.sql
SELECT pg_temp.first_id('public.brand', 'brand_id', 51) \g /dev/null
\ir seed/brand.sql
\ir seed/characteristic.sql
\ir seed/detail_group.sql
\ir seed/employee_role.sql
\ir seed/payment_status.sql
\ir seed/repair_status.sql
\ir seed/supply_status.sql
SELECT pg_temp.first_id('public.service', 'service_id', 45) \g /dev/null
\ir seed/service.sql

-- Reference data
\ir seed/city.sql
SELECT pg_temp.first_id('public.car', 'car_id', 301) \g /dev/null
\ir seed/car.sql
\ir seed/client.sql
\ir seed/detail_type.sql
\ir seed/manufacturer.sql
\ir seed/service_station.sql
\ir seed/supplier.sql
\ir seed/detail.sql
\ir seed/detail_characteristic.sql
\ir seed/employee.sql
\ir seed/inventory.sql
\ir seed/car_client.sql

-- Transactions
SELECT pg_temp.first_id('public.repair_order', 'repair_id', 1201) \g /dev/null
\ir seed/repair_order.sql
SELECT pg_temp.first_id('public.service_repairs', 'service_repairs_id', 4) \g /dev/null
\ir seed/service_repairs.sql
SELECT pg_temp.first_id('public.service_inventory', 'service_inventory_id', 2164) \g /dev/null
\ir seed/service_inventory.sql
\ir seed/supply.sql
SELECT pg_temp.first_id('public.supply_details', 'supply_details_id', 109) \g /dev/null
\ir seed/supply_details.sql
\ir seed/supply_invoice.sql
