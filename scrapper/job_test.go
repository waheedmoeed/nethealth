package scrapper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadAndSavePDF(t *testing.T) {
	// Create a test server to mock the PDF download

	// Set up test file path and auth key
	filePath := "test/testfile"
	authKey := "2EB4D0A5831369F8CEC00E35CC9710F6C7A028962688963F7457774CC09AA1A0A189382775CD0E515D0AA654964450330C4EDAD6EF5E40273767D4F7B63DF64A467E5A95A76C8DB48016A84027FD7EA280D3F925691B35B1FDE20A44F288F99486B5E8BA87B06F25F318A289FFAD675E8D91A9B9A11AA179D1C694CB0AE91ABD9E47D0E6172C100E27D575B2B0F73D54FA37BD578766C94E245A89E53974710F6970EB6B6A1C2FF3F76E57D2F8CA10BF66B4077F1A5797B23C68B1AC4C49BB43"

	// Call the function
	err := downloadAndSavePDF("https://p13006.therapy.nethealth.com/Financials.Service/api/deposits/167575/pdf?customerId=81808", filePath, authKey)

	// Assert no error
	assert.NoError(t, err)

	// Check if the file has been created
	_, err = os.Stat(filePath + ".pdf")
	assert.NoError(t, err)

	// Clean up the file
	defer os.Remove(filePath + ".pdf")
}
