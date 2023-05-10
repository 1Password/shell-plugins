package credselect

import "github.com/1Password/shell-plugins/sdk/schema"

func New(selectionID string, opts ...CredentialSelectionOpt) schema.CredentialSelection {
	selection := schema.CredentialSelection{
		SelectionID: selectionID,
	}
	for _, opt := range opts {
		opt(selection)
	}
	return selection
}

func Any() CredentialSelectionOpt {
	return func(cs schema.CredentialSelection) {
		cs.Selector = any
	}
}

func PredefinedList(list []schema.CredentialUsage) CredentialSelectionOpt {
	return func(cs schema.CredentialSelection) {
		cs.PredefinedList = list
	}
}

func AllowMultiple() CredentialSelectionOpt {
	return func(cs schema.CredentialSelection) {
		cs.AllowMultiple = true
	}
}

const any = schema.CredentialSelector("any")

type CredentialSelectionOpt func(schema.CredentialSelection)

func Matches(cred schema.CredentialType, selector schema.CredentialSelector) bool {
	switch selector {
	case any:
		return true
	default:
		return false
	}
}
