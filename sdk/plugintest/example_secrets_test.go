package plugintest

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/stretchr/testify/assert"
)

func TestSecretContainsSuffix(t *testing.T) {
	v := schema.ValueComposition{
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
	v := schema.ValueComposition{
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
			v := schema.ValueComposition{
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
	v := schema.ValueComposition{
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
	v := schema.ValueComposition{
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
	v := schema.ValueComposition{
		Length: expectedLength,
		Charset: schema.Charset{
			Uppercase: true,
		},
	}
	result := ExampleSecretFromComposition(v)

	assert.Equal(t, expectedLength, len(result), fmt.Sprintf("should have %d chars length", expectedLength))
}

func TestSecretWithFixedLengthAndPrefix(t *testing.T) {
	v := schema.ValueComposition{
		Length: 40,
		Prefix: "ghp_",
		Charset: schema.Charset{
			Uppercase: true,
			Lowercase: true,
			Digits:    true,
		},
	}
	result := ExampleSecretFromComposition(v)

	assert.Equal(t, 40, len(result), "should have exactly the declared length")
	assert.True(t, strings.HasPrefix(result, "ghp_"), "should contain ghp_ prefix")
	assert.True(t, strings.HasSuffix(result, secretExampleSuffix), fmt.Sprintf("should contain %s suffix", secretExampleSuffix))
	assertOnlyCharsFrom(t, strings.TrimPrefix(result, "ghp_"), v.Charset)
}

func TestSecretWithoutLengthIsComposedAdditively(t *testing.T) {
	cases := map[string]struct {
		prefix         string
		charset        schema.Charset
		expectedSuffix string
	}{
		"no prefix, letters": {
			charset:        schema.Charset{Uppercase: true, Digits: true},
			expectedSuffix: secretExampleSuffix,
		},
		"letters, digits, @ and .": {
			charset:        schema.Charset{Lowercase: true, Uppercase: true, Digits: true, Specific: []rune{'@', '.'}},
			expectedSuffix: secretExampleSuffix,
		},
		"letters, digits and symbols": {
			charset:        schema.Charset{Lowercase: true, Uppercase: true, Digits: true, Symbols: true},
			expectedSuffix: secretExampleSuffix,
		},
		"pypi- prefix": {
			prefix:         "pypi-",
			charset:        schema.Charset{Uppercase: true, Lowercase: true, Digits: true, Specific: []rune{'-', '_'}},
			expectedSuffix: secretExampleSuffix,
		},
		"lowercase letters": {
			prefix:         "abc_",
			charset:        schema.Charset{Lowercase: true},
			expectedSuffix: strings.ToLower(secretExampleSuffix),
		},
		"no letters": {
			prefix:  "12-",
			charset: schema.Charset{Digits: true},
		},
		"prefix longer than the default body": {
			prefix:         strings.Repeat("p", defaultBodyLength+10),
			charset:        schema.Charset{Uppercase: true, Lowercase: true},
			expectedSuffix: secretExampleSuffix,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			v := schema.ValueComposition{Prefix: tc.prefix, Charset: tc.charset}

			var result string
			assert.NotPanics(t, func() { result = ExampleSecretFromComposition(v) })

			assert.Equal(t, len(tc.prefix)+defaultBodyLength+len(tc.expectedSuffix), len(result))
			assert.True(t, strings.HasPrefix(result, tc.prefix))
			assert.True(t, strings.HasSuffix(result, tc.expectedSuffix))
			body := strings.TrimSuffix(strings.TrimPrefix(result, tc.prefix), tc.expectedSuffix)
			assert.Len(t, body, defaultBodyLength)
			assertOnlyCharsFrom(t, body, tc.charset)
		})
	}
}

func TestSecretWithShortLengthAndPrefix(t *testing.T) {
	cases := map[string]struct {
		length         int
		prefix         string
		expectedLength int
		hasSuffix      bool
	}{
		"suffix is dropped to keep the declared length": {
			length:         8,
			prefix:         "pypi-",
			expectedLength: 8,
		},
		"suffix is dropped when it would leave no random characters": {
			length:         12,
			prefix:         "pypi-",
			expectedLength: 12,
		},
		"suffix is kept once a random character fits": {
			length:         13,
			prefix:         "pypi-",
			expectedLength: 13,
			hasSuffix:      true,
		},
		"prefix is kept whole when longer than length": {
			length:         3,
			prefix:         "pypi-",
			expectedLength: 5,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			v := schema.ValueComposition{
				Length:  tc.length,
				Prefix:  tc.prefix,
				Charset: schema.Charset{Uppercase: true, Lowercase: true},
			}

			var result string
			assert.NotPanics(t, func() { result = ExampleSecretFromComposition(v) })

			assert.Equal(t, tc.expectedLength, len(result))
			assert.True(t, strings.HasPrefix(result, tc.prefix))
			assert.Equal(t, tc.hasSuffix, strings.HasSuffix(result, secretExampleSuffix))
			assertOnlyCharsFrom(t, strings.TrimPrefix(result, tc.prefix), v.Charset)
		})
	}
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

func assertOnlyCharsFrom(t *testing.T, str string, c schema.Charset) {
	t.Helper()
	allowed := charsToUse(c)
	for _, r := range str {
		assert.True(t, strings.ContainsRune(allowed, r), fmt.Sprintf("character %q in %q is not in the charset", r, str))
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
