// Package testdb gives integration tests a throwaway copy of the seeded
// database. It clones service_station_template, which
// db/init/13_test_template.sh creates when the compose database starts.
//
//	docker compose up -d --wait
//	make test-integration
package testdb

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// New returns a connection to a fresh copy of the seeded database and drops
// it when the test ends. It skips the test when TEST_DATABASE_URL is unset.
func New(t testing.TB) *sql.DB {
	t.Helper()
	adminURL := os.Getenv("TEST_DATABASE_URL")
	if adminURL == "" {
		t.Skip("TEST_DATABASE_URL is not set; start the database with `docker compose up -d` and run `make test-integration`")
	}
	admin, err := sql.Open("postgres", adminURL)
	if err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("test_%d", time.Now().UnixNano())
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
		src, err := os.ReadFile(filepath.Join(initDir(), f))
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

// initDir is db/init, found relative to this source file.
func initDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "db", "init")
}

// Exec runs a statement and fails the test on error.
func Exec(t testing.TB, db *sql.DB, query string, args ...any) {
	t.Helper()
	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("%v\n%s", err, query)
	}
}

// Scalar returns the single value a query selects.
func Scalar[T any](t testing.TB, db *sql.DB, query string, args ...any) T {
	t.Helper()
	var v T
	if err := db.QueryRow(query, args...).Scan(&v); err != nil {
		t.Fatalf("%v\n%s", err, query)
	}
	return v
}
