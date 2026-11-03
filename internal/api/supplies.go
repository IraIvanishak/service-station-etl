package api

import (
	"database/sql"
	"fmt"
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
	status := int32(supplyDelivered)
	if req.Status == "Pending" {
		status = supplyPending
	}
	for i, d := range req.Details {
		if d.StationID == 0 {
			return badRequest("details[%d]: StationID is required", i)
		}
	}

	err := s.inTx(c.UserContext(), func(q *dbs.Queries) error {
		ctx := c.UserContext()
		supplyID, err := q.NewSupply(ctx, dbs.NewSupplyParams{
			TotalSum:   sql.NullString{String: "0", Valid: true}, // summed from the lines below
			StatusID:   sql.NullInt32{Int32: status, Valid: true},
			SupplyDate: req.SupplyDate,
			SupplierID: sql.NullInt32{Int32: req.SupplierID, Valid: true},
			EmployeeID: sql.NullInt32{Int32: req.EmployeeID, Valid: true},
			// Numbered from the id below; the column is NOT NULL.
			SupplyNumber: "",
		})
		if err != nil {
			return err
		}
		if err := q.SetSupplyNumber(ctx, dbs.SetSupplyNumberParams{
			SupplyID:     supplyID,
			SupplyNumber: dto.SupplyNumber(supplyID),
		}); err != nil {
			return err
		}

		for _, d := range req.Details {
			lineID, err := q.NewSupplyDetail(ctx, dbs.NewSupplyDetailParams{
				SupplyID: supplyID,
				DetailID: d.DetailID,
				Quantity: d.Quantity,
			})
			if err != nil {
				return fmt.Errorf("detail %d: %w", d.DetailID, err)
			}
			inventoryID, err := q.StoreInventory(ctx, dbs.StoreInventoryParams{
				DetailID:         sql.NullInt32{Int32: d.DetailID, Valid: true},
				ServiceStationID: sql.NullInt32{Int32: d.StationID, Valid: true},
				Quantity:         sql.NullInt32{Int32: d.Quantity, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("stock for detail %d: %w", d.DetailID, err)
			}
			if status == supplyDelivered {
				if _, err := q.NewSupplyInvoice(ctx, dbs.NewSupplyInvoiceParams{
					InventoryID:     inventoryID,
					SupplyDetailsID: sql.NullInt32{Int32: lineID, Valid: true},
					InvoiceNumber:   sql.NullString{String: dto.InvoiceNumber(supplyID, lineID), Valid: true},
				}); err != nil {
					return err
				}
			}
		}
		return q.UpdateSupplyTotal(ctx, supplyID)
	})
	if err != nil {
		return err
	}
	return c.JSON(req)
}

// receiveSupply marks a pending supply as delivered and issues its invoices.
func (s *Server) receiveSupply(c *fiber.Ctx) error {
	id, err := idParam(c, "id")
	if err != nil {
		return err
	}
	err = s.inTx(c.UserContext(), func(q *dbs.Queries) error {
		ctx := c.UserContext()
		lines, err := q.GetSupplyDetailsByID(ctx, sql.NullInt32{Int32: id, Valid: true})
		if err != nil {
			return err
		}
		for _, line := range lines {
			inventoryID, err := q.GetInventoryIDbyDetailsID(ctx, line.DetailID)
			if err != nil {
				return fmt.Errorf("stock for detail %d: %w", line.DetailID.Int32, err)
			}
			if _, err := q.NewSupplyInvoice(ctx, dbs.NewSupplyInvoiceParams{
				InventoryID:     inventoryID,
				SupplyDetailsID: sql.NullInt32{Int32: line.SupplyDetailsID, Valid: true},
				InvoiceNumber:   sql.NullString{String: dto.InvoiceNumber(id, line.SupplyDetailsID), Valid: true},
			}); err != nil {
				return err
			}
		}
		_, err = q.SetSupplyStatus(ctx, dbs.SetSupplyStatusParams{
			StatusID: sql.NullInt32{Int32: supplyDelivered, Valid: true},
			SupplyID: id,
		})
		return err
	})
	if err != nil {
		return err
	}
	return c.JSON("Updated supply with id: " + strconv.Itoa(int(id)))
}

func (s *Server) deleteSupply(c *fiber.Ctx) error {
	id, err := idParam(c, "id")
	if err != nil {
		return err
	}
	err = s.inTx(c.UserContext(), func(q *dbs.Queries) error {
		ctx := c.UserContext()
		if err := q.DeleteSupplyDetailsBySupplyId(ctx, sql.NullInt32{Int32: id, Valid: true}); err != nil {
			return err
		}
		return deleted(q.DeleteSupply(ctx, id))
	})
	if err != nil {
		return err
	}
	return c.SendString("Deleted")
}
