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
	return db
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
