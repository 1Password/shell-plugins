package schema

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsTitleCaseWord_ReturnsTrue(t *testing.T) {
	isTitleCase := IsTitleCaseWord("Title")
	assert.Equal(t, true, isTitleCase, fmt.Sprint("should return true"))
}

func TestIsTitleCaseWord_ReturnsFalse(t *testing.T) {
	cases := map[string]struct {
		word string
	}{
		"when empty string provided": {
			word: "",
		},
		"when starts with special character": {
			word: "_",
		},
		"when string with more than one str provided": {
			word: "Access Token",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			isTitleCase := IsTitleCaseWord(tc.word)
			assert.Equal(t, false, isTitleCase, fmt.Sprint("should return false"))
		})
	}
}

func TestIsTitleCaseString_ReturnsTrue(t *testing.T) {
	cases := map[string]struct {
		str string
	}{
		"when one str string provided": {
			str: "Title",
		},
		"when string has more than one str": {
			str: "Title Case String",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			isTitleCase := IsTitleCaseString(tc.str)
			assert.Equal(t, true, isTitleCase, fmt.Sprint("should return true"))
		})
	}
}

func TestIsTitleCaseString_ReturnsFalse(t *testing.T) {
	cases := map[string]struct {
		str string
	}{
		"when empty string provided": {
			str: "",
		},
		"when first word in lowercase": {
			str: "title Case String",
		},
		"when middle word in lowercase": {
			str: "Title case String",
		},
		"when last word in lowercase": {
			str: "Title Case string",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			isTitleCase := IsTitleCaseString(tc.str)
			assert.Equal(t, false, isTitleCase, fmt.Sprint("should return false"))
		})
	}
}

func TestContainsLowercaseLettersOrDigits_ReturnsTrue(t *testing.T) {
	cases := map[string]struct {
		str string
	}{
		"when there are only lowercase letters": {
			str: "aws",
		},
		"when there are only digits": {
			str: "911",
		},
		"when there are lowercase letters and digits": {
			str: "aws911",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			isValid := ContainsLowercaseLettersOrDigits(tc.str)
			assert.Equal(t, true, isValid, fmt.Sprint("should return true"))
		})
	}
}

func TestContainsLowercaseLettersOrDigits_ReturnsFalse(t *testing.T) {
	cases := map[string]struct {
		str string
	}{
		"when empty string provided": {
			str: "",
		},
		"when there is Capital letter": {
			str: "Aws",
		},
		"when string has more than one word": {
			str: "aws911 test",
		},
		"when string contains special char": {
			str: "aws911_.test",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			isValid := ContainsLowercaseLettersOrDigits(tc.str)
			assert.Equal(t, false, isValid, fmt.Sprint("should return true"))
		})
	}
}
