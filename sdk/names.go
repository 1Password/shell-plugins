package sdk

import (
	"regexp"
	"strings"
)

// FieldName represents a name of credential field. It should be title-cased.
// Examples: "Password", "Token", "API Key".
type FieldName string

func (n FieldName) String() string {
	return string(n)
}

// CredentialName represents a name of a credential type. It should be title-cased.
// Examples: "Personal Access Token", "API Key".
type CredentialName string

func (n CredentialName) String() string {
	return string(n)
}

func (n CredentialName) ID() CredentialTypeID {
	return CredentialTypeID(credentialNameToSnakeCase(n))
}

func credentialNameToSnakeCase(name CredentialName) string {
	const underscore = "_"

	str := name.String()

	charsToReplace := regexp.MustCompile(`[-/,. ]`)
	str = charsToReplace.ReplaceAllLiteralString(str, underscore)

	multipleUnderscores := regexp.MustCompile(`_+`)
	str = multipleUnderscores.ReplaceAllString(str, underscore)

	str = strings.TrimPrefix(str, "_")
	str = strings.TrimSuffix(str, "_")

	return strings.ToLower(str)
}

type CredentialTypeID string

func (i CredentialTypeID) String() string {
	return string(i)
}
