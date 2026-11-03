package api

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/IraIvanishak/service-station-etl/dto"
	"github.com/IraIvanishak/service-station-etl/storage/dbs"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) listRepairs(c *fiber.Ctx) error {
	rows, err := s.q.ListOrders(c.UserContext())
	if err != nil {
		return err
	}
	return c.JSON(dto.BuildRepairs(rows))
}

func (s *Server) createRepair(c *fiber.Ctx) error {
	var repair dto.Repair
	if err := parseBody(c, &repair); err != nil {
		return err
	}
	err := s.inTx(c.UserContext(), func(q *dbs.Queries) error {
		ctx := c.UserContext()
		carClientID, err := q.StoreCarClient(ctx, dbs.StoreCarClientParams{
			CarID:              repair.CarID,
			ClientID:           repair.ClientID,
			Description:        sql.NullString{String: repair.Description, Valid: true},
			RegistrationNumber: sql.NullString{String: repair.RegistrationNumber, Valid: true},
			Vin:                sql.NullString{String: repair.VIN, Valid: true},
			Color:              sql.NullString{String: repair.Color, Valid: true},
		})
		if err != nil {
			return err
		}
		opened := time.Now()
		repairID, err := q.StoreOrder(ctx, dbs.StoreOrderParams{
			EmployeeID:      sql.NullInt32{Int32: repair.AssignedEmployeeID, Valid: true},
			StartDate:       sql.NullTime{Time: opened, Valid: true},
			RepairStatusID:  sql.NullInt32{Int32: 1, Valid: true},
			Notes:           sql.NullString{String: repair.Notes, Valid: true},
			PaymentStatusID: sql.NullInt32{Int32: 1, Valid: true},
			CarClientID:     sql.NullInt32{Int32: carClientID, Valid: true},
			TotalSum:        sql.NullString{String: "0", Valid: true},
		})
		if err != nil {
			return err
		}
		return q.SetRepairNumber(ctx, dbs.SetRepairNumberParams{
			RepairID:     repairID,
			RepairNumber: sql.NullString{String: dto.RepairNumber(repairID, opened), Valid: true},
		})
	})
	if err != nil {
		return err
	}
	return c.JSON(repair)
}

func (s *Server) deleteRepair(c *fiber.Ctx) error {
	id, err := idParam(c, "id")
	if err != nil {
		return err
	}
	err = s.inTx(c.UserContext(), func(q *dbs.Queries) error {
		ctx := c.UserContext()
		if err := q.DeleteRepairParts(ctx, sql.NullInt32{Int32: id, Valid: true}); err != nil {
			return err
		}
		if err := q.DeleteRepairServices(ctx, sql.NullInt32{Int32: id, Valid: true}); err != nil {
			return err
		}
		return deleted(q.DeleteOrder(ctx, id))
	})
	if err != nil {
		return err
	}
	return c.SendString("Deleted")
}

func (s *Server) addRepairService(c *fiber.Ctx) error {
	repairID, err := idParam(c, "id")
	if err != nil {
		return err
	}
	var service dto.ServiceRepair
	if err := parseBody(c, &service); err != nil {
		return err
	}
	_, err = s.q.StoreServiceRepair(c.UserContext(), dbs.StoreServiceRepairParams{
		ServiceID:      service.ServiceID,
		Quantity:       service.Quantity,
		RepairID:       repairID,
		CompletionDate: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("service %d: %w", service.ServiceID, err)
	}
	return c.JSON(service)
}

// repairPart is the body of POST /repairsdetails, e.g.
// {"repairId":1505,"serviceName":"Engine Belt Replacement","detailId":14,"quantity":2}
type repairPart struct {
	RepairID    int32  `json:"repairId"`
	DetailId    int32  `json:"detailId"`
	ServiceName string `json:"serviceName"`
	Quantity    int32  `json:"quantity"`
}

func (s *Server) addRepairPart(c *fiber.Ctx) error {
	var part repairPart
	if err := parseBody(c, &part); err != nil {
		return err
	}
	ctx := c.UserContext()

	serviceRepairID, err := s.q.GetServiceRepairByRepairIdAndServiceName(ctx, dbs.GetServiceRepairByRepairIdAndServiceNameParams{
		RepairID:    sql.NullInt32{Int32: part.RepairID, Valid: true},
		ServiceName: part.ServiceName,
	})
	if err != nil {
		return err
	}
	inventoryID, err := s.q.GetInventoryIDbyDetailsID(ctx, sql.NullInt32{Int32: part.DetailId, Valid: true})
	if err != nil {
		return err
	}
	err = s.q.StoreServiceInventory(ctx, dbs.StoreServiceInventoryParams{
		ServiceRepairID: sql.NullInt32{Int32: serviceRepairID, Valid: true},
		InventoryID:     sql.NullInt32{Int32: inventoryID, Valid: true},
		Amount:          sql.NullString{String: strconv.Itoa(int(part.Quantity)), Valid: true},
	})
	if err != nil {
		return err
	}
	return c.JSON(part)
}
