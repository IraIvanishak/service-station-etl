-- name: ClearStorage :exec
CALL storage.clearstorage();

-- name: OltpToStage :exec
CALL stage.transfer_data_to_stage();

-- name: StageToOlap :exec
CALL storage.fill_storage_from_stage();

-- name: GetStageData :many
SELECT 'inventory' AS table_name, COUNT(*) AS row_count FROM stage.inventory
UNION ALL
SELECT 'service_inventory', COUNT(*) FROM stage.service_inventory
UNION ALL
SELECT 'service', COUNT(*) FROM stage.service
UNION ALL
SELECT 'supplier', COUNT(*) FROM stage.supplier
UNION ALL
SELECT 'supply', COUNT(*) FROM stage.supply
UNION ALL
SELECT 'supply_details', COUNT(*) FROM stage.supply_details
UNION ALL
SELECT 'service_repair', COUNT(*) FROM stage.service_repair
UNION ALL
SELECT 'detail', COUNT(*) FROM stage.detail
UNION ALL
SELECT 'repair', COUNT(*) FROM stage.repair;


-- name: OlapToStageIncremental :exec
CALL stage.fill_staging_from_oltp_incremental();

-- name: StageToOlapIncremental :exec
CALL storage.fill_olap_from_staging_incremental();

-- name: GetLastDataLoadHistoryID :one
SELECT data_load_history_id
FROM meta.dataloadhistory
ORDER BY data_load_history_id DESC  
LIMIT 1;


-- name: UpdateDataLoadHistory :exec
UPDATE meta.dataloadhistory
SET 
    load_rows = @load_rows,
    affected_table_count = @affected_table_count
WHERE
    data_load_history_id = @data_load_history_id;