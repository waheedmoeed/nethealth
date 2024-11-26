package storage

import (
	"testing"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/stretchr/testify/assert"
)

// TestStoreTransactionsBreakdownToPDF verifies that the StoreTransactionsBreakdownToPDF
// function successfully creates a PDF from a slice of TransactionBreakdown objects
// without returning an error.
func TestStoreTransactionsBreakdownToPDF(t *testing.T) {
	transactionDetils := []*model.TransactionDetail{
		{
			ServiceDate: "2020-12-01",
			ServiceCode: "12345",
			Description: "Service Description",
			Charge:      "$100.00",
			Balance:     "$200.00",
			TransactionBreakdown: []*model.TransactionBreakdown{
				{
					Date:         "2020-12-01",
					ResonCode:    "12345",
					Description:  "Breakdown Description",
					Amount:       "$50.00",
					Reference:    "123456789",
					Payer:        "Payer A",
					Batch:        "12345",
					PDFLink:      "https://example.com/pdf1.pdf",
					UploadedLink: "https://example.com/pdf1.pdf",
				},
			},
		},
		{
			ServiceDate: "2020-12-01",
			ServiceCode: "12345",
			Description: "Service Description",
			Charge:      "$100.00",
			Balance:     "$200.00",
			TransactionBreakdown: []*model.TransactionBreakdown{
				{
					Date:         "2020-12-01",
					ResonCode:    "12345",
					Description:  "Breakdown Description",
					Amount:       "$50.00",
					Reference:    "123456789",
					Payer:        "Payer A",
					Batch:        "12345",
					PDFLink:      "https://example.com/pdf2.pdf",
					UploadedLink: "https://example.com/pdf2.pdf",
				},
			},
		},
		{
			ServiceDate: "2020-12-01",
			ServiceCode: "12345",
			Description: "Service Description",
			Charge:      "$100.00",
			Balance:     "$200.00",
			TransactionBreakdown: []*model.TransactionBreakdown{
				{
					Date:         "2020-12-01",
					ResonCode:    "12345",
					Description:  "Breakdown Description",
					Amount:       "$50.00",
					Reference:    "123456789",
					Payer:        "Payer A",
					Batch:        "12345",
					PDFLink:      "https://example.com/pdf3.pdf",
					UploadedLink: "https://example.com/pdf3.pdf",
				},
			},
		},	
	}
	err := StoreTransactionDetailsToPDF("test/test_transactions_breakdown.pdf", transactionDetils)
	assert.NoError(t, err)
}
