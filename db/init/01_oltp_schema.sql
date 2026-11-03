CREATE TABLE public.brand (
	brand_id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT brand_pkey PRIMARY KEY (brand_id)
);

CREATE TABLE public.characteristic (
	characteristic_id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT characteristic_pkey PRIMARY KEY (characteristic_id)
);

CREATE TABLE public.client (
	client_id serial4 NOT NULL,
	first_name varchar(255) NOT NULL,
	last_name varchar(255) NOT NULL,
	phone_number varchar(20) NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT client_pkey PRIMARY KEY (client_id)
);

CREATE TABLE public.country (
	country_id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT country_pkey PRIMARY KEY (country_id)
);

CREATE TABLE public.detail_group (
	detail_group_id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT detail_group_pkey PRIMARY KEY (detail_group_id)
);

CREATE TABLE public.employee_role (
	role_id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT employee_role_pkey PRIMARY KEY (role_id)
);

CREATE TABLE public.payment_status (
	status_id serial4 NOT NULL,
	title varchar(255) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT payment_status_pkey PRIMARY KEY (status_id)
);

CREATE TABLE public.repair_status (
	status_id serial4 NOT NULL,
	title varchar(255) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT repair_status_pkey PRIMARY KEY (status_id)
);

CREATE TABLE public.service (
	service_id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	description text NULL,
	price numeric(10, 2) NULL,
	estimated_duration interval NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT service_pkey PRIMARY KEY (service_id)
);

CREATE TABLE public.supply_status (
	status_id serial4 NOT NULL,
	title varchar(255) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT supply_status_pkey PRIMARY KEY (status_id)
);

CREATE TABLE public.car (
	car_id serial4 NOT NULL,
	"year" int4 NULL,
	body_type varchar(50) NULL,
	model varchar(255) NULL,
	brand_id int4 NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT car_pkey PRIMARY KEY (car_id),
	CONSTRAINT car_brand_id_fkey FOREIGN KEY (brand_id) REFERENCES public.brand(brand_id)
);

CREATE TABLE public.car_client (
	car_client_id serial4 NOT NULL,
	car_id int4 NOT NULL,
	client_id int4 NOT NULL,
	description text NULL,
	registration_number varchar(50) NULL,
	vin varchar(17) NULL,
	color varchar(50) NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT car_client_pkey PRIMARY KEY (car_client_id),
	CONSTRAINT car_client_car_id_fkey FOREIGN KEY (car_id) REFERENCES public.car(car_id),
	CONSTRAINT car_client_client_id_fkey FOREIGN KEY (client_id) REFERENCES public.client(client_id)
);

CREATE TABLE public.city (
	city_id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	country_id int4 NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT city_pkey PRIMARY KEY (city_id),
	CONSTRAINT city_country_id_fkey FOREIGN KEY (country_id) REFERENCES public.country(country_id)
);

CREATE TABLE public.detail_type (
	detail_type_id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	detail_group_id int4 NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT detail_type_pkey PRIMARY KEY (detail_type_id),
	CONSTRAINT detail_type_detail_group_id_fkey FOREIGN KEY (detail_group_id) REFERENCES public.detail_group(detail_group_id)
);

CREATE TABLE public.manufacturer (
	manufacturer_id serial4 NOT NULL,
	country_id int4 NULL,
	"name" varchar(255) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT manufacturer_pkey PRIMARY KEY (manufacturer_id),
	CONSTRAINT manufacturer_country_id_fkey FOREIGN KEY (country_id) REFERENCES public.country(country_id)
);

CREATE TABLE public.service_station (
	station_id serial4 NOT NULL,
	address text NOT NULL,
	city_id int4 NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT service_station_pkey PRIMARY KEY (station_id),
	CONSTRAINT service_station_city_id_fkey FOREIGN KEY (city_id) REFERENCES public.city(city_id)
);

CREATE TABLE public.supplier (
	supplier_id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	city_id int4 NULL,
	phone_number varchar(20) NULL,
	address text NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT supplier_pkey PRIMARY KEY (supplier_id),
	CONSTRAINT supplier_city_id_fkey FOREIGN KEY (city_id) REFERENCES public.city(city_id)
);

