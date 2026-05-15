package main

import (
	"strings"
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

type FieldNameSplitFromCredNameSplit struct {
	FieldNameSplit []string
	CredNameSplit  []string
}

func TestValidateFieldNameSplitFromCredNameSplit(t *testing.T) {
	lengthCutoff := 7
	cases := []FieldNameSplitFromCredNameSplit{
		{
			CredNameSplit:  []string{"Personal", "Access", "Token"},
			FieldNameSplit: []string{"Token"},
		},
		{
			CredNameSplit:  []string{"Secret", "Key"},
			FieldNameSplit: []string{"Key"},
		},
		{
			CredNameSplit:  []string{"API", "Key"},
			FieldNameSplit: []string{"API", "Key"},
		},
		{
			CredNameSplit:  []string{"GitHub", "API", "Key"},
			FieldNameSplit: []string{"API", "Key"},
		},
		{
			CredNameSplit:  []string{"Credentials"},
			FieldNameSplit: []string{"Credentials"},
		},
	}

	for _, c := range cases {
		t.Run(strings.Join(c.CredNameSplit, " "), func(t *testing.T) {
			fieldNameSplit := fieldNameSplitFromCredNameSplit(c.CredNameSplit, lengthCutoff)
			assert.Equal(t, c.FieldNameSplit, fieldNameSplit)
		})
	}
}
