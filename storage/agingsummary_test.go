package storage

import (
	"testing"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/stretchr/testify/assert"
)

func TestStoreAgingSummaryToPDF(t *testing.T) {
	agingSummary := []*model.AgingSummary{
		{
			PayerPlan:      "Payer Plan",
			Credits:        "Credits",
			Total:          "Total",
			Current:        "Current",
			ThirtyDays:     "30Days",
			SixtyDays:      "60Days",
			NinetyDays:     "90Days",
			OneTwentyDays:  "120Days",
			OneEightyDays:  "180Days",
			MoreThanEigthy: "Greater 180Days",
		},
	}
	err := StoreAgingSummaryToPDF("test/agingSummary.pdf", agingSummary)
	assert.NoError(t, err)
}
