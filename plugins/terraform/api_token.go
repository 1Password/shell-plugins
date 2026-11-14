package terraform

import (
	"context"
	"fmt"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func APIToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIToken,
		DocsURL:       sdk.URL("https://developer.hashicorp.com/terraform/cloud-docs/users-teams-organizations/users#tokens"),
		ManagementURL: sdk.URL("https://app.terraform.io/app/settings/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Terraform Cloud.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 90,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Host,
				MarkdownDescription: "Hostname used to authenticate to which Terraform Cloud instance (default: app.terraform.io).",
			},
		},
		DefaultProvisioner: &TerraformProvisioner{},
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryTerraformConfigFile(),
		)}
}

type TerraformProvisioner struct {
	sdk.Provisioner

	Schema map[string]sdk.FieldName
}

func (p *TerraformProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	if hostValue, ok := in.ItemFields[fieldname.Host]; ok {
		if tokenValue, ok := in.ItemFields[fieldname.Token]; ok {
			envVarName := fmt.Sprintf("TF_TOKEN_%s", strings.ReplaceAll(hostValue, ".", "_"))
			out.AddEnvVar(envVarName, tokenValue)
		}
	}
}

func (p *TerraformProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p TerraformProvisioner) Description() string {
	var envVarNames []string
	for envVarName := range p.Schema {
		envVarNames = append(envVarNames, envVarName)
	}

	return fmt.Sprintf("Provision environment variables: %s", strings.Join(envVarNames, ", "))
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
    "TF_TOKEN_app_terraform_io": fieldname.Token,
}

func TryTerraformConfigFile() sdk.Importer {
	return importer.TryFile("~/.terraform.d/credentials.tfrc.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		if len(config.Credentials) == 0 {
			return
		}

		for host, cred := range config.Credentials {
			out.AddCandidate(sdk.ImportCandidate{
				Fields: map[sdk.FieldName]string{
					fieldname.Host:  host,
					fieldname.Token: cred.Token,
				},
			})
		}
	})
}

type Config struct {
	Credentials map[string]Credential `json:"credentials"`
}

type Credential struct {
	Token string `json:"token"`
}
