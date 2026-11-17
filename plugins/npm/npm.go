package npm

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func NPMCLI() schema.Executable {
	return packageManagerCLI("npm CLI", "npm", "https://docs.npmjs.com/cli/")
}

func PNPMCLI() schema.Executable {
	executable := packageManagerCLI("pnpm CLI", "pnpm", "https://pnpm.io/cli")
	executable.Uses[0].Provisioner = pnpmProvisioner()
	return executable
}

func packageManagerCLI(name, command, docsURL string) schema.Executable {
	return schema.Executable{
		Name:    name,
		Runs:    []string{command},
		DocsURL: sdk.URL(docsURL),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.IfAny(
				needsauth.ForCommand("access"),
				needsauth.ForCommand("add"),
				needsauth.ForCommand("audit"),
				needsauth.ForCommand("ci"),
				needsauth.ForCommand("cit"),
				needsauth.ForCommand("deprecate"),
				needsauth.ForCommand("dist-tag"),
				needsauth.ForCommand("fetch"),
				needsauth.ForCommand("i"),
				needsauth.ForCommand("exec"),
				needsauth.ForCommand("install"),
				needsauth.ForCommand("install-ci-test"),
				needsauth.ForCommand("install-test"),
				needsauth.ForCommand("it"),
				needsauth.ForCommand("owner"),
				needsauth.ForCommand("pack"),
				needsauth.ForCommand("profile"),
				needsauth.ForCommand("publish"),
				needsauth.ForCommand("search"),
				needsauth.ForCommand("stage"),
				needsauth.ForCommand("team"),
				needsauth.ForCommand("token"),
				needsauth.ForCommand("unpublish"),
				needsauth.ForCommand("update"),
				needsauth.ForCommand("view"),
				needsauth.ForCommand("whoami"),
			),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AccessToken,
			},
		},
	}
}