CREATE TABLE public.detail (
	detail_id serial4 NOT NULL,
	description text NOT NULL,
	manufacturer_id int4 NULL,
	detail_type_id int4 NULL,
	article int4 NULL,
	price numeric(10, 2) NULL,
	dimensions varchar(255) NULL,
	weight numeric(10, 2) NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT detail_pkey PRIMARY KEY (detail_id),
	CONSTRAINT detail_detail_type_id_fkey FOREIGN KEY (detail_type_id) REFERENCES public.detail_type(detail_type_id),
	CONSTRAINT detail_manufacturer_id_fkey FOREIGN KEY (manufacturer_id) REFERENCES public.manufacturer(manufacturer_id)
);

CREATE TABLE public.detail_characteristic (
	characteristic_id int4 NOT NULL,
	detail_id int4 NOT NULL,
	value text NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT detail_characteristic_pkey PRIMARY KEY (characteristic_id, detail_id),
	CONSTRAINT detail_characteristic_characteristic_id_fkey FOREIGN KEY (characteristic_id) REFERENCES public.characteristic(characteristic_id),
	CONSTRAINT detail_characteristic_detail_id_fkey FOREIGN KEY (detail_id) REFERENCES public.detail(detail_id)
);

CREATE TABLE public.employee (
	employee_id serial4 NOT NULL,
	first_name varchar(255) NOT NULL,
	last_name varchar(255) NOT NULL,
	phone_number varchar(20) NULL,
	role_id int4 NULL,
	service_station_id int4 NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT employee_pkey PRIMARY KEY (employee_id),
	CONSTRAINT employee_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.employee_role(role_id),
	CONSTRAINT employee_service_station_id_fkey FOREIGN KEY (service_station_id) REFERENCES public.service_station(station_id)
);

CREATE TABLE public.inventory (
	inventory_id serial4 NOT NULL,
	detail_id int4 NULL,
	quantity int4 NULL,
	service_station_id int4 NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT inventory_pkey PRIMARY KEY (inventory_id),
	CONSTRAINT inventory_detail_id_fkey FOREIGN KEY (detail_id) REFERENCES public.detail(detail_id),
	CONSTRAINT inventory_service_station_id_fkey FOREIGN KEY (service_station_id) REFERENCES public.service_station(station_id)
);

CREATE TABLE public.repair_order (
	repair_id serial4 NOT NULL,
	employee_id int4 NULL,
	start_date timestamp NULL,
	repair_status_id int4 NULL,
	notes text NULL,
	repair_number varchar(255) NULL,
	payment_status_id int4 NULL,
	car_client_id int4 NULL,
	total_sum numeric(10, 2) NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT repair_pkey PRIMARY KEY (repair_id),
	CONSTRAINT repair_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.employee(employee_id),
	CONSTRAINT repair_order_car_client_id_fkey FOREIGN KEY (car_client_id) REFERENCES public.car_client(car_client_id),
	CONSTRAINT repair_order_payment_status_id_fkey FOREIGN KEY (payment_status_id) REFERENCES public.payment_status(status_id),
	CONSTRAINT repair_status_id_fkey FOREIGN KEY (repair_status_id) REFERENCES public.repair_status(status_id)
);

CREATE TABLE public.service_repairs (
	service_repairs_id serial4 NOT NULL,
	service_id int4 NULL,
	repair_id int4 NULL,
	quantity int4 NULL,
	completion_date date NULL,
	service_total_price numeric(10, 2) NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT service_repairs_pkey PRIMARY KEY (service_repairs_id),
	CONSTRAINT service_repairs_repair_id_fkey FOREIGN KEY (repair_id) REFERENCES public.repair_order(repair_id),
	CONSTRAINT service_repairs_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.service(service_id)
);

