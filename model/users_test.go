package model

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadUsersFromCSVFile(t *testing.T) {
	users, err := ReadUsersFromCSVFile(context.Background(), "/Users/abdulwaheed/Workspace/Bionomical/NetHealth/userscvs/Ageility at 75 State Street_users.csv", "Ageility at 75 State Street")
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}
