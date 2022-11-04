package strgenerator

import (
	"fmt"
	"math/rand"
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

func ExampleSecretFromComposition(v *schema.ValueComposition) (string, error) {
	var exampleStr string
	chars := charsToUse(v.Charset)

	if v.Prefix != "" {
		exampleStr += v.Prefix
	}
	generatedStr, err := stringFromCharset(v.Length-len(v.Prefix)-len(secretExampleSuffix), chars)
	if err != nil {
		return "", err
	}
	exampleStr += generatedStr
	exampleStr += secretExampleSuffix

	return exampleStr, nil
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

	if c.Uppercase == true {
		chars += capitalCaseLetters
	}

	if c.Lowercase == true {
		chars += lowerCaseLetters
	}

	if c.Digits == true {
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