CREATE TABLE public.supply (
	supply_id serial4 NOT NULL,
	supply_number varchar(255) NOT NULL,
	total_sum numeric(10, 2) NULL,
	supply_date date NOT NULL,
	status_id int4 NULL,
	supplier_id int4 NULL,
	employee_id int4 NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT supply_pkey PRIMARY KEY (supply_id),
	CONSTRAINT supply_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.employee(employee_id),
	CONSTRAINT supply_status_id_fkey FOREIGN KEY (status_id) REFERENCES public.supply_status(status_id),
	CONSTRAINT supply_supplier_id_fkey FOREIGN KEY (supplier_id) REFERENCES public.supplier(supplier_id)
);

CREATE TABLE public.supply_details (
	supply_details_id serial4 NOT NULL,
	detail_id int4 NULL,
	quantity int4 NULL,
	price_per_unit numeric(10, 2) NULL,
	supply_id int4 NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT supply_details_pkey PRIMARY KEY (supply_details_id),
	CONSTRAINT supply_details_detail_id_fkey FOREIGN KEY (detail_id) REFERENCES public.detail(detail_id),
	CONSTRAINT supply_details_supply_id_fkey FOREIGN KEY (supply_id) REFERENCES public.supply(supply_id)
);

CREATE TABLE public.supply_invoice (
	supply_invoice_id serial4 NOT NULL,
	supply_details_id int4 NULL,
	invoice_number varchar(255) NULL,
	inventory_id int4 NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT supply_invoice_pkey PRIMARY KEY (supply_invoice_id),
	CONSTRAINT fk_inventory FOREIGN KEY (inventory_id) REFERENCES public.inventory(inventory_id),
	CONSTRAINT supply_invoice_supply_details_id_fkey FOREIGN KEY (supply_details_id) REFERENCES public.supply_details(supply_details_id)
);

CREATE TABLE public.service_inventory (
	service_inventory_id serial4 NOT NULL,
	service_repair_id int4 NULL,
	inventory_id int4 NULL,
	amount numeric(10, 2) NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT service_inventory_pkey PRIMARY KEY (service_inventory_id),
	CONSTRAINT fk_iinventory FOREIGN KEY (inventory_id) REFERENCES public.inventory(inventory_id),
	CONSTRAINT fk_inventory FOREIGN KEY (inventory_id) REFERENCES public.inventory(inventory_id),
	CONSTRAINT service_inventory_inventory_id_fkey FOREIGN KEY (inventory_id) REFERENCES public.inventory(inventory_id),
	CONSTRAINT service_inventory_service_repair_id_fkey FOREIGN KEY (service_repair_id) REFERENCES public.service_repairs(service_repairs_id)
);


ALTER TABLE public.inventory
ADD CONSTRAINT detail_service_unique UNIQUE (detail_id, service_station_id);


DO $$
DECLARE
    rec record;
BEGIN
    FOR rec IN SELECT table_name, column_name, pg_get_serial_sequence(table_name, column_name) as seq_name
               FROM information_schema.columns
               WHERE column_default LIKE 'nextval(%' AND table_schema = 'public'
    LOOP
        IF rec.seq_name IS NOT NULL THEN
            EXECUTE format('SELECT setval(%L, COALESCE(MAX(%I), 0) + 1, false) FROM %I', 
                           rec.seq_name, rec.column_name, rec.table_name);
        END IF;
    END LOOP;
END $$;


-- Keep updated_at current on every UPDATE; the incremental ETL relies on it.
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    RAISE NOTICE 'Trigger called for table % on column % at time %', TG_TABLE_NAME, 'updated_at', CURRENT_TIMESTAMP;
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;




DO $$
DECLARE
    row record;
BEGIN
    FOR row IN SELECT table_name FROM information_schema.tables
               WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
    LOOP
        -- Перевірка чи існує колонка `updated_at`
        IF EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_schema = 'public' AND table_name = row.table_name AND column_name = 'updated_at') THEN
            EXECUTE format('CREATE TRIGGER update_%1$I_updated_at BEFORE UPDATE ON public.%1$I 
                            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()', row.table_name);
        END IF;
    END LOOP;
END;
$$;
