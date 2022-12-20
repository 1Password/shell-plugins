package curl

import (
	"context"
	"fmt"

	"github.com/1Password/shell-plugins/sdk"
)

func CurlProvisioner(httpProvisioner sdk.HTTPProvisioner) sdk.Provisioner {
	return curlProvisioner{httpProvisioner: httpProvisioner}
}

type curlProvisioner struct {
	sdk.Provisioner

	httpProvisioner sdk.HTTPProvisioner
}

func (p curlProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	httpIn := sdk.HTTPProvisionInput{
		HomeDir:    in.HomeDir,
		Cache:      in.Cache,
		ItemFields: in.ItemFields,
	}

	httpOut := sdk.HTTPProvisionOutput{
		Headers:     map[string]string{},
		QueryParams: map[string]string{},
	}

	p.httpProvisioner.Provision(ctx, httpIn, &httpOut)

	if len(httpOut.QueryParams) > 0 {
		urlIndex := len(out.CommandLine) - 1
		url := out.CommandLine[urlIndex] // TODO: Find URL arg from curl command
		for queryParamName, queryParamValue := range httpOut.QueryParams {
			// TODO: Properly parse URL and add URL encoded query params
			url = fmt.Sprintf("%s?%s=%s", url, queryParamName, queryParamValue)
			out.CommandLine[urlIndex] = url
		}
	}

	for headerName, headerValue := range httpOut.Headers {
		out.AddArgs("-H", fmt.Sprintf("%s: %s", headerName, headerValue))
	}
}

func (p curlProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
}

func (p curlProvisioner) Description() string {
	return ""
}
