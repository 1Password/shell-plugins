package schema

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsTitleCase(t *testing.T) {
	isTitleCase, _ := IsTitleCase("Test")
	assert.Equal(t, false, isTitleCase, fmt.Sprint("should return title case string"))
}
