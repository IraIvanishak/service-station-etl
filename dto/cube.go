package dto

import (
	"database/sql"
	"fmt"

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
