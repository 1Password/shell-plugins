package http

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
)

type httpHeadersProvisioner struct {
	sdk.HTTPProvisioner

	headers map[string]ValueFunc
}

func HTTPHeadersProvisioner(headers map[string]ValueFunc) sdk.HTTPProvisioner {
	return httpHeadersProvisioner{
		headers: headers,
	}
}

func (p httpHeadersProvisioner) Description() string {
	return ""
}

func (p httpHeadersProvisioner) Provision(ctx context.Context, in sdk.HTTPProvisionInput, out *sdk.HTTPProvisionOutput) {
	for name, valueFunc := range p.headers {
		value, err := valueFunc(ctx, in)
		if err != nil {
			// TODO: handle error
		}
		out.Headers[name] = value
	}
}
