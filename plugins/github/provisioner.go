package github

import (
	"context"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

const (
	githubHost    = "github.com"
	localhostHost = "github.localhost"
	tenancyHost   = "ghe.com"
)

// GitHubCLIProvisioner provisions auth tokens using the environment variables expected by the GitHub CLI.
// GitHub.com and GitHub Enterprise Cloud (*.ghe.com) use GH_TOKEN and GITHUB_TOKEN.
// GitHub Enterprise Server uses GH_ENTERPRISE_TOKEN and GITHUB_ENTERPRISE_TOKEN with GH_HOST.
type GitHubCLIProvisioner struct{}

func (p GitHubCLIProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	token, ok := in.ItemFields[fieldname.Token]
	if !ok {
		return
	}

	host := strings.ToLower(strings.TrimSpace(in.ItemFields[fieldname.Host]))

	if usesEnterpriseToken(host) {
		out.AddEnvVar("GH_ENTERPRISE_TOKEN", token)
		out.AddEnvVar("GITHUB_ENTERPRISE_TOKEN", token)
		if host != "" {
			out.AddEnvVar("GH_HOST", host)
		}
		return
	}

	out.AddEnvVar("GH_TOKEN", token)
	out.AddEnvVar("GITHUB_TOKEN", token)
	if host != "" && host != githubHost {
		out.AddEnvVar("GH_HOST", host)
	}
}

func usesEnterpriseToken(host string) bool {
	if host == "" || host == githubHost || host == localhostHost {
		return false
	}

	return !isTenancyHost(host)
}

func isTenancyHost(host string) bool {
	return strings.HasSuffix(host, "."+tenancyHost)
}

func (p GitHubCLIProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Environment variables get wiped automatically when the process exits.
}

func (p GitHubCLIProvisioner) Description() string {
	return "Provision GitHub CLI environment variables: GH_TOKEN, GITHUB_TOKEN, GH_ENTERPRISE_TOKEN, GITHUB_ENTERPRISE_TOKEN, GH_HOST"
}
