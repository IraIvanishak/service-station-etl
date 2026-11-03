package api

import (
	"database/sql"
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
	ctx := c.UserContext()

	carClientID, err := s.q.StoreCarClient(ctx, dbs.StoreCarClientParams{
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

	err = s.q.StoreOrder(ctx, dbs.StoreOrderParams{
		EmployeeID:      sql.NullInt32{Int32: repair.AssignedEmployeeID, Valid: true},
		StartDate:       sql.NullTime{Time: time.Now(), Valid: true},
		RepairStatusID:  sql.NullInt32{Int32: 1, Valid: true},
		Notes:           sql.NullString{String: repair.Notes, Valid: true},
		RepairNumber:    sql.NullString{String: dto.CreateRepairNumber(), Valid: true},
		PaymentStatusID: sql.NullInt32{Int32: 1, Valid: true},
		CarClientID:     sql.NullInt32{Int32: carClientID, Valid: true},
		TotalSum:        sql.NullString{String: "0", Valid: true},
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
	if err := s.q.DeleteOrderCascade(c.UserContext(), id); err != nil {
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
	ctx := c.UserContext()

	price, _ := s.q.GetServicePriceById(ctx, service.ServiceID)
	unitPrice := dto.ParseFloat(price.String)
	err = s.q.StoreServiceRepair(ctx, dbs.StoreServiceRepairParams{
		ServiceID:      sql.NullInt32{Int32: service.ServiceID, Valid: true},
		Quantity:       sql.NullInt32{Int32: service.Quantity, Valid: true},
		RepairID:       sql.NullInt32{Int32: repairID, Valid: true},
		CompletionDate: sql.NullTime{Time: time.Now(), Valid: true},
		ServiceTotalPrice: sql.NullString{
			String: strconv.FormatFloat(unitPrice*float64(service.Quantity), 'f', 6, 64),
			Valid:  true,
		},
	})
	if err != nil {
		return err
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
