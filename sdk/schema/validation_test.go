package schema

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsTitleCaseWord(t *testing.T) {
	cases := map[string]struct {
		word     string
		expected bool
	}{
		"when title case string provided": {
			word:     "Title",
			expected: true,
		},
		"when empty string provided": {
			word:     "",
			expected: false,
		},
		"when starts with special character": {
			word:     "_",
			expected: false,
		},
		"when string with more than one str provided": {
			word:     "Access Token",
			expected: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			isTitleCase := IsTitleCaseWord(tc.word)
			assert.Equal(t, tc.expected, isTitleCase, fmt.Sprintf("should return %t", tc.expected))
		})
	}
}

func TestIsTitleCaseString(t *testing.T) {
	cases := map[string]struct {
		str      string
		expected bool
	}{
		"when one str string provided": {
			str:      "Title",
			expected: true,
		},
		"when string has more than one str": {
			str:      "Title Case String",
			expected: true,
		},
		"when empty string provided": {
			str:      "",
			expected: false,
		},
		"when first word in lowercase": {
			str:      "title Case String",
			expected: false,
		},
		"when middle word in lowercase": {
			str:      "Title case String",
			expected: false,
		},
		"when last word in lowercase": {
			str:      "Title Case string",
			expected: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			isTitleCase := IsTitleCaseString(tc.str)
			assert.Equal(t, tc.expected, isTitleCase, fmt.Sprintf("should return %t", tc.expected))
		})
	}
}

func TestContainsLowercaseLettersOrDigits(t *testing.T) {
	cases := map[string]struct {
		str      string
		expected bool
	}{
		"when there are only lowercase letters": {
			str:      "aws",
			expected: true,
		},
		"when there are only digits": {
			str:      "911",
			expected: true,
		},
		"when there are lowercase letters and digits": {
			str:      "aws911",
			expected: true,
		},
		"when empty string provided": {
			str:      "",
			expected: false,
		},
		"when there is Capital letter": {
			str:      "Aws",
			expected: false,
		},
		"when string has more than one word": {
			str:      "aws911 test",
			expected: false,
		},
		"when string contains special char": {
			str:      "aws911_.test",
			expected: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			isValid := ContainsLowercaseLettersOrDigits(tc.str)
			assert.Equal(t, tc.expected, isValid, fmt.Sprintf("should return %t", tc.expected))
		})
	}
}

func TestPluginValidateHasHeading(t *testing.T) {
	expectedHeading := "Plugin: test"
	p := Plugin{Name: "test"}
	_, report := p.Validate()

	assert.Equal(t, expectedHeading, report.Heading, fmt.Sprintf("plugin should have heading %s", expectedHeading))
}

func TestPluginValidateEachReportFieldHasError(t *testing.T) {
	p := Plugin{}
	_, report := p.Validate()
	c := report.Checks[0]

	assert.False(t, c.Assertion, fmt.Sprintf("\"%s\" validation is erroneous", c.Description))
}
