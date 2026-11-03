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

func GroupSupplyData(rows []dbs.ListSupplyRow) ([]SupplyDTO, error) {
	supplyMap := make(map[int32]*SupplyDTO)

	for _, row := range rows {
		if _, exists := supplyMap[row.SupplyID]; !exists {
			supplyMap[row.SupplyID] = &SupplyDTO{
				SupplyID:     row.SupplyID,
				SupplyNumber: row.SupplyNumber,
				Details:      []DetailDTO{},
				SupplierName: row.SupplierName.String,
				Employee:     EmployeeDTO{FirstName: row.EmployeeFirstName.String, LastName: row.EmployeeLastName.String},
				SupplyStatus: row.SupplyStatus.String,
			}
		}

		price := ParseFloat(row.PricePerUnit.String)
		detail := DetailDTO{
			DetailID:       row.DetailID.Int32,
			Quantity:       row.Quantity.Int32,
			PricePerUnit:   price,
			Article:        row.Article.Int32,
			SupplyDate:     row.SupplyDate,
			Description:    row.Description.String,
			InvoiceNumber:  row.InvoiceNumber.String,
			ServiceStation: row.ServiceStation.String,
		}
		supplyMap[row.SupplyID].Details = append(supplyMap[row.SupplyID].Details, detail)
	}

	supplies := make([]SupplyDTO, 0, len(supplyMap))
	for _, supply := range supplyMap {
		supplies = append(supplies, *supply)
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

// func to generate smt like SN-000900 baced on id
func CreateSupplyNumber(id int) string {
	return fmt.Sprintf("SN-000%d", id)
}

// func to gebereate smt like INV-636217
func CreateInvoiceNumber() string {
	t := time.Now()
	return fmt.Sprintf("INV-%d", t.Unix())
}
