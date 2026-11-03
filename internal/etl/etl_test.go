package etl_test

// Integration tests for the ETL procedures. Each test clones the
// service_station_template database created by db/init/13_test_template.sh,
// so it starts from the seed data and never touches the dev database.
//
//	docker compose up -d --wait
//	make test-integration

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/IraIvanishak/service-station-etl/internal/etl"
	_ "github.com/lib/pq"
)

func TestFullLoadFillsDimensions(t *testing.T) {
	db := testDB(t)
	run(t, etl.NewRunner(db).Full)

	for dim, source := range map[string]string{
		"storage.service_dim":            "public.service",
		"storage.supplier_dim":           "public.supplier",
		"storage.repair_dim":             "public.repair_order",
		"storage.service_repair_dim":     "public.service_repairs",
		"storage.service_inventory_fact": "public.service_inventory",
	} {
		got := scalar[int](t, db, "SELECT count(*) FROM "+dim)
		want := scalar[int](t, db, "SELECT count(*) FROM "+source)
		if got != want {
			t.Errorf("%s has %d rows, want %d (one per row of %s)", dim, got, want, source)
		}
	}
	if n := scalar[int](t, db, "SELECT count(*) FROM meta.dataloadhistory"); n != 1 {
		t.Errorf("meta.dataloadhistory has %d rows after one load, want 1", n)
	}
}

// The load history must record the rows written to all warehouse tables.
func TestFullLoadRecordsTotalRowCount(t *testing.T) {
	db := testDB(t)
	run(t, etl.NewRunner(db).Full)

	want := scalar[int](t, db, `
		SELECT (SELECT count(*) FROM storage.date_dim)
		     + (SELECT count(*) FROM storage.inventory_dim)
		     + (SELECT count(*) FROM storage.service_dim)
		     + (SELECT count(*) FROM storage.supplier_dim)
		     + (SELECT count(*) FROM storage.repair_dim)
		     + (SELECT count(*) FROM storage.service_repair_dim)
		     + (SELECT count(*) FROM storage.service_inventory_fact)
		     + (SELECT count(*) FROM storage.supply_details_fact)`)
	got := scalar[int](t, db, "SELECT load_rows FROM meta.dataloadhistory")
	if got != want {
		t.Errorf("load_rows = %d, want %d (rows in all warehouse tables)", got, want)
	}
}

// In the seed data inventory_id happens to equal detail_id. Stock the same
// detail at a second station so they differ, and use that stock in a repair.
func TestFullLoadKeysPartsUsageByDetail(t *testing.T) {
	db := testDB(t)
	inventoryID := scalar[int](t, db, `
		INSERT INTO public.inventory (detail_id, quantity, service_station_id)
		VALUES (1, 5, 2) RETURNING inventory_id`)
	exec(t, db, `
		INSERT INTO public.service_inventory (service_repair_id, inventory_id, amount)
		VALUES ((SELECT min(service_repairs_id) FROM public.service_repairs), $1, 2)`, inventoryID)

	run(t, etl.NewRunner(db).Full)

	mismatched := scalar[int](t, db, `
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
	db := testDB(t)
	run(t, etl.NewRunner(db).Full)

	want := scalar[int](t, db, `
		SELECT count(DISTINCT (sd.detail_id, s.supplier_id, s.supply_date))
		FROM public.supply_details sd JOIN public.supply s USING (supply_id)`)
	got := scalar[int](t, db, "SELECT count(*) FROM storage.supply_details_fact")
	if got != want {
		t.Errorf("supply_details_fact has %d rows, want %d", got, want)
	}
}

// Quantities and costs in the supply fact must add up to the OLTP totals.
// 277 seed supplies have more than one line, which is what used to break it.
func TestFullLoadSupplyTotalsMatchOLTP(t *testing.T) {
	db := testDB(t)
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
	db := testDB(t)
	loads := etl.NewRunner(db)
	run(t, loads.Full)

	exec(t, db, "UPDATE public.service SET price = 99.50 WHERE service_id = 45")
	exec(t, db, "UPDATE public.service_repairs SET quantity = quantity + 1 WHERE service_repairs_id = 4")
	run(t, loads.Incremental)

	if price := scalar[string](t, db, "SELECT price::text FROM storage.service_dim WHERE bk_service_id = 45"); price != "99.50" {
		t.Errorf("service_dim price = %s, want 99.50", price)
	}
	want := scalar[int](t, db, "SELECT quantity FROM public.service_repairs WHERE service_repairs_id = 4")
	if got := scalar[int](t, db, "SELECT quantity FROM storage.service_repair_dim WHERE bk_service_repair = 4"); got != want {
		t.Errorf("service_repair_dim quantity = %d, want %d", got, want)
	}
	if n := scalar[int](t, db, "SELECT count(*) FROM meta.dataloadhistory"); n != 2 {
		t.Errorf("meta.dataloadhistory has %d rows after two loads, want 2", n)
	}
}

// --- helpers ---------------------------------------------------------------

// testDB returns a connection to a throwaway copy of the seeded database.
func testDB(t *testing.T) *sql.DB {
	t.Helper()
	adminURL := os.Getenv("TEST_DATABASE_URL")
	if adminURL == "" {
		t.Skip("TEST_DATABASE_URL is not set; start the database with `docker compose up -d` and run `make test-integration`")
	}
	admin, err := sql.Open("postgres", adminURL)
	if err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("etl_test_%d", time.Now().UnixNano())
	if _, err := admin.Exec("CREATE DATABASE " + name + " TEMPLATE service_station_template"); err != nil {
		t.Fatalf("clone template database: %v", err)
	}

	u, err := url.Parse(adminURL)
	if err != nil {
		t.Fatal(err)
	}
	u.Path = "/" + name
	db, err := sql.Open("postgres", u.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
		if _, err := admin.Exec("DROP DATABASE IF EXISTS " + name + " WITH (FORCE)"); err != nil {
			t.Errorf("drop %s: %v", name, err)
		}
		admin.Close()
	})

	// The template holds the procedures from when the container was created;
	// load the working-tree versions so SQL edits are tested without a
	// `make db-reset`. (Schema changes still need one.)
	for _, f := range procedureFiles {
		src, err := os.ReadFile(filepath.Join("..", "..", "db", "init", f))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(string(src)); err != nil {
			t.Fatalf("load %s: %v", f, err)
		}
	}
	return db
}

var procedureFiles = []string{
	"06_meta_procedures.sql",
	"07_stage_full_load.sql",
	"08_stage_incremental_load.sql",
	"09_warehouse_full_load.sql",
	"10_warehouse_incremental_load.sql",
	"11_warehouse_clear.sql",
}

func run(t *testing.T, load func(context.Context) error) {
	t.Helper()
	if err := load(context.Background()); err != nil {
		t.Fatalf("load failed: %v", err)
	}
}

func exec(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()
	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("%v\n%s", err, query)
	}
}

func scalar[T any](t *testing.T, db *sql.DB, query string, args ...any) T {
	t.Helper()
	var v T
	if err := db.QueryRow(query, args...).Scan(&v); err != nil {
		t.Fatalf("%v\n%s", err, query)
	}
	return v
}
