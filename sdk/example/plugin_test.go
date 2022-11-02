package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlugin(t *testing.T) {
	_, errs := New().Validate()
	assert.Empty(t, errs)
}
