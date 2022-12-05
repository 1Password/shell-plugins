package plugintest

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/1Password/shell-plugins/sdk/schema"
)

const (
	lowerCaseLetters    = "abcdefghijklmnopqrstuvwxyz"
	capitalCaseLetters  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits              = "0123456789"
	symbols             = "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"
	secretExampleSuffix = "EXAMPLE"
)

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func ExampleSecretFromComposition(v schema.ValueComposition) string {
	prefix := getPrefix(v)
	suffix := getSuffix(v)
	base := generateBase(v, v.Length-len(prefix)-len(suffix))

	return prefix + base + suffix
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

func getSuffix(v schema.ValueComposition) string {
	var suffix string

	if v.Length > len(secretExampleSuffix) && (v.Charset.Uppercase || v.Charset.Lowercase) {
		suffix = secretExampleSuffix
		if v.Charset.Lowercase && !v.Charset.Uppercase {
			suffix = strings.ToLower(secretExampleSuffix)
		}
	}

	return suffix
}

func stringFromCharset(length int, charset string) (string, error) {
	if charset == "" {
		return "", fmt.Errorf("invalid charset provided")
	}

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
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
