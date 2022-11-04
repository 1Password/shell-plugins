package strgenerator

import (
	"fmt"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
)

func TestExampleSecretGenerationFails(t *testing.T) {
	v := &schema.ValueComposition{
		Length:  10,
		Charset: schema.Charset{},
	}
	_, err := ExampleSecretFromComposition(v)
	if err == nil {
		t.Fail()
	}
}

func TestExampleSecretGenerationContainsSuffix(t *testing.T) {
	v := &schema.ValueComposition{
		Length: 10,
		Charset: schema.Charset{
			Uppercase: true,
		},
	}
	result, _ := ExampleSecretFromComposition(v)
	hasExampleSuffix := strings.HasSuffix(result, secretExampleSuffix)

	assert.Equal(t, true, hasExampleSuffix, fmt.Sprintf("should contain %s suffix", secretExampleSuffix))
}

func TestExampleSecretGenerationContainsPrefix(t *testing.T) {
	prefix := "ABC"
	v := &schema.ValueComposition{
		Length: 10,
		Prefix: prefix,
		Charset: schema.Charset{
			Uppercase: true,
		},
	}
	result, _ := ExampleSecretFromComposition(v)
	hasExampleSuffix := strings.HasPrefix(result, prefix)

	assert.Equal(t, true, hasExampleSuffix, fmt.Sprintf("should contain %s prefix", prefix))
}

func TestExampleSecretGenerationWithExpectedLength(t *testing.T) {
	expectedLength := 20
	v := &schema.ValueComposition{
		Length: expectedLength,
		Charset: schema.Charset{
			Uppercase: true,
		},
	}
	result, _ := ExampleSecretFromComposition(v)

	assert.Equal(t, expectedLength, len(result), fmt.Sprintf("should have %d chars length", expectedLength))
}

func TestStingFromCharsetReturnErrorWhenNoCharsetProvided(t *testing.T) {
	_, err := stringFromCharset(10, "")
	if err == nil {
		t.FailNow()
	}
}

func TestStingFromCharsetHasOnlyCapitalLetters(t *testing.T) {
	stringLength := 20
	c := schema.Charset{
		Uppercase: true,
	}
	charset := charsToUse(c)
	result, _ := stringFromCharset(stringLength, charset)
	hasOnlyCapital, err := containsOnlyCapitalLetters(result)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	assert.Equal(t, true, hasOnlyCapital, "should contain capital letters only")
}

func TestStringFromCharsetHasOnlyLowercaseLetters(t *testing.T) {
	stringLength := 20
	c := schema.Charset{
		Lowercase: true,
	}
	charset := charsToUse(c)
	result, _ := stringFromCharset(stringLength, charset)
	hasOnlyLowercase, err := containsOnlyLowercaseLetters(result)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	assert.Equal(t, true, hasOnlyLowercase, "should contain lowercase letters only")
}

func TestStringFromCharsetHasOnlyDigits(t *testing.T) {
	stringLength := 20
	c := schema.Charset{
		Digits: true,
	}
	charset := charsToUse(c)
	result, _ := stringFromCharset(stringLength, charset)
	hasOnlyDigits, err := containsOnlyDigits(result)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	assert.Equal(t, true, hasOnlyDigits, "should contain digits only")
}

func containsOnlyCapitalLetters(str string) (bool, error) {
	return regexp.Match("^[A-Z]+$", []byte(str))
}

func containsOnlyLowercaseLetters(str string) (bool, error) {
	return regexp.Match("^[a-z]+$", []byte(str))
}

func containsOnlyDigits(str string) (bool, error) {
	return regexp.Match("^[0-9]+$", []byte(str))
}
