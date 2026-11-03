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

func AggregateServices(services []ServiceDTO) []ServiceDTO {
	aggregated := make(map[string]*ServiceDTO)
	for _, service := range services {
		if aggService, exists := aggregated[service.ServiceName]; exists {
			aggService.Quantity += service.Quantity
			aggService.TotalPrice += service.TotalPrice
			aggService.InventoryUsed = append(aggService.InventoryUsed, service.InventoryUsed...)
		} else {
			// Make a deep copy of the service to store in the map
			newService := service
			aggregated[service.ServiceName] = &newService
		}
	}

	var aggregatedServices []ServiceDTO
	for _, service := range aggregated {
		aggregatedServices = append(aggregatedServices, *service)
	}

	return aggregatedServices
}

// Convert data rows to DTO
func ConvertToDTO(rows []dbs.ListOrdersRow) []RepairOrderDTO {
	repairOrders := make(map[int32]*RepairOrderDTO)

	for _, row := range rows {
		if _, exists := repairOrders[row.RepairID]; !exists {
			repairOrders[row.RepairID] = &RepairOrderDTO{
				Address:          row.Address.String,
				RepairID:         row.RepairID,
				RepairNumber:     row.RepairNumber.String,
				StartDate:        row.StartDate.Time,
				Notes:            row.Notes.String,
				TotalSum:         ParseFloat(row.TotalSum.String),
				RepairStatus:     row.RepairStatus.String,
				PaymentStatus:    row.PaymentStatus.String,
				AssignedEmployee: row.AssignedEmployee.(string),
				VehicleInfo: VehicleInfo{
					RegistrationNumber: row.RegistrationNumber.String,
					Description:        row.VehicleDescription.String,
					Model:              row.VehicleModel.String,
					Year:               int32(row.VehicleYear.Int32),
				},
				ClientName: row.ClientName.(string),
				Services:   []ServiceDTO{},
			}
		}

		service := ServiceDTO{
			ServiceName:    row.ServiceName.String,
			Quantity:       int32(row.ServiceQuantity.Int32),
			CompletionDate: row.CompletionDate.Time,
			TotalPrice:     ParseFloat(row.ServiceTotalPrice.String),
			InventoryUsed:  []InventoryDTO{},
		}

		inventory := InventoryDTO{
			Amount:          ParseFloat(row.InventoryAmount.String),
			ItemDescription: row.InventoryItemDescription.String,
			PricePerUnit:    ParseFloat(row.DetailItemPricePerUnit.String),
			Quantity:        int32(row.InventoryItemQuantity.Int32),
		}

		service.InventoryUsed = append(service.InventoryUsed, inventory)
		repairOrders[row.RepairID].Services = append(repairOrders[row.RepairID].Services, service)
	}

	// Convert map to slice
	dtos := make([]RepairOrderDTO, 0, len(repairOrders))
	for _, dto := range repairOrders {
		dtos = append(dtos, *dto)
	}
	return dtos
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
