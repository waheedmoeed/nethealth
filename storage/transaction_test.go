package storage

import (
	"testing"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/stretchr/testify/assert"
)

func TestStoreTransactionsToPDF(t *testing.T) {
	transactions := []*model.Transaction{
		{
			ServiceDate: "2023-10-01",
			ServiceCode: "SC001",
			Description: "Service Description 1",
			ClaimType:   "Claim Type 1",
			Units:       "10",
			Rate:        "15.00",
			Charge:      "150.00",
			Payer:       "Payer A",
			Batch:       "Batch001",
			Balance:     "150.00",
			Entity:      "Entity A",
		},
	}

	err := StoreTransactionsToPDF("test/test_transactions.pdf", transactions)
	assert.NoError(t, err)
}
