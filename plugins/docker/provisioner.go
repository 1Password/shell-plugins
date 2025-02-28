package docker

import (
	"context"
	"github.com/1Password/shell-plugins/sdk"
)

type DockerProvisioner struct {
}

func dockerConfig(in sdk.ProvisionInput) ([]byte, error) {
	return []byte(`{
        "auths": {
                "https://index.docker.io/v1/": {
                        "credHelper": "envvars"
                }
        },
        "currentContext": "desktop-linux"
}`), nil
}

func (d DockerProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (d DockerProvisioner) Description() string {
	return "Provision environment variables with temporary STS credentials AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN"
}
