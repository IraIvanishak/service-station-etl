package dto

import (
	"fmt"
	"time"

	"github.com/IraIvanishak/service-station-etl/storage/dbs"
)

//	{
//		"status": "Pending",
//		"details": [
//		  {
//			"DetailID": 3,
//			"quantity": 4,
//			"StationID": 2
//		  },
//		  {
//			"quantity": 5,
//			"DetailID": 4,
//			"StationID": 3
//		  }
//		],
//		"SupplierID": 5,
//		"supplyDate": "2024-05-15",
//		"EmployeeID": 3
//	  }
type NewSupplyDTO struct {
	Status     string       `json:"status"`
	Details    []DetailsDTO `json:"details"`
	SupplyDate time.Time    `json:"supplyDate"`
	SupplierID int32        `json:"SupplierID"`
	EmployeeID int32        `json:"EmployeeID"`
}

type DetailsDTO struct {
	DetailID  int32 `json:"DetailID"`
	Quantity  int32 `json:"quantity"`
	StationID int32 `json:"StationID"`
}

// GroupSupplyData groups the ListSupply rows (one per supply line) into
// supplies, keeping the query order (by supply id).
func GroupSupplyData(rows []dbs.ListSupplyRow) ([]SupplyDTO, error) {
	var supplies []SupplyDTO
	index := make(map[int32]int)

	for _, row := range rows {
		i, ok := index[row.SupplyID]
		if !ok {
			i = len(supplies)
			index[row.SupplyID] = i
			supplies = append(supplies, SupplyDTO{
				SupplyID:     row.SupplyID,
				SupplyNumber: row.SupplyNumber,
				Details:      []DetailDTO{},
				SupplierName: row.SupplierName.String,
				Employee:     EmployeeDTO{FirstName: row.EmployeeFirstName.String, LastName: row.EmployeeLastName.String},
				SupplyStatus: row.SupplyStatus.String,
			})
		}
		if !row.DetailID.Valid { // supply without lines
			continue
		}
		supplies[i].Details = append(supplies[i].Details, DetailDTO{
			DetailID:       row.DetailID.Int32,
			Quantity:       row.Quantity.Int32,
			PricePerUnit:   ParseFloat(row.PricePerUnit.String),
			Article:        row.Article.Int32,
			SupplyDate:     row.SupplyDate,
			Description:    row.Description.String,
			InvoiceNumber:  row.InvoiceNumber.String,
			ServiceStation: row.ServiceStation.String,
		})
	}
	return supplies, nil
}

type SupplyDTO struct {
	SupplyID     int32       `json:"supply_id"`
	SupplyNumber string      `json:"supply_number"`
	Details      []DetailDTO `json:"details"`
	SupplierName string      `json:"supplier_name"`
	Employee     EmployeeDTO `json:"employee"`
	SupplyStatus string      `json:"supply_status"`
}

type DetailDTO struct {
	DetailID       int32     `json:"detail_id"`
	Quantity       int32     `json:"quantity"`
	PricePerUnit   float64   `json:"price_per_unit"`
	Article        int32     `json:"article"`
	Description    string    `json:"description"`
	InvoiceNumber  string    `json:"invoice_number"`
	ServiceStation string    `json:"service_station"`
	SupplyDate     time.Time `json:"supply_date"`
}

type EmployeeDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// SupplyNumber formats a supply number from its id, e.g. SN-001001.
func SupplyNumber(id int32) string {
	return fmt.Sprintf("SN-%06d", id)
}

// InvoiceNumber formats an invoice number from the supply and supply line
// ids, e.g. INV-1001-2121. It used to be the current Unix second, so every
// line of one supply got the same number.
func InvoiceNumber(supplyID, lineID int32) string {
	return fmt.Sprintf("INV-%d-%d", supplyID, lineID)
}
