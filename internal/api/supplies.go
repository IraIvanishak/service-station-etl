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

	err := s.inTx(c.UserContext(), func(q *dbs.Queries) error {
		ctx := c.UserContext()
		prices := make([]float64, len(req.Details))
		var total float64
		for i, d := range req.Details {
			price, err := q.GetSupplyDetailPrice(ctx, d.DetailID)
			if err != nil {
				return fmt.Errorf("detail %d: %w", d.DetailID, err)
			}
			prices[i] = dto.ParseFloat(price.String)
			total += float64(d.Quantity) * prices[i]
		}

		existing, err := q.ListSupply(ctx)
		if err != nil {
			return err
		}
		supplyID, err := q.NewSupply(ctx, dbs.NewSupplyParams{
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

		for i, d := range req.Details {
			lineID, err := q.NewSupplyDetail(ctx, dbs.NewSupplyDetailParams{
				Quantity:     sql.NullInt32{Int32: d.Quantity, Valid: true},
				SupplyID:     sql.NullInt32{Int32: supplyID, Valid: true},
				DetailID:     sql.NullInt32{Int32: d.DetailID, Valid: true},
				PricePerUnit: sql.NullString{String: strconv.FormatFloat(prices[i], 'f', 6, 64), Valid: true},
			})
			if err != nil {
				return err
			}
			inventoryID, err := q.StoreInventory(ctx, dbs.StoreInventoryParams{
				Quantity: sql.NullInt32{Int32: d.Quantity, Valid: true},
				DetailID: sql.NullInt32{Int32: d.DetailID, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("stock for detail %d: %w", d.DetailID, err)
			}
			if status == supplyDelivered {
				if _, err := q.NewSupplyInvoice(ctx, dbs.NewSupplyInvoiceParams{
					InventoryID:     inventoryID,
					SupplyDetailsID: sql.NullInt32{Int32: lineID, Valid: true},
					InvoiceNumber:   sql.NullString{String: dto.CreateInvoiceNumber(), Valid: true},
				}); err != nil {
					return err
				}
			}
		}
		return nil
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
				InvoiceNumber:   sql.NullString{String: dto.CreateInvoiceNumber(), Valid: true},
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
		return q.DeleteSupply(ctx, id)
	})
	if err != nil {
		return err
	}
	return c.SendString("Deleted")
}
