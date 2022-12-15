package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCredentialName(t *testing.T) {
	err := validateCredentialName(any("Custom Credential Token"))
	assert.Nil(t, err)
}

func TestValidateCredentialNameReturnError(t *testing.T) {
	cases := map[string]string{
		"when no value provided":           "",
		"when first word not capitalized":  "custom Credential Token",
		"when middle word not capitalized": "Custom credential Token",
		"when last word not capitalized":   "Custom Credential token",
		"when lowercase string provided":   "custom credential token",
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := validateCredentialName(any(tc))
			assert.NotNil(t, err)
		})
	}
}
