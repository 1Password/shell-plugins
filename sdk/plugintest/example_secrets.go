package plugintest

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/1Password/shell-plugins/sdk/schema"
)

const (
	lowerCaseLetters    = "abcdefghijklmnopqrstuvwxyz"
	capitalCaseLetters  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits              = "0123456789"
	symbols             = "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"
	secretExampleSuffix = "EXAMPLE"
	defaultBodyLength   = 24
)

// ExampleSecretFromComposition returns a Length-character value (just the prefix if that's longer).
// Unset Length: prefix + defaultBodyLength random characters + EXAMPLE if the charset has letters.
func ExampleSecretFromComposition(v schema.ValueComposition) string {
	prefix := getPrefix(v)
	suffix := getSuffix(v.Charset)

	if v.Length == 0 {
		return prefix + generateBase(v, defaultBodyLength) + suffix
	}

	// Only add the suffix if at least one random character still fits.
	if v.Length <= len(prefix)+len(suffix) {
		suffix = ""
	}

	baseLength := v.Length - len(prefix) - len(suffix)
	if baseLength < 0 {
		baseLength = 0
	}

	return prefix + generateBase(v, baseLength) + suffix
}

func getPrefix(v schema.ValueComposition) string {
	if v.Prefix != "" {
		return v.Prefix
	}

	return ""
}

func generateBase(v schema.ValueComposition, baseLength int) string {
	chars := charsToUse(v.Charset)
	generatedStr, err := stringFromCharset(baseLength, chars)

	if err != nil {
		log.Fatalf("Error while generating secret: %v", err)
	}

	return generatedStr
}

func getSuffix(c schema.Charset) string {
	if c.Uppercase {
		return secretExampleSuffix
	}

	if c.Lowercase {
		return strings.ToLower(secretExampleSuffix)
	}

	return ""
}

func stringFromCharset(length int, charset string) (string, error) {
	if charset == "" {
		return "", fmt.Errorf("invalid charset provided")
	}
	max := big.NewInt(int64(len(charset)))
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("reading random: %w", err)
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}

func charsToUse(c schema.Charset) string {
	var chars string

	if c.Uppercase {
		chars += capitalCaseLetters
	}

	if c.Lowercase {
		chars += lowerCaseLetters
	}

	if c.Digits {
		chars += digits
	}

	if c.Symbols {
		chars += symbols
	}

	if len(c.Specific) > 0 {
		for _, r := range c.Specific {
			chars += string(r)
		}
	}

	return chars
}
