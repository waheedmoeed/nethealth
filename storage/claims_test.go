package storage

import (
	"testing"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/stretchr/testify/assert"
)

func TestStoreClaimsToPDF(t *testing.T) {
	claims := []*model.Claim{
		{
			CreationDate:    "2020-12-01",
			ServicesFrom:    "2020-12-01",
			ServicesThrough: "2020-12-31",
			ClaimNumber:     "123456",
			ClaimType:       "medical",
			BatchNumber:     "111111",
			Entity:          "Abdul",
			PayingAgency:    "BlueCross   gkugg kuygu",
			PayerPlan:       "PPO",
			PayerSequence:   "2",
			ClaimAmount:     "$100.00",
			PDFLink:         "https://example.com/claim123456.pdf",
			UploadedLink:    "https://example.com/claim123456.pdf",
		},
	}
	err := StoreClaimsToPDF("test/claims.pdf", claims)
	assert.NoError(t, err)
}
