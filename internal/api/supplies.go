package api

import (
	"database/sql"
	"strconv"

	"github.com/IraIvanishak/service-station-etl/dto"
	"github.com/IraIvanishak/service-station-etl/storage/dbs"
	"github.com/gofiber/fiber/v2"
)

// Supply status ids from public.supply_status.
const (
	supplyPending   = 1
	supplyDelivered = 2
)

func (s *Server) listSupplies(c *fiber.Ctx) error {
	rows, err := s.q.ListSupply(c.UserContext())
	if err != nil {
		return err
	}
	supplies, err := dto.GroupSupplyData(rows)
	if err != nil {
		return err
	}
	return c.JSON(supplies)
}

func (s *Server) createSupply(c *fiber.Ctx) error {
	var req dto.NewSupplyDTO
	if err := parseBody(c, &req); err != nil {
		return err
	}
	ctx := c.UserContext()

	status := int32(supplyDelivered)
	if req.Status == "Pending" {
		status = supplyPending
	}

	var total float64
	for _, d := range req.Details {
		price, _ := s.q.GetSupplyDetailPrice(ctx, d.DetailID)
		total += float64(d.Quantity) * dto.ParseFloat(price.String)
	}

	existing, _ := s.q.ListSupply(ctx)
	supplyID, err := s.q.NewSupply(ctx, dbs.NewSupplyParams{
		TotalSum:     sql.NullString{String: strconv.FormatFloat(total, 'f', 6, 64), Valid: true},
		StatusID:     sql.NullInt32{Int32: status, Valid: true},
		SupplyDate:   req.SupplyDate,
		SupplierID:   sql.NullInt32{Int32: req.SupplierID, Valid: true},
		EmployeeID:   sql.NullInt32{Int32: req.EmployeeID, Valid: true},
		SupplyNumber: dto.CreateSupplyNumber(len(existing) + 1),
	})
	if err != nil {
		return err
	}

	if status == supplyDelivered {
		if _, err := s.q.SetSupplyStatus(ctx, dbs.SetSupplyStatusParams{
			StatusID: sql.NullInt32{Int32: supplyDelivered, Valid: true},
			SupplyID: supplyID,
		}); err != nil {
			return err
		}
	}

	for _, d := range req.Details {
		price, _ := s.q.GetSupplyDetailPrice(ctx, d.DetailID)
		lineID, err := s.q.NewSupplyDetail(ctx, dbs.NewSupplyDetailParams{
			Quantity:     sql.NullInt32{Int32: d.Quantity, Valid: true},
			SupplyID:     sql.NullInt32{Int32: supplyID, Valid: true},
			DetailID:     sql.NullInt32{Int32: d.DetailID, Valid: true},
			PricePerUnit: sql.NullString{String: strconv.FormatFloat(dto.ParseFloat(price.String), 'f', 6, 64), Valid: true},
		})
		if err != nil {
			return err
		}
		inventoryID, err := s.q.StoreInventory(ctx, dbs.StoreInventoryParams{
			Quantity: sql.NullInt32{Int32: d.Quantity, Valid: true},
			DetailID: sql.NullInt32{Int32: d.DetailID, Valid: true},
		})
		if err != nil {
			return err
		}
		if status == supplyDelivered {
			if _, err := s.q.NewSupplyInvoice(ctx, dbs.NewSupplyInvoiceParams{
				InventoryID:     inventoryID,
				SupplyDetailsID: sql.NullInt32{Int32: lineID, Valid: true},
				InvoiceNumber:   sql.NullString{String: dto.CreateInvoiceNumber(), Valid: true},
			}); err != nil {
				return err
			}
		}
	}
	return c.JSON(req)
}

// receiveSupply marks a pending supply as delivered and issues its invoices.
func (s *Server) receiveSupply(c *fiber.Ctx) error {
	id, err := idParam(c, "id")
	if err != nil {
		return err
	}
	ctx := c.UserContext()

	lines, err := s.q.GetSupplyDetailsByID(ctx, sql.NullInt32{Int32: id, Valid: true})
	if err != nil {
		return err
	}
	for _, line := range lines {
		inventoryID, err := s.q.GetInventoryIDbyDetailsID(ctx, line.DetailID)
		if err != nil {
			return err
		}
		if _, err := s.q.NewSupplyInvoice(ctx, dbs.NewSupplyInvoiceParams{
			InventoryID:     inventoryID,
			SupplyDetailsID: sql.NullInt32{Int32: line.SupplyDetailsID, Valid: true},
			InvoiceNumber:   sql.NullString{String: dto.CreateInvoiceNumber(), Valid: true},
		}); err != nil {
			return err
		}
	}

	if _, err := s.q.SetSupplyStatus(ctx, dbs.SetSupplyStatusParams{
		StatusID: sql.NullInt32{Int32: supplyDelivered, Valid: true},
		SupplyID: id,
	}); err != nil {
		return err
	}
	return c.JSON("Updated supply with id: " + strconv.Itoa(int(id)))
}

func (s *Server) deleteSupply(c *fiber.Ctx) error {
	id, err := idParam(c, "id")
	if err != nil {
		return err
	}
	ctx := c.UserContext()
	if err := s.q.DeleteSupplyDetailsBySupplyId(ctx, sql.NullInt32{Int32: id, Valid: true}); err != nil {
		return err
	}
	if err := s.q.DeleteSupply(ctx, id); err != nil {
		return err
	}
	return c.SendString("Deleted")
}
