package dto

import (
	"testing"
	"time"
)

func TestDocumentNumbers(t *testing.T) {
	opened := time.Date(2023, 4, 10, 15, 0, 0, 0, time.UTC)
	for got, want := range map[string]string{
		RepairNumber(1201, opened): "RO-20230410-1201",
		RepairNumber(7, opened):    "RO-20230410-007",
		SupplyNumber(1001):         "SN-001001",
		InvoiceNumber(1001, 2121):  "INV-1001-2121",
	} {
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
}
