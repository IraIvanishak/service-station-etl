package etl_test

// Integration tests for the ETL procedures. Each test runs against its own
// copy of the seeded database (see internal/testdb).

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/IraIvanishak/service-station-etl/internal/etl"
	"github.com/IraIvanishak/service-station-etl/internal/testdb"
	_ "github.com/lib/pq"
)

func TestFullLoadFillsDimensions(t *testing.T) {
	db := testdb.New(t)
	run(t, etl.NewRunner(db).Full)

	for dim, source := range map[string]string{
		"storage.service_dim":            "public.service",
		"storage.supplier_dim":           "public.supplier",
		"storage.repair_dim":             "public.repair_order",
		"storage.service_repair_dim":     "public.service_repairs",
		"storage.service_inventory_fact": "public.service_inventory",
	} {
		got := testdb.Scalar[int](t, db, "SELECT count(*) FROM "+dim)
		want := testdb.Scalar[int](t, db, "SELECT count(*) FROM "+source)
		if got != want {
			t.Errorf("%s has %d rows, want %d (one per row of %s)", dim, got, want, source)
		}
	}
	if n := testdb.Scalar[int](t, db, "SELECT count(*) FROM meta.dataloadhistory"); n != 1 {
		t.Errorf("meta.dataloadhistory has %d rows after one load, want 1", n)
	}
}

// The load history must record the rows written to all warehouse tables.
func TestFullLoadRecordsTotalRowCount(t *testing.T) {
	db := testdb.New(t)
	run(t, etl.NewRunner(db).Full)

	want := testdb.Scalar[int](t, db, `
		SELECT (SELECT count(*) FROM storage.date_dim)
		     + (SELECT count(*) FROM storage.inventory_dim)
		     + (SELECT count(*) FROM storage.service_dim)
		     + (SELECT count(*) FROM storage.supplier_dim)
		     + (SELECT count(*) FROM storage.repair_dim)
		     + (SELECT count(*) FROM storage.service_repair_dim)
		     + (SELECT count(*) FROM storage.service_inventory_fact)
		     + (SELECT count(*) FROM storage.supply_details_fact)`)
	got := testdb.Scalar[int](t, db, "SELECT load_rows FROM meta.dataloadhistory")
	if got != want {
		t.Errorf("load_rows = %d, want %d (rows in all warehouse tables)", got, want)
	}
}

// In the seed data inventory_id happens to equal detail_id. Stock the same
// detail at a second station so they differ, and use that stock in a repair.
func TestFullLoadKeysPartsUsageByDetail(t *testing.T) {
	db := testdb.New(t)
	inventoryID := testdb.Scalar[int](t, db, `
		INSERT INTO public.inventory (detail_id, quantity, service_station_id)
		VALUES (1, 5, 2) RETURNING inventory_id`)
	testdb.Exec(t, db, `
		INSERT INTO public.service_inventory (service_repair_id, inventory_id, amount)
		VALUES ((SELECT min(service_repairs_id) FROM public.service_repairs), $1, 2)`, inventoryID)

	run(t, etl.NewRunner(db).Full)

	mismatched := testdb.Scalar[int](t, db, `
		SELECT count(*)
		FROM storage.service_inventory_fact f
		JOIN public.service_inventory si ON si.service_inventory_id = f.bk_service_inventory_id
		JOIN public.inventory i ON i.inventory_id = si.inventory_id
		JOIN storage.inventory_dim d ON d.inventory_dim_id = f.inventory_dim_id
		WHERE d.bk_detail_id <> i.detail_id`)
	if mismatched != 0 {
		t.Errorf("%d parts-usage facts point at the wrong detail", mismatched)
	}
}

// Every (detail, supplier, supply date) group in the OLTP data must become
// one supply fact row, including details that are not in stock anywhere.
func TestFullLoadKeepsEverySupplyGroup(t *testing.T) {
	db := testdb.New(t)
	run(t, etl.NewRunner(db).Full)

	want := testdb.Scalar[int](t, db, `
		SELECT count(DISTINCT (sd.detail_id, s.supplier_id, s.supply_date))
		FROM public.supply_details sd JOIN public.supply s USING (supply_id)`)
	got := testdb.Scalar[int](t, db, "SELECT count(*) FROM storage.supply_details_fact")
	if got != want {
		t.Errorf("supply_details_fact has %d rows, want %d", got, want)
	}
}

