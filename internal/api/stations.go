package api

import (
	"database/sql"

	"github.com/IraIvanishak/service-station-etl/dto"
	"github.com/IraIvanishak/service-station-etl/storage/dbs"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) createStation(c *fiber.Ctx) error {
	var station dto.Station
	if err := parseBody(c, &station); err != nil {
		return err
	}
	err := s.q.AddServiceStation(c.UserContext(), dbs.AddServiceStationParams{
		Address: station.Address,
		CityID:  sql.NullInt32{Int32: station.City, Valid: true},
	})
	if err != nil {
		return err
	}
	return c.JSON(station)
}

func (s *Server) updateStation(c *fiber.Ctx) error {
	var station dto.Station
	if err := parseBody(c, &station); err != nil {
		return err
	}
	err := s.q.UpdateServiceStation(c.UserContext(), dbs.UpdateServiceStationParams{
		StationID: station.Id,
		Address:   station.Address,
		CityID:    sql.NullInt32{Int32: station.City, Valid: true},
	})
	if err != nil {
		return err
	}
	return c.JSON(station)
}

func (s *Server) deleteStation(c *fiber.Ctx) error {
	id, err := idParam(c, "id")
	if err != nil {
		return err
	}
	if err := deleted(s.q.DeleteServiceStation(c.UserContext(), id)); err != nil {
		return err
	}
	return c.SendString("Deleted")
}
