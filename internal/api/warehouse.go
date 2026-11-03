package api

import (
	"github.com/IraIvanishak/service-station-etl/dto"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) serviceInventoryCube(c *fiber.Ctx) error {
	rows, err := s.q.GetCubeServiceInventory(c.UserContext())
	if err != nil {
		return err
	}
	out := make([]dto.ServiceInventoryDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, dto.ConvertToUserFriendly(row))
	}
	return c.JSON(out)
}

func (s *Server) supplyCube(c *fiber.Ctx) error {
	rows, err := s.q.GetCubeSupply(c.UserContext())
	if err != nil {
		return err
	}
	out := make([]dto.SupplyDetailsDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, dto.ToDTOSupply(row))
	}
	return c.JSON(out)
}

func (s *Server) fullLoad(c *fiber.Ctx) error {
	if err := s.loads.Full(c.UserContext()); err != nil {
		return err
	}
	return c.SendString("Initialized")
}

func (s *Server) detectChanges(c *fiber.Ctx) error {
	staged, err := s.loads.DetectChanges(c.UserContext())
	if err != nil {
		return err
	}
	return c.JSON(staged)
}

func (s *Server) incrementalLoad(c *fiber.Ctx) error {
	if err := s.loads.Incremental(c.UserContext()); err != nil {
		return err
	}
	return c.SendString("Incremental")
}

func (s *Server) clearWarehouse(c *fiber.Ctx) error {
	if err := s.loads.Clear(c.UserContext()); err != nil {
		return err
	}
	return c.SendString("Cleared")
}