// Quantities and costs in the supply fact must add up to the OLTP totals.
// 277 seed supplies have more than one line, which is what used to break it.
func TestFullLoadSupplyTotalsMatchOLTP(t *testing.T) {
	db := testdb.New(t)
	run(t, etl.NewRunner(db).Full)

	var wantQty, gotQty int
	var wantCost, gotCost string
	if err := db.QueryRow(`
		SELECT sum(sd.quantity), sum(sd.quantity * sd.price_per_unit)::text
		FROM public.supply_details sd JOIN public.supply s USING (supply_id)`).Scan(&wantQty, &wantCost); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`
		SELECT sum(detail_quantity), sum(detail_total_cost)::text
		FROM storage.supply_details_fact`).Scan(&gotQty, &gotCost); err != nil {
		t.Fatal(err)
	}
	if gotQty != wantQty || gotCost != wantCost {
		t.Errorf("warehouse totals: %d units / %s, OLTP: %d units / %s", gotQty, gotCost, wantQty, wantCost)
	}
}

func TestIncrementalLoadAppliesChanges(t *testing.T) {
	db := testdb.New(t)
	loads := etl.NewRunner(db)
	run(t, loads.Full)

	testdb.Exec(t, db, "UPDATE public.service SET price = 99.50 WHERE service_id = 45")
	testdb.Exec(t, db, "UPDATE public.service_repairs SET quantity = quantity + 1 WHERE service_repairs_id = 4")
	run(t, loads.Incremental)

	if price := testdb.Scalar[string](t, db, "SELECT price::text FROM storage.service_dim WHERE bk_service_id = 45"); price != "99.50" {
		t.Errorf("service_dim price = %s, want 99.50", price)
	}
	want := testdb.Scalar[int](t, db, "SELECT quantity FROM public.service_repairs WHERE service_repairs_id = 4")
	if got := testdb.Scalar[int](t, db, "SELECT quantity FROM storage.service_repair_dim WHERE bk_service_repair = 4"); got != want {
		t.Errorf("service_repair_dim quantity = %d, want %d", got, want)
	}
	if n := testdb.Scalar[int](t, db, "SELECT count(*) FROM meta.dataloadhistory"); n != 2 {
		t.Errorf("meta.dataloadhistory has %d rows after two loads, want 2", n)
	}
}

// A repair is created before any service is added to it (POST /repairs).
func TestIncrementalLoadPicksUpRepairWithoutServices(t *testing.T) {
	db := testdb.New(t)
	loads := etl.NewRunner(db)
	run(t, loads.Full)

	repairID := testdb.Scalar[int](t, db, `
		INSERT INTO public.repair_order (employee_id, start_date, repair_status_id, notes, repair_number, payment_status_id, total_sum)
		VALUES (1, now(), 1, 'New repair', 'RO-TEST', 1, 0) RETURNING repair_id`)
	run(t, loads.Incremental)

	if n := testdb.Scalar[int](t, db, "SELECT count(*) FROM storage.repair_dim WHERE bk_repair_id = $1", repairID); n != 1 {
		t.Errorf("repair %d is not in repair_dim after the incremental load", repairID)
	}
}

// After an incremental load the warehouse must look exactly like a full
// rebuild from the same OLTP data.
func TestIncrementalLoadMatchesFullLoad(t *testing.T) {
	db := testdb.New(t)
	loads := etl.NewRunner(db)
	run(t, loads.Full)

	// 1. A repair uses a part whose stock and detail rows did not change.
	testdb.Exec(t, db, `
		INSERT INTO public.service_inventory (service_repair_id, inventory_id, amount)
		VALUES (4, 1, 3)`)

	// 2. A second delivery from the same supplier on the same day for a detail
	//    that already has a line: the fact row must hold both lines.
	testdb.Exec(t, db, `
		WITH line AS (
			SELECT sd.detail_id, sd.price_per_unit, s.supplier_id, s.supply_date, s.employee_id, s.status_id
			FROM public.supply_details sd JOIN public.supply s USING (supply_id)
			ORDER BY sd.supply_details_id LIMIT 1),
		new_supply AS (
			INSERT INTO public.supply (supply_number, total_sum, supply_date, status_id, supplier_id, employee_id)
			SELECT 'SN-TEST', 5 * price_per_unit, supply_date, status_id, supplier_id, employee_id FROM line
			RETURNING supply_id)
		INSERT INTO public.supply_details (detail_id, quantity, price_per_unit, supply_id)
		SELECT line.detail_id, 5, line.price_per_unit, new_supply.supply_id FROM line, new_supply`)

	// 3. One service line of a multi-service repair changes; the estimated
	//    duration of the whole repair, and every parts fact of it, changes.
	testdb.Exec(t, db, `
		UPDATE public.service_repairs SET quantity = quantity + 2
		WHERE service_repairs_id = (
			SELECT min(service_repairs_id) FROM public.service_repairs
			WHERE repair_id IN (SELECT repair_id FROM public.service_repairs GROUP BY repair_id HAVING count(*) > 1))`)

	run(t, loads.Incremental)
	incremental := warehouseSnapshot(t, db)
	run(t, loads.Full)
	full := warehouseSnapshot(t, db)

	for table, want := range full {
		if diff := firstDiff(incremental[table], want); diff != "" {
			t.Errorf("%s differs after incremental load (-incremental +full):\n%s", table, diff)
		}
	}
}

