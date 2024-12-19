package scrapper

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSterializeAccountNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"Account: 12345", 12345},
		{"Account: 67890", 67890},
		{"Invalid: abcde", 0}, // Test case for invalid input
		{"Account: ", 0},      // Test case for missing number
	}

	for _, test := range tests {
		result := sterializeAccountNumber(test.input)
		assert.Equal(t, test.expected, result, "they should be equal")
	}
}

func TestSample(t *testing.T) {
	claimsUrl := "https://p13006.therapy.nethealth.com/Financials#patient/details/351/transactions"

	claimsUrl = fmt.Sprintf("%s/claims", claimsUrl[:strings.LastIndex(claimsUrl, "/")])
	assert.Equal(t, "https://p13006.therapy.nethealth.com/Financials#patient/details/351/claims", claimsUrl)

	claimsUrl = "https://p13006.therapy.nethealth.com/Financials#patient/details/351/transactions"

	claimsUrl = fmt.Sprintf("%s/agingSummary", claimsUrl[:strings.LastIndex(claimsUrl, "/")])
	assert.Equal(t, "https://p13006.therapy.nethealth.com/Financials#patient/details/351/agingSummary", claimsUrl)
}
