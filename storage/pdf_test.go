package storage

import (
	"os"
	"testing"

	"github.com/abdulwaheed/nethealth/model"
)

func TestStoreLaggerGroupsToPDF(t *testing.T) {
	laggerGroups := []*model.LaggerGroup{
		{
			ServiceDate: "2023-10-10",
			Laggers: []*model.Lagger{
				{TxDate: "2023-10-09", Type: "Ty", ControlNumber: "12345, 6456456, 4564564, 564564, 564564 56456, 5645645, 645645645, 6456456", Description: "Description1 4564564 56456456 45645645, 6456456456, 456456456", Seq: "001", ServiceDate: "2023-10-10", Category: "Cat1", DBAmount: "100.00", CRAmount: "50.00", Balance: "50.00", PDFLink: "https://google.com"},
				{TxDate: "2023-10-09", Type: "Ty", ControlNumber: "12345, 6456456, 4564564, 564564, 564564 56456, 5645645, 645645645, 6456456", Description: "Description1 4564564 56456456 45645645, 6456456456, 456456456", Seq: "001", ServiceDate: "2023-10-10", Category: "Cat1", DBAmount: "100.00", CRAmount: "50.00", Balance: "50.00", PDFLink: ""},
			},
			EstimatedAdjustments: []*model.Lagger{
				{TxDate: "2023-10-09", Type: "Ty", ControlNumber: "4564564, 56456456, 456456", Description: "Description2", Seq: "002", ServiceDate: "2023-10-10", Category: "Cat2", DBAmount: "200.00", CRAmount: "100.00", Balance: "100.00"},
			},
		},
		{
			ServiceDate: "2024-10-10",
			Laggers: []*model.Lagger{
				{TxDate: "2023-10-09", Type: "Ty", ControlNumber: "12345, 6456456, 4564564, 564564, 564564 56456, 5645645, 645645645, 6456456", Description: "Description1 4564564 56456456 45645645, 6456456456, 456456456", Seq: "001", ServiceDate: "2023-10-10", Category: "Cat1", DBAmount: "100.00", CRAmount: "50.00", Balance: "50.00", PDFLink: "https://google.com"},
				{TxDate: "2023-10-09", Type: "Ty", ControlNumber: "12345, 6456456, 4564564, 564564, 564564 56456, 5645645, 645645645, 6456456", Description: "Description1 4564564 56456456 45645645, 6456456456, 456456456", Seq: "001", ServiceDate: "2023-10-10", Category: "Cat1", DBAmount: "100.00", CRAmount: "50.00", Balance: "50.00", PDFLink: "https://google.com"},
			},
			EstimatedAdjustments: []*model.Lagger{
				{TxDate: "2023-10-09", Type: "Ty", ControlNumber: "4564564, 56456456, 456456", Description: "Description2", Seq: "002", ServiceDate: "2023-10-10", Category: "Cat2", DBAmount: "200.00", CRAmount: "100.00", Balance: "100.00"},
			},
		},
	}

	fileName := "test_laggers.pdf"
	err := StoreLaggerGroupsToPDF(fileName, laggerGroups)
	if err != nil {
		t.Fatalf("Failed to store lagger groups to PDF: %v", err)
	}

	// Check if the PDF file was created
	_, err = os.Stat(fileName + ".pdf")
	if os.IsNotExist(err) {
		t.Fatalf("PDF file was not created")
	}
	// Clean up the created PDF file after test
	defer os.Remove(fileName + ".pdf")
}