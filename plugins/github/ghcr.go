package github

import (
	"github.com/1Password/shell-plugins/plugins/docker"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func GHCRProvisioner() sdk.Provisioner {
	return docker.ConfigFileProvisioner(
		docker.WithURLField(fieldname.URL),
		docker.WithStaticURL("ghcr.io"),
		docker.WithUsernameField(fieldname.Username),
		docker.WithPasswordField(fieldname.Token),
	)
}
