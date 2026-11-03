package api

import (
	"database/sql"

	"github.com/IraIvanishak/service-station-etl/storage/dbs"
	"github.com/gofiber/fiber/v2"
)

type newClient struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Phone     string `json:"PhoneNumber"`
}

func (s *Server) createClient(c *fiber.Ctx) error {
	var client newClient
	if err := parseBody(c, &client); err != nil {
		return err
	}
	err := s.q.SoreClient(c.UserContext(), dbs.SoreClientParams{
		FirstName:   client.FirstName,
		LastName:    client.LastName,
		PhoneNumber: sql.NullString{String: client.Phone, Valid: true},
	})
	if err != nil {
		return err
	}
	return c.JSON(client)
}
