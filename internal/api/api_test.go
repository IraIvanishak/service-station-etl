package api_test

// Integration tests for the HTTP API. Each test gets its own copy of the
// seeded database (see internal/testdb).

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/IraIvanishak/service-station-etl/internal/api"
	"github.com/IraIvanishak/service-station-etl/internal/testdb"
	"github.com/gofiber/fiber/v2"
)

// A supply whose second line fails must not leave the supply, its first
// line or a stock change behind.
func TestCreateSupplyIsAllOrNothing(t *testing.T) {
	app, db := newServer(t)
	supplies := testdb.Scalar[int](t, db, "SELECT count(*) FROM supply")
	lines := testdb.Scalar[int](t, db, "SELECT count(*) FROM supply_details")
	stock := testdb.Scalar[int](t, db, "SELECT sum(quantity) FROM inventory")

	code, body := call(t, app, "POST", "/supplies", map[string]any{
		"status":     "Delivered",
		"supplyDate": "2024-05-15T00:00:00Z",
		"SupplierID": 5,
		"EmployeeID": 3,
		"details": []map[string]any{
			{"DetailID": 3, "quantity": 4, "StationID": 3},
			{"DetailID": 999999, "quantity": 5, "StationID": 3},
		},
	})
	if code < 400 {
		t.Fatalf("POST /supplies with an unknown detail: %d %s, want an error", code, body)
	}
	if n := testdb.Scalar[int](t, db, "SELECT count(*) FROM supply"); n != supplies {
		t.Errorf("supply rows: %d, want %d", n, supplies)
	}
	if n := testdb.Scalar[int](t, db, "SELECT count(*) FROM supply_details"); n != lines {
		t.Errorf("supply_details rows: %d, want %d", n, lines)
	}
	if n := testdb.Scalar[int](t, db, "SELECT sum(quantity) FROM inventory"); n != stock {
		t.Errorf("total stock: %d, want %d", n, stock)
	}
}

// A repair whose order insert fails must not leave a car_client behind.
func TestCreateRepairIsAllOrNothing(t *testing.T) {
	app, db := newServer(t)
	carClients := testdb.Scalar[int](t, db, "SELECT count(*) FROM car_client")

	code, body := call(t, app, "POST", "/repairs", map[string]any{
		"assignedEmployeeId": 999999, // no such employee
		"clientId":           1,
		"carId":              301,
		"registrationNumber": "AA0001AA",
	})
	if code < 400 {
		t.Fatalf("POST /repairs with an unknown employee: %d %s, want an error", code, body)
	}
	if n := testdb.Scalar[int](t, db, "SELECT count(*) FROM car_client"); n != carClients {
		t.Errorf("car_client rows: %d, want %d", n, carClients)
	}
}

// --- helpers ---------------------------------------------------------------

func newServer(t *testing.T) (*fiber.App, *sql.DB) {
	t.Helper()
	db := testdb.New(t)
	return api.New(db, "*"), db
}

// call sends a JSON request and returns the status code and body.
func call(t *testing.T, app *fiber.App, method, path string, body any) (int, string) {
	t.Helper()
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		r = bytes.NewReader(b)
	}
	req := httptest.NewRequest(method, path, r)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return resp.StatusCode, string(out)
}
