package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlugin(t *testing.T) {
	_, report := New().Validate()
	for _, f := range report.Fields {
		assert.Empty(t, f.Errors)
	}
}
