package sdk

import "strings"

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
	str := name.String()
	str = strings.ReplaceAll(str, " ", "_")
	str = strings.ReplaceAll(str, "-", "_")
	str = strings.ReplaceAll(str, "/", "_")
	return strings.ToLower(str)
}

type CredentialTypeID string

func (i CredentialTypeID) String() string {
	return string(i)
}
