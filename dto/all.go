package dto

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/IraIvanishak/service-station-etl/storage/dbs"
)

// SupplyDetailsDTO - структура для user-friendly формату.
type SupplyDetailsDTO struct {
	DetailID        int32
	Quantity        string
	TotalCost       string
	Inventory       string
	Article         string
	Price           string
	Dimensions      string
	Weight          string
	SupplierName    string
	SupplierContact string
	DateRange       string
}

// ToDTO конвертує GetCubeSupplyRow в SupplyDetailsDTO.
func ToDTOSupply(row dbs.GetCubeSupplyRow) SupplyDetailsDTO {
	return SupplyDetailsDTO{
		DetailID:        row.DetailID,
		Quantity:        nullInt32ToString(row.DetailQuantity),
		TotalCost:       nullStringToString(row.DetailTotalCost),
		Inventory:       nullStringToString(row.InventoryDescription),
		Article:         nullStringToString(row.Article),
		Price:           nullStringToString(row.Price),
		Dimensions:      nullStringToString(row.Dimensions),
		Weight:          nullStringToString(row.Weight),
		SupplierName:    nullStringToString(row.SupplierName),
		SupplierContact: formatSupplierContact(row.SupplierPhoneNumber, row.SupplierAddress),
		DateRange:       formatDateRange(row.StartYear, row.StartMonth, row.StartDay, row.EndYear, row.EndMonth, row.EndDay),
	}
}

// Helper functions to handle sql.NullString and sql.NullInt32 types.
func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return "N/A"
}

func nullInt32ToString(ni sql.NullInt32) string {
	if ni.Valid {
		return fmt.Sprintf("%d", ni.Int32)
	}
	return "N/A"
}

func formatSupplierContact(phone sql.NullString, address sql.NullString) string {
	return fmt.Sprintf("Phone: %s, Address: %s", nullStringToString(phone), nullStringToString(address))
}

func formatDateRange(startYear, startMonth, startDay, endYear, endMonth, endDay sql.NullInt32) string {
	start := fmt.Sprintf("%s-%s-%s", nullInt32ToString(startYear), nullInt32ToString(startMonth), nullInt32ToString(startDay))
	end := fmt.Sprintf("%s-%s-%s", nullInt32ToString(endYear), nullInt32ToString(endMonth), nullInt32ToString(endDay))
	return fmt.Sprintf("%s to %s", start, end)
}

// ----------------------------

type ServiceInventoryDTO struct {
	ServiceInventoryFactID   int32
	PartsCost                string
	TotalCost                string
	PartsCostPercentage      string
	EstimatedServiceDuration string
	ActualServiceDuration    string
	ServiceDurationDeviation string
	InventoryAmount          string
	ServiceRepairDimID       int32
	ServiceID                int32
	RepairID                 int32
	ServiceQuantity          int
	ServiceTotalPrice        string
	ServiceName              string
	ServiceDescription       string
	ServicePrice             string
	ServiceEstimatedDuration string
	RepairNotes              string
	RepairNumber             string
	InventoryDimID           int32
	InventoryBkDetailID      int32
	InventoryDescription     string
	InventoryArticle         string
	InventoryPrice           string
	InventoryDimensions      string
	InventoryWeight          string
}

func ConvertToUserFriendly(row dbs.GetCubeServiceInventoryRow) ServiceInventoryDTO {
	return ServiceInventoryDTO{
		ServiceInventoryFactID:   row.ServiceInventoryFactID,
		PartsCost:                row.PartsCost.String,
		TotalCost:                row.TotalCost.String,
		PartsCostPercentage:      row.PartsCostPercentage.String,
		EstimatedServiceDuration: row.EstimatedServiceDuration.String,
		ActualServiceDuration:    row.ActualServiceDuration.String,
		ServiceDurationDeviation: row.ServiceDurationDeviation.String,
		InventoryAmount:          row.InventoryAmount.String,
		ServiceRepairDimID:       row.ServiceRepairDimID,
		ServiceID:                row.ServiceID,
		RepairID:                 row.RepairID,
		ServiceQuantity:          int(row.ServiceQuantity.Int32),
		ServiceTotalPrice:        row.ServiceTotalPrice.String,
		ServiceName:              row.ServiceName.String,
		ServiceDescription:       row.ServiceDescription.String,
		ServicePrice:             row.ServicePrice.String,
		ServiceEstimatedDuration: row.ServiceEstimatedDuration.String,
		RepairNotes:              row.RepairNotes.String,
		RepairNumber:             row.RepairNumber.String,
		InventoryDimID:           row.InventoryDimID,
		InventoryBkDetailID:      row.InventoryBkDetailID,
		InventoryDescription:     row.InventoryDescription.String,
		InventoryArticle:         row.InventoryArticle.String,
		InventoryPrice:           row.InventoryPrice.String,
		InventoryDimensions:      row.InventoryDimensions.String,
		InventoryWeight:          row.InventoryWeight.String,
	}
}

// func to gebereate smt like INV-636217
func CreateInvoiceNumber() string {
	t := time.Now()
	return fmt.Sprintf("INV-%d", t.Unix())
}

// func to generate smt like SN-000900 baced on id
func CreateSupplyNumber(id int) string {
	return fmt.Sprintf("SN-000%d", id)
}

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
type ServiceRepair struct {
	ServiceID int32 `json:"serviceId"`
	Quantity  int32 `json:"quantity"`
	RepairID  int32 `json:"repairId"`
}

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

// func that generates smt in format like RO-20230410-008
func CreateRepairNumber() string {
	t := time.Now()
	return fmt.Sprintf("RO-%d%02d%02d-%03d", t.Year(), t.Month(), t.Day(), 8)
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

func ParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return f
}

type Station struct {
	Id      int32  `json:"Id"`
	Address string `json:"Address"`
	City    int32  `json:"CityId"`
}
