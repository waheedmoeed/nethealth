package scrapper

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareDataRoomDir(t *testing.T) {
	tempDir := "../data/AgeilityatBearCreek/abnerclaire_8108" // creates a temporary directory for testing

	err := prepareDataRoomDir(tempDir)
	assert.NoError(t, err)

	directories := []string{"transactions", "ledger", "claims", "agingsummary", "benefits", "transactionbreakdowns"}
	for _, dir := range directories {
		path := filepath.Join(tempDir, dir)
		info, err := os.Stat(path)
		assert.NoError(t, err)
		assert.True(t, info.IsDir(), "Expected %s to be a directory", path)
	}
}