// With no load history yet, an incremental load has to take everything.
func TestIncrementalLoadWithoutHistoryLoadsEverything(t *testing.T) {
	db := testdb.New(t)
	loads := etl.NewRunner(db)

	run(t, loads.Incremental)
	incremental := warehouseSnapshot(t, db)
	run(t, loads.Full)
	full := warehouseSnapshot(t, db)

	for table, want := range full {
		if diff := firstDiff(incremental[table], want); diff != "" {
			t.Errorf("%s differs (-incremental +full):\n%s", table, diff)
		}
	}
}

// warehouseSnapshot renders every warehouse table by business keys and
// values, leaving out surrogate ids that differ between load strategies.
func warehouseSnapshot(t *testing.T, db *sql.DB) map[string][]string {
	t.Helper()
	queries := map[string]string{
		"service_dim":  `SELECT bk_service_id, name, price, estimated_duration FROM storage.service_dim`,
		"supplier_dim": `SELECT bk_supplier_id, name, phone_number, address FROM storage.supplier_dim`,
		"inventory_dim": `SELECT bk_detail_id, description, article, price, dimensions, weight
			FROM storage.inventory_dim`,
		"repair_dim": `SELECT r.bk_repair_id, r.notes, r.repair_number, make_date(d.year, d.month, d.day)
			FROM storage.repair_dim r LEFT JOIN storage.date_dim d ON d.date_id = r.start_repair_date_dim_id`,
		"service_repair_dim": `SELECT sr.bk_service_repair, s.bk_service_id, r.bk_repair_id, sr.quantity,
				sr.service_total_price, make_date(d.year, d.month, d.day)
			FROM storage.service_repair_dim sr
			JOIN storage.service_dim s ON s.service_dim_id = sr.service_id
			JOIN storage.repair_dim r ON r.repair_dim_id = sr.repair_id
			LEFT JOIN storage.date_dim d ON d.date_id = sr.complete_date_dim_id`,
		"service_inventory_fact": `SELECT f.bk_service_inventory_id, sr.bk_service_repair, i.bk_detail_id,
				f.parts_cost, f.total_cost, f.parts_cost_percentage, f.actual_service_duration,
				f.estimated_service_duration, f.service_duration_deviation, f.amount
			FROM storage.service_inventory_fact f
			JOIN storage.service_repair_dim sr ON sr.service_repair_dim_id = f.service_repair_dim_id
			JOIN storage.inventory_dim i ON i.inventory_dim_id = f.inventory_dim_id`,
		"supply_details_fact": `SELECT i.bk_detail_id, s.bk_supplier_id,
				make_date(sd.year, sd.month, sd.day), make_date(ed.year, ed.month, ed.day),
				f.detail_quantity, f.detail_total_cost
			FROM storage.supply_details_fact f
			JOIN storage.inventory_dim i ON i.inventory_dim_id = f.detail_dim_id
			JOIN storage.supplier_dim s ON s.supplier_dim_id = f.supplier_id
			JOIN storage.date_dim sd ON sd.date_id = f.start_date_dim_id
			JOIN storage.date_dim ed ON ed.date_id = f.end_date_dim_id`,
	}
	snap := make(map[string][]string, len(queries))
	for table, q := range queries {
		rows, err := db.Query(`SELECT t::text FROM (` + q + `) t ORDER BY 1`)
		if err != nil {
			t.Fatalf("%s: %v", table, err)
		}
		for rows.Next() {
			var line string
			if err := rows.Scan(&line); err != nil {
				t.Fatal(err)
			}
			snap[table] = append(snap[table], line)
		}
		if err := rows.Close(); err != nil {
			t.Fatal(err)
		}
	}
	return snap
}

// firstDiff returns up to five lines present in only one of the sorted slices.
func firstDiff(got, want []string) string {
	var out []string
	i, j := 0, 0
	for (i < len(got) || j < len(want)) && len(out) < 5 {
		switch {
		case j == len(want) || (i < len(got) && got[i] < want[j]):
			out = append(out, "- "+got[i])
			i++
		case i == len(got) || want[j] < got[i]:
			out = append(out, "+ "+want[j])
			j++
		default:
			i++
			j++
		}
	}
	if len(out) == 0 {
		return ""
	}
	return fmt.Sprintf("%s\n(%d rows incremental, %d rows full)", strings.Join(out, "\n"), len(got), len(want))
}

// --- helpers ---------------------------------------------------------------

func run(t *testing.T, load func(context.Context) error) {
	t.Helper()
	if err := load(context.Background()); err != nil {
		t.Fatalf("load failed: %v", err)
	}
}
