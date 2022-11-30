package sdk

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
