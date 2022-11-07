package plugintest

import (
	"fmt"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
)

func TestSecretContainsSuffix(t *testing.T) {
	v := &schema.ValueComposition{
		Length: 10,
		Charset: schema.Charset{
			Uppercase: true,
		},
	}
	result := ExampleSecretFromComposition(v)
	hasExampleSuffix := strings.HasSuffix(result, secretExampleSuffix)

	assert.Equal(t, true, hasExampleSuffix, fmt.Sprintf("should contain %s suffix", secretExampleSuffix))
}

func TestSecretContainsLowercaseSuffix(t *testing.T) {
	v := &schema.ValueComposition{
		Length: 10,
		Charset: schema.Charset{
			Lowercase: true,
		},
	}
	result := ExampleSecretFromComposition(v)
	hasExampleSuffix := strings.HasSuffix(result, strings.ToLower(secretExampleSuffix))

	assert.Equal(t, true, hasExampleSuffix, fmt.Sprintf("should contain lowercase %s suffix", secretExampleSuffix))
}

func TestSecretHasNoSuffix(t *testing.T) {
	cases := map[string]struct {
		secretLength int
		charset      schema.Charset
	}{
		"when secret length equal suffix length": {
			secretLength: len(secretExampleSuffix),
			charset:      schema.Charset{Uppercase: true},
		},
		"when secret length less than suffix length": {
			secretLength: len(secretExampleSuffix) - 1,
			charset:      schema.Charset{Uppercase: true},
		},
		"when no Uppercase or Lowercase charset": {
			secretLength: len(secretExampleSuffix) + 1,
			charset:      schema.Charset{Symbols: true},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			v := &schema.ValueComposition{
				Length:  tc.secretLength,
				Charset: tc.charset,
			}
			result := ExampleSecretFromComposition(v)
			hasNoSuffix := !strings.HasSuffix(result, secretExampleSuffix)

			assert.Equal(t, true, hasNoSuffix, fmt.Sprintln("should not has suffix"))
		})
	}
}

func TestSecretContainsDigitsOnly(t *testing.T) {
	v := &schema.ValueComposition{
		Length: 10,
		Charset: schema.Charset{
			Digits: true,
		},
	}
	result := ExampleSecretFromComposition(v)
	hasDigitsOnly, _ := containsOnlyDigits(result)

	assert.Equal(t, true, hasDigitsOnly, fmt.Sprintln("should contain digits only"))
}

func TestSecretContainsPrefix(t *testing.T) {
	prefix := "ABC"
	v := &schema.ValueComposition{
		Length: 10,
		Prefix: prefix,
		Charset: schema.Charset{
			Uppercase: true,
		},
	}
	result := ExampleSecretFromComposition(v)
	hasExampleSuffix := strings.HasPrefix(result, prefix)

	assert.Equal(t, true, hasExampleSuffix, fmt.Sprintf("should contain %s prefix", prefix))
}

func TestSecretWithExpectedLength(t *testing.T) {
	expectedLength := 20
	v := &schema.ValueComposition{
		Length: expectedLength,
		Charset: schema.Charset{
			Uppercase: true,
		},
	}
	result := ExampleSecretFromComposition(v)

	assert.Equal(t, expectedLength, len(result), fmt.Sprintf("should have %d chars length", expectedLength))
}

func TestStingFromCharsetReturnErrorWhenNoCharsetProvided(t *testing.T) {
	_, err := stringFromCharset(10, "")
	if err == nil {
		t.FailNow()
	}
}

func TestStringFromCharsetContainsOnly(t *testing.T) {
	cases := map[string]struct {
		charset      schema.Charset
		containsFunc func(str string) (bool, error)
	}{
		"capital letters": {
			charset:      schema.Charset{Uppercase: true},
			containsFunc: containsOnlyCapitalLetters,
		},
		"lowercase letters": {
			charset:      schema.Charset{Lowercase: true},
			containsFunc: containsOnlyLowercaseLetters,
		},
		"digits": {
			charset:      schema.Charset{Digits: true},
			containsFunc: containsOnlyDigits,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			stringLength := 20
			charset := charsToUse(tc.charset)
			result, _ := stringFromCharset(stringLength, charset)
			hasOnly, err := tc.containsFunc(result)
			if err != nil {
				t.Log(err)
				t.FailNow()
			}

			assert.Equal(t, true, hasOnly)
		})
	}
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
