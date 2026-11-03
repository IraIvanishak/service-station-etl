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

// New repair, supply and invoice numbers must not repeat. (The seed data
// already has duplicate supply numbers, so only new rows are checked.)
func TestDocumentNumbersAreUnique(t *testing.T) {
	app, db := newServer(t)
	lastRepair := testdb.Scalar[int](t, db, "SELECT max(repair_id) FROM repair_order")
	lastSupply := testdb.Scalar[int](t, db, "SELECT max(supply_id) FROM supply")
	lastInvoice := testdb.Scalar[int](t, db, "SELECT max(supply_invoice_id) FROM supply_invoice")

	for i := 0; i < 2; i++ {
		mustCall(t, app, "POST", "/repairs", map[string]any{
			"assignedEmployeeId": 1, "clientId": 1, "carId": 301, "registrationNumber": "AA0001AA",
		})
		mustCall(t, app, "POST", "/supplies", map[string]any{
			"status": "Delivered", "supplyDate": "2024-05-15T00:00:00Z", "SupplierID": 5, "EmployeeID": 3,
			"details": []map[string]any{
				{"DetailID": 3, "quantity": 1, "StationID": 3},
				{"DetailID": 4, "quantity": 1, "StationID": 4},
			},
		})
	}

	for _, c := range []struct {
		table, id, number string
		after             int
	}{
		{"repair_order", "repair_id", "repair_number", lastRepair},
		{"supply", "supply_id", "supply_number", lastSupply},
		{"supply_invoice", "supply_invoice_id", "invoice_number", lastInvoice},
	} {
		dup := testdb.Scalar[int](t, db, `
			SELECT count(*) FROM `+c.table+` n
			WHERE n.`+c.id+` > $1 AND (n.`+c.number+` IS NULL OR EXISTS (
				SELECT 1 FROM `+c.table+` o WHERE o.`+c.number+` = n.`+c.number+` AND o.`+c.id+` <> n.`+c.id+`))`, c.after)
		if dup != 0 {
			t.Errorf("%s: %d new rows have a missing or duplicate %s", c.table, dup, c.number)
		}
	}
}

// Delivered parts go to the station named in the request, and only there.
func TestSupplyStocksTheChosenStation(t *testing.T) {
	app, db := newServer(t)
	// Detail 3 is stocked at station 3 only; stock it at station 2 as well.
	testdb.Exec(t, db, "INSERT INTO inventory (detail_id, quantity, service_station_id) VALUES (3, 10, 2)")
	stock := func(station int) int {
		return testdb.Scalar[int](t, db, `
			SELECT coalesce(sum(quantity), 0) FROM inventory WHERE detail_id = 3 AND service_station_id = $1`, station)
	}
	before := map[int]int{1: stock(1), 2: stock(2), 3: stock(3)}

	mustCall(t, app, "POST", "/supplies", map[string]any{
		"status": "Delivered", "supplyDate": "2024-05-15T00:00:00Z", "SupplierID": 5, "EmployeeID": 3,
		"details": []map[string]any{
			{"DetailID": 3, "quantity": 4, "StationID": 2},
			{"DetailID": 3, "quantity": 6, "StationID": 1}, // not stocked there yet
		},
	})

	want := map[int]int{1: before[1] + 6, 2: before[2] + 4, 3: before[3]}
	for station, w := range want {
		if got := stock(station); got != w {
			t.Errorf("detail 3 at station %d: %d units, want %d", station, got, w)
		}
	}
	if code, body := call(t, app, "POST", "/supplies", map[string]any{
		"status": "Delivered", "supplyDate": "2024-05-15T00:00:00Z", "SupplierID": 5, "EmployeeID": 3,
		"details": []map[string]any{{"DetailID": 3, "quantity": 1}},
	}); code != 400 {
		t.Errorf("line without StationID: %d %s, want 400", code, body)
	}
}

// --- helpers ---------------------------------------------------------------

func newServer(t *testing.T) (*fiber.App, *sql.DB) {
	t.Helper()
	db := testdb.New(t)
	return api.New(db, "*"), db
}

// mustCall is call for requests that have to succeed.
func mustCall(t *testing.T, app *fiber.App, method, path string, body any) string {
	t.Helper()
	code, out := call(t, app, method, path, body)
	if code >= 300 {
		t.Fatalf("%s %s: %d %s", method, path, code, out)
	}
	return out
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
