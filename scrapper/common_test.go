package scrapper

import (
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
