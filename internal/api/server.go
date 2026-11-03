// Package api is the HTTP API used by the Vue front end (service_station/).
package api

import (
	"context"
	"database/sql"

	"github.com/IraIvanishak/service-station-etl/internal/etl"
	"github.com/IraIvanishak/service-station-etl/storage/dbs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Server struct {
	db    *sql.DB
	q     *dbs.Queries
	loads *etl.Runner
}

// New returns the Fiber app with every route registered.
func New(db *sql.DB, corsOrigins string) *fiber.App {
	s := &Server{db: db, q: dbs.New(db), loads: etl.NewRunner(db)}

	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})
	app.Use(cors.New(cors.Config{
		AllowOrigins: corsOrigins,
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))
	s.routes(app)
	return app
}

func (s *Server) routes(app *fiber.App) {
	// Reference data
	app.Get("/city", list(s.q.ListCities))
	app.Get("/services", list(s.q.ListServices))
	app.Get("/details", list(s.q.ListDetails))
	app.Get("/suppliers", list(s.q.ListSuppliers))
	app.Get("/roles", list(s.q.ListEmployeesRoles))
	app.Get("/cars", list(s.q.ListCars))

	// Stations
	app.Get("/stations", list(s.q.ListServiceStations))
	app.Post("/stations", s.createStation)
	app.Put("/stations", s.updateStation)
	app.Delete("/stations/:id", s.deleteStation)

	// People
	app.Get("/employees", list(s.q.ListEmployees))
	app.Get("/employees-full", list(s.q.ListEmployees)) // Employee.vue; same data as /employees
	app.Get("/clients", list(s.q.ListClients))
	app.Post("/client", s.createClient)
	app.Get("/carclient", list(s.q.ListCarClients))

	// Repairs
	app.Get("/repairs", s.listRepairs)
	app.Post("/repairs", s.createRepair)
	app.Delete("/repairs/:id", s.deleteRepair)
	app.Post("/repairs/:id/services", s.addRepairService)
	app.Post("/repairsdetails", s.addRepairPart)

	// Supplies
	app.Get("/supplies", s.listSupplies)
	app.Post("/supplies", s.createSupply)
	app.Patch("/supplies/:id", s.receiveSupply)
	app.Delete("/supplies/:id", s.deleteSupply)

	// Warehouse and ETL
	app.Get("/meta", list(s.q.MetaData))
	app.Get("/service-invenroty-cube", s.serviceInventoryCube)
	app.Get("/supply-cube", s.supplyCube)
	app.Post("/init", s.fullLoad)
	app.Get("/stage", s.detectChanges)
	app.Post("/incremental", s.incrementalLoad)
	app.Post("/clearmeta", s.clearWarehouse)
}

// list serves the result of a query that takes no arguments.
func list[T any](query func(context.Context) ([]T, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := query(c.UserContext())
		if err != nil {
			return err
		}
		return c.JSON(rows)
	}
}

// inTx runs fn in one transaction, so a request that writes several rows
// either writes all of them or none.
func (s *Server) inTx(ctx context.Context, fn func(q *dbs.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(s.q.WithTx(tx)); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
