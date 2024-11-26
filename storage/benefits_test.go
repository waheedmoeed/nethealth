package storage

import (
	"testing"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/stretchr/testify/assert"
)

func TestStoreBenefitsToPDF(t *testing.T) {
	benefits := []*model.Benefit{
		{Text: `Case: 915-0012769 (03/27/2024 - (open) - Ageility at Bear Creek)
				Payer Order/Type: 1: (Medicare Part B)
				Paying Agency: Medicare
				Payer Plan: Med B Novitas NJ JL
				Entitlement Date: 03/27/2024
				Policy No: 4QG1QJ6UQ26
				Subscriber: SELF
				Subscriber's DOB: 7/1/1937
				Benefits: (03/27/2024 - (open))
				Network Usage: In Network
				Applies to Discipline(s): PT, OT, ST

				Payer Order/Type: 2: (Commercial Insurance)
				Paying Agency: Commercial
				Payer Plan: AARP Coinsurance
				Entitlement Date: 03/27/2024
				Policy No: 07993390911
				Subscriber: SELF
				Subscriber's DOB: 7/1/1937
				Benefits: (03/27/2024 - (open))
				Network Usage: Out Of Network
				Applies to Discipline(s): PT, OT, ST

				Payer Order/Type: 3: (FRP)
				Paying Agency:
				FRP: Abner, Claire`},
		{Text: "Benefit 2"},
	}

	err := StoreBenefitsToPDF("test/test_benefits.pdf", benefits)
	assert.NoError(t, err)

	// Add additional checks here to verify the PDF content if needed
}
