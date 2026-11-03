// Package etl drives the warehouse loads implemented as stored procedures
// (db/init/07..11). Every load runs in a single transaction, so a failure
// leaves the stage, warehouse and metadata schemas as they were.
package etl

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/IraIvanishak/service-station-etl/storage/dbs"
)

type Runner struct {
	db *sql.DB
}

func NewRunner(db *sql.DB) *Runner {
	return &Runner{db: db}
}

// Full rebuilds the warehouse from the OLTP tables: OLTP -> stage -> storage.
func (r *Runner) Full(ctx context.Context) error {
	return r.inTx(ctx, func(q *dbs.Queries) error {
		if err := q.OltpToStage(ctx); err != nil {
			return fmt.Errorf("stage OLTP data: %w", err)
		}
		if err := q.StageToOlap(ctx); err != nil {
			return fmt.Errorf("load warehouse: %w", err)
		}
		return nil
	})
}

// DetectChanges stages the OLTP rows changed since the last load and returns
// how many rows each stage table received.
func (r *Runner) DetectChanges(ctx context.Context) ([]dbs.GetStageDataRow, error) {
	var rows []dbs.GetStageDataRow
	err := r.inTx(ctx, func(q *dbs.Queries) error {
		if err := q.OlapToStageIncremental(ctx); err != nil {
			return fmt.Errorf("stage changed rows: %w", err)
		}
		var err error
		rows, err = q.GetStageData(ctx)
		return err
	})
	return rows, err
}

// Incremental loads the OLTP rows changed since the last load and records
// how many changed rows and tables it saw in meta.dataloadhistory.
func (r *Runner) Incremental(ctx context.Context) error {
	return r.inTx(ctx, func(q *dbs.Queries) error {
		if err := q.OlapToStageIncremental(ctx); err != nil {
			return fmt.Errorf("stage changed rows: %w", err)
		}
		staged, err := q.GetStageData(ctx)
		if err != nil {
			return fmt.Errorf("count staged rows: %w", err)
		}
		var changedRows, changedTables int32
		for _, t := range staged {
			changedRows += int32(t.RowCount)
			if t.RowCount > 0 {
				changedTables++
			}
		}

		if err := q.StageToOlapIncremental(ctx); err != nil {
			return fmt.Errorf("load warehouse: %w", err)
		}

		loadID, err := q.GetLastDataLoadHistoryID(ctx)
		if err != nil {
			return fmt.Errorf("find load history entry: %w", err)
		}
		return q.UpdateDataLoadHistory(ctx, dbs.UpdateDataLoadHistoryParams{
			LoadRows:           sql.NullInt32{Int32: changedRows, Valid: true},
			AffectedTableCount: sql.NullInt32{Int32: changedTables, Valid: true},
			DataLoadHistoryID:  loadID,
		})
	})
}

// Clear empties every warehouse table.
func (r *Runner) Clear(ctx context.Context) error {
	return r.inTx(ctx, func(q *dbs.Queries) error {
		return q.ClearStorage(ctx)
	})
}

func (r *Runner) inTx(ctx context.Context, fn func(q *dbs.Queries) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(dbs.New(tx)); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
