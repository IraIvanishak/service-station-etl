package dto

import (
	"fmt"
	"time"

	"github.com/IraIvanishak/service-station-etl/storage/dbs"
)

type Repair struct {
	AssignedEmployeeID int32  `json:"assignedEmployeeId"`
	ClientID           int32  `json:"clientId"`
	CarID              int32  `json:"carId"`
	Notes              string `json:"notes"`
	PaymentStatus      string `json:"paymentStatus"`
	RepairStatus       string `json:"repairStatus"`
	Description        string `json:"description"`
	RegistrationNumber string `json:"registrationNumber"`
	VIN                string `json:"vin"`
	Color              string `json:"color"`
}

type ServiceRepair struct {
	ServiceID int32 `json:"serviceId"`
	Quantity  int32 `json:"quantity"`
	RepairID  int32 `json:"repairId"`
}

// Define DTO structures
type RepairOrderDTO struct {
	RepairID         int32        `json:"repairId"`
	RepairNumber     string       `json:"repairNumber"`
	StartDate        time.Time    `json:"startDate"`
	Notes            string       `json:"notes"`
	TotalSum         float64      `json:"totalSum"`
	RepairStatus     string       `json:"repairStatus"`
	PaymentStatus    string       `json:"paymentStatus"`
	AssignedEmployee string       `json:"assignedEmployee"`
	VehicleInfo      VehicleInfo  `json:"vehicleInfo"`
	ClientName       string       `json:"clientName"`
	Services         []ServiceDTO `json:"services"`
	Address          string       `json:"address"`
}

type VehicleInfo struct {
	RegistrationNumber string `json:"registrationNumber"`
	Description        string `json:"description"`
	Model              string `json:"model"`
	Year               int32  `json:"year"`
}

type ServiceDTO struct {
	ServiceName    string         `json:"serviceName"`
	Quantity       int32          `json:"quantity"`
	CompletionDate time.Time      `json:"completionDate"`
	TotalPrice     float64        `json:"totalPrice"`
	InventoryUsed  []InventoryDTO `json:"inventoryUsed"`
}

type InventoryDTO struct {
	Amount          float64 `json:"amount"`
	ItemDescription string  `json:"itemDescription"`
	PricePerUnit    float64 `json:"pricePerUnit"`
	Quantity        int32   `json:"quantity"`
}

// AggregateServices merges service lines with the same name (a service
// added to a repair twice), keeping the order they first appear in.
func AggregateServices(services []ServiceDTO) []ServiceDTO {
	merged := make([]ServiceDTO, 0, len(services))
	index := make(map[string]int)
	for _, service := range services {
		if i, ok := index[service.ServiceName]; ok {
			merged[i].Quantity += service.Quantity
			merged[i].TotalPrice += service.TotalPrice
			merged[i].InventoryUsed = append(merged[i].InventoryUsed, service.InventoryUsed...)
			continue
		}
		index[service.ServiceName] = len(merged)
		service.InventoryUsed = append([]InventoryDTO(nil), service.InventoryUsed...)
		merged = append(merged, service)
	}
	return merged
}

// ConvertToDTO groups the ListOrders rows (one per repair x service line x
// part used) into repairs, in query order. Each service line is added once,
// with all of its parts; previously every part row added the line again,
// multiplying its quantity and price by the number of parts.
func ConvertToDTO(rows []dbs.ListOrdersRow) []RepairOrderDTO {
	var repairs []RepairOrderDTO
	repairIndex := make(map[int32]int)
	lineIndex := make(map[int32]int) // service_repairs_id -> index in its repair's Services

	for _, row := range rows {
		ri, ok := repairIndex[row.RepairID]
		if !ok {
			ri = len(repairs)
			repairIndex[row.RepairID] = ri
			repairs = append(repairs, RepairOrderDTO{
				Address:          row.Address.String,
				RepairID:         row.RepairID,
				RepairNumber:     row.RepairNumber.String,
				StartDate:        row.StartDate.Time,
				Notes:            row.Notes.String,
				TotalSum:         ParseFloat(row.TotalSum.String),
				RepairStatus:     row.RepairStatus.String,
				PaymentStatus:    row.PaymentStatus.String,
				AssignedEmployee: row.AssignedEmployee,
				VehicleInfo: VehicleInfo{
					RegistrationNumber: row.RegistrationNumber.String,
					Description:        row.VehicleDescription.String,
					Model:              row.VehicleModel.String,
					Year:               row.VehicleYear.Int32,
				},
				ClientName: row.ClientName,
				Services:   []ServiceDTO{},
			})
		}
		if !row.ServiceRepairsID.Valid { // repair without services yet
			continue
		}
		repair := &repairs[ri]

		li, ok := lineIndex[row.ServiceRepairsID.Int32]
		if !ok {
			li = len(repair.Services)
			lineIndex[row.ServiceRepairsID.Int32] = li
			repair.Services = append(repair.Services, ServiceDTO{
				ServiceName:    row.ServiceName.String,
				Quantity:       row.ServiceQuantity.Int32,
				CompletionDate: row.CompletionDate.Time,
				TotalPrice:     ParseFloat(row.ServiceTotalPrice.String),
				InventoryUsed:  []InventoryDTO{},
			})
		}
		if !row.ServiceInventoryID.Valid { // service without parts
			continue
		}
		repair.Services[li].InventoryUsed = append(repair.Services[li].InventoryUsed, InventoryDTO{
			Amount:          ParseFloat(row.InventoryAmount.String),
			ItemDescription: row.InventoryItemDescription.String,
			PricePerUnit:    ParseFloat(row.DetailItemPricePerUnit.String),
			Quantity:        row.InventoryItemQuantity.Int32,
		})
	}
	return repairs
}

// RepairNumber formats an order number from its opening date and id,
// e.g. RO-20230410-1201. (The sequence suffix used to be a constant 008.)
func RepairNumber(id int32, opened time.Time) string {
	return fmt.Sprintf("RO-%s-%03d", opened.Format("20060102"), id)
}

// BuildRepairs groups the flat ListOrders rows into repair orders with their
// services and parts, and computes each order's total.
func BuildRepairs(rows []dbs.ListOrdersRow) []RepairOrderDTO {
	repairs := ConvertToDTO(rows)
	for i := range repairs {
		repairs[i].Services = AggregateServices(repairs[i].Services)
		sum := 0.0
		for _, service := range repairs[i].Services {
			sum += service.TotalPrice
			for _, inventory := range service.InventoryUsed {
				sum += inventory.Amount * inventory.PricePerUnit
			}
		}
		repairs[i].TotalSum = sum
	}
	return repairs
}
