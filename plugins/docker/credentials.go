package docker

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func Credentials() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.Credentials,
		DocsURL:       sdk.URL("https://docs.docker.com/engine/reference/commandline/login"),
		ManagementURL: sdk.URL("https://hub.docker.com/settings/security"),
		Fields: []schema.CredentialField{
			{
				Name: fieldname.Username,
				MarkdownDescription: "Username used in Docker registries.",
				Secret: false,
			},
			{
				Name:                fieldname.Secret,
				AlternativeNames: []string{fieldname.Password.String(), fieldname.AccessToken.String()},
				MarkdownDescription: "Password or access token used to authenticate to a Docker registry.",
				Secret:              true,
			},
			{
				Name: fieldname.Host,
				MarkdownDescription: "URL of the Docker registry server.",
				Optional: true, // Defaults to Docker Hub registry
				Secret: false,
			},
		},
		DefaultProvisioner: dockerProvisioner{},
		Importer: importer.TryAll(
			TryDockerConfigFile(),
		)}
}

func TryDockerConfigFile() sdk.Importer {
	// The config file likely points to a keychain, but might hold the plaintext credentials.
	return importer.TryFile("~/.docker/config.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		 var config Config
		 if err := contents.ToJSON(&config); err != nil {
		 	out.AddError(err)
		 	return
		 }

		 if config.Username == "" || config.Secret == "" {
		 	return
		 }

		 out.AddCandidate(sdk.ImportCandidate{
		 	Fields: map[sdk.FieldName]string{
		 		fieldname.Username: config.Username,
				fieldname.Secret: config.Secret,
				fieldname.Host: config.ServerUrl,
		 	},
		 })
	})
}

type dockerProvisioner struct{}

func (p dockerProvisioner) Description() string {
	return "Docker login credentials provisioner"
}

func (p dockerProvisioner) Provision(ctx context.Context, input sdk.ProvisionInput, output *sdk.ProvisionOutput) {
	if registry := input.ItemFields[fieldname.Host]; registry != "" {
		output.AddArgs(registry)
	}
	// Docker will emit a warning that using --password is insecure, as it in unaware that the password is not being
	// typed in manually on the command line
	output.AddArgs("--username", input.ItemFields[fieldname.Username], "--password", input.ItemFields[fieldname.Secret])
}

func (p dockerProvisioner) Deprovision(ctx context.Context, input sdk.DeprovisionInput, output *sdk.DeprovisionOutput) {
	// No-op: nothing to delete
}

type Config struct {
	Username string `json:"Username"`
	Secret string `json:"Secret"`
	ServerUrl string `json:"ServerURL"`
}
