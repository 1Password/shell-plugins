package sdk

import (
	"context"
)

type HTTPProvisioner interface {
	Description() string

	Provision(ctx context.Context, input HTTPProvisionInput, output *HTTPProvisionOutput)
}

type HTTPProvisionInput struct {
	// HomeDir is the path to current user's home directory.
	HomeDir string

	// Cache can contain data that got added in the provision step from previous runs for this credential.
	Cache CacheState

	// ItemFields contains the field names and their corresponding (sensitive) values.
	ItemFields map[FieldName]string
}

// HTTPProvisionOutput contains the sensitive values that the HTTPProvisioner outputs.
type HTTPProvisionOutput struct {
	Headers     map[string]string
	QueryParams map[string]string

	Cache CacheOperations

	Diagnostics Diagnostics
}
