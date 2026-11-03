CREATE OR REPLACE PROCEDURE meta.initialize_metadata()
LANGUAGE plpgsql
AS $$
BEGIN
    -- Deleting data from tables and resetting identities within the 'meta' schema
    DELETE FROM meta.dataloadhistory;
    ALTER SEQUENCE meta.dataloadhistory_data_load_history_id_seq RESTART WITH 1;

    DELETE FROM meta.transformation;
    ALTER SEQUENCE meta.transformation_transformation_id_seq RESTART WITH 1;

    DELETE FROM meta.factmetric;
    ALTER SEQUENCE meta.factmetric_fact_metric_id_seq RESTART WITH 1;

    DELETE FROM meta.dimensionattributes;

    DELETE FROM meta.fact; 
    ALTER SEQUENCE meta.fact_fact_id_seq RESTART WITH 1;

    DELETE FROM meta.dimension;
    ALTER SEQUENCE meta.dimension_dimension_id_seq RESTART WITH 1;

    DELETE FROM meta.sourcecolumn;
    ALTER SEQUENCE meta.sourcecolumn_source_column_id_seq RESTART WITH 1;

    DELETE FROM meta.sourcetable;
    ALTER SEQUENCE meta.sourcetable_source_table_id_seq RESTART WITH 1;

    DELETE FROM meta.dwattributecolumn;
    ALTER SEQUENCE meta.dwattributecolumn_dw_attribute_column_id_seq RESTART WITH 1;

    DELETE FROM meta.dwtable;
    ALTER SEQUENCE meta.dwtable_dw_table_id_seq RESTART WITH 1;

    DELETE FROM meta.sourcedb ;
    ALTER SEQUENCE meta.sourcedb_source_db_id_seq RESTART WITH 1;

    -- Insert into source database table
    INSERT INTO meta.sourcedb(source_db_id, source_db_name)
    VALUES (1, 'public');

    -- Insert into source table information with primary key lookup
    INSERT INTO meta.sourcetable(source_db_id, source_table_name, key_name)
    SELECT 
        1, 
        t.table_name,
        (
            SELECT kcu.column_name AS key_name
            FROM information_schema.key_column_usage AS kcu
            JOIN information_schema.table_constraints AS tc 
            ON kcu.constraint_name = tc.constraint_name
            AND tc.constraint_type = 'PRIMARY KEY'
            WHERE kcu.table_name = t.table_name
            AND kcu.table_schema = 'public'
            LIMIT 1
        ) AS key_name
    FROM information_schema.tables AS t
    WHERE t.table_schema = 'public'; 

    -- Insert into source column table
    INSERT INTO meta.sourcecolumn(source_table_id, source_column_name, data_type)
    SELECT
        (SELECT source_table_id FROM meta.sourcetable WHERE source_table_name = t.table_name),
        t.column_name,
        t.data_type
    FROM information_schema.columns AS t
    WHERE t.table_schema = 'public';

    -- Insert into DW attribute column table
    INSERT INTO meta.dwattributecolumn(dw_attribute_column_name, dw_attribute_column_datatype)
    SELECT
        t.column_name,
        t.data_type
    FROM information_schema.columns AS t
    WHERE t.table_schema = 'storage';

    -- Insert into DW table
    INSERT INTO meta.dwtable(dw_table_name)
    SELECT
        t.table_name
    FROM information_schema.tables AS t
    WHERE t.table_schema = 'storage';

    -- Insert into fact table
    INSERT INTO meta.fact(key_name, dw_table_id, fact_type)
    SELECT
        (
            SELECT kcu.column_name
            FROM information_schema.key_column_usage AS kcu
            JOIN information_schema.table_constraints AS tc 
            ON kcu.constraint_name = tc.constraint_name
            AND tc.constraint_type = 'PRIMARY KEY'
            WHERE kcu.table_name = t.table_name
            AND kcu.table_schema = 'storage'
            LIMIT 1
        ) AS key_name,
        (SELECT dw_table_id FROM meta.dwtable WHERE dw_table_name = t.table_name) AS dw_table_id,
        'Type description here'
    FROM information_schema.tables AS t
    WHERE t.table_schema = 'storage' AND t.table_name LIKE '%fact';

    -- Insert into dimension table
    INSERT INTO meta.dimension(key_name, dw_table_id, dimension_name, dimension_type)
    SELECT
        (
            SELECT kcu.column_name
            FROM information_schema.key_column_usage AS kcu
            JOIN information_schema.table_constraints AS tc 
            ON kcu.constraint_name = tc.constraint_name
            AND tc.constraint_type = 'PRIMARY KEY'
            WHERE kcu.table_name = t.table_name
            AND kcu.table_schema = 'storage'
            LIMIT 1
        ) AS key_name,
        (SELECT dw_table_id FROM meta.dwtable WHERE dw_table_name = t.table_name) AS dw_table_id,
        t.table_name AS dimension_name,
        'Dimension type here'
    FROM information_schema.tables AS t
    WHERE t.table_schema = 'storage' AND t.table_name LIKE '%dim';

    -- Insert into dimension attributes table
    INSERT INTO meta.dimensionattributes(dimension_id, dw_attribute_column_id)
    SELECT
        (
            SELECT dimension_id FROM meta.dimension 
            WHERE dimension_name = t.table_name
            LIMIT 1  -- Ensure only one row is returned
        ),
        (
            SELECT dw_attribute_column_id FROM meta.dwattributecolumn 
            WHERE dw_attribute_column_name = t.column_name
            LIMIT 1  -- Ensure only one row is returned
        )
    FROM information_schema.columns AS t
    WHERE t.table_schema = 'storage' AND t.table_name LIKE '%dim';


    -- Insert into fact metric table
    INSERT INTO meta.factmetric(fact_id, dw_attribute_column_id, fact_metric_description)
    SELECT
        (
            SELECT fact_id FROM meta.fact f
            JOIN meta.dwtable dwt ON f.dw_table_id = dwt.dw_table_id
            WHERE dwt.dw_table_name = t.table_name
            LIMIT 1  -- Added LIMIT 1 to ensure only one row is returned
        ),
        (
            SELECT dw_attribute_column_id FROM meta.dwattributecolumn
            WHERE dw_attribute_column_name = t.column_name
            LIMIT 1  -- Added LIMIT 1 to ensure only one row is returned
        ),
        'Metric description here'
    FROM information_schema.columns AS t
    WHERE t.table_schema = 'storage' AND t.table_name LIKE '%fact';


    -- Insert into transformation table
    INSERT INTO meta.transformation(dw_attribute_column_id, source_column_id, transformation_rule)
    SELECT
        dwac.dw_attribute_column_id,
        sc.source_column_id,
        'Transformation rule description here'
    FROM meta.dwattributecolumn AS dwac
    JOIN meta.sourcecolumn AS sc ON dwac.dw_attribute_column_name = sc.source_column_name;

    -- Optionally log the completion of initialization
    RAISE NOTICE 'Metadata initialization completed successfully.';
END;
$$;
