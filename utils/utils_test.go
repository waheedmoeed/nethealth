package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveBetweenChars(t *testing.T) {
	assert.Equal(t, "hello ", RemoveBetweenChars("hello <p>world <p> </p></p>", "<", ">"))
	assert.Equal(t, "hello ", RemoveBetweenChars("hello ", "<", ">"))
}
