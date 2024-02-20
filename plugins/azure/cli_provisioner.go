package azure

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

type CLIProvisioner struct {
}

func (p CLIProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	args := []string{
		"login", "--service-principal",
		"--user", in.ItemFields[fieldname.ClientID],
		"--password", in.ItemFields[fieldname.ClientSecret],
	}

	if tid, ok := in.ItemFields[fieldname.TenantID]; ok {
		args = append(args, []string{"--tenant", tid}...)
	}

	cmd := exec.Command("az", args...)
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func (p CLIProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	cmd := exec.Command("az", "logout")
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func (p CLIProvisioner) Description() string {
	return "Log into Azure using a service principal"
}

func ExecuteSilently[input interface{}, output interface{}, e error](f func(input) (output, e)) func(input) (output, e) {
	return func(i input) (output, e) {
		log.SetOutput(io.Discard)
		defer log.SetOutput(os.Stderr)
		return f(i)
	}
}
