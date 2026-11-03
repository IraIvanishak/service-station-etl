CREATE SCHEMA IF NOT EXISTS meta;

-- Створення таблиці для відстеження історії завантаження даних
CREATE TABLE meta.dataloadhistory (
    data_load_history_id serial4 PRIMARY KEY,
    load_datetime timestamp without time zone,
    load_time time without time zone,
    load_rows integer,
    affected_table_count integer,
    source_table_count integer
);

-- Створення таблиці для опису таблиць у сховищі даних
CREATE TABLE meta.dwtable (
    dw_table_id serial4 PRIMARY KEY,
    data_load_history_id integer,
    dw_table_name character varying,
    FOREIGN KEY (data_load_history_id) REFERENCES meta.dataloadhistory(data_load_history_id)
);

-- Створення таблиці для опису бази даних джерела
CREATE TABLE meta.sourcedb (
    source_db_id serial4 PRIMARY KEY,
    source_db_name character varying
);

-- Створення таблиці для опису таблиць джерела даних
CREATE TABLE meta.sourcetable (
    source_table_id serial4 PRIMARY KEY,
    source_db_id integer,
    key_name character varying,
    source_table_name character varying,
    FOREIGN KEY (source_db_id) REFERENCES meta.sourcedb(source_db_id)
);

-- Створення таблиці для опису вимірювань (вимірювальних таблиць)
CREATE TABLE meta.dimension (
    dimension_id serial4 PRIMARY KEY,
    dw_table_id integer,
    key_name character varying,
    dimension_type character varying,
    dimension_name character varying,
    FOREIGN KEY (dw_table_id) REFERENCES meta.dwtable(dw_table_id)
);

-- Створення таблиці для опису колонок атрибутів
CREATE TABLE meta.dwattributecolumn (
    dw_attribute_column_id serial4 PRIMARY KEY,
    dw_attribute_column_name character varying,
    dw_attribute_column_datatype character varying
);

-- Створення таблиці для опису атрибутів вимірювань
CREATE TABLE meta.dimensionattributes (
    dimension_id integer,
    dw_attribute_column_id integer,
    FOREIGN KEY (dimension_id) REFERENCES meta.dimension(dimension_id),
    FOREIGN KEY (dw_attribute_column_id) REFERENCES meta.dwattributecolumn(dw_attribute_column_id)

);
-- Створення таблиці для опису фактів
CREATE TABLE meta.fact (
    fact_id serial4 PRIMARY KEY,
    key_name character varying,
    dw_table_id integer,
    fact_type character varying,
    FOREIGN KEY (dw_table_id) REFERENCES meta.dwtable(dw_table_id)
);

-- Створення таблиці для опису метрик фактів
CREATE TABLE meta.factmetric (
    fact_metric_id serial4 PRIMARY KEY,
    fact_id integer,
    dw_attribute_column_id integer,
    fact_metric_description character varying,
    FOREIGN KEY (fact_id) REFERENCES meta.fact(fact_id),
    FOREIGN KEY (dw_attribute_column_id) REFERENCES meta.dwattributecolumn(dw_attribute_column_id)
);

-- Створення таблиці для опису колонок джерела
CREATE TABLE meta.sourcecolumn (
    source_column_id serial4 PRIMARY KEY,
    source_table_id integer,
    source_column_name character varying,
    data_type character varying,
    FOREIGN KEY (source_table_id) REFERENCES meta.sourcetable(source_table_id)
);

-- Створення таблиці для опису трансформацій
CREATE TABLE meta.transformation (
    transformation_id serial4 PRIMARY KEY,
    dw_attribute_column_id integer,
    transformation_rule character varying,
    source_column_id integer,
    FOREIGN KEY (dw_attribute_column_id) REFERENCES meta.dwattributecolumn(dw_attribute_column_id),
    FOREIGN KEY (source_column_id) REFERENCES meta.sourcecolumn(source_column_id)
);
