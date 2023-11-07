package nomad

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func AuthToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AuthToken,
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to HashiCorp Nomad.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 36,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Address,
				MarkdownDescription: "Address of the HashiCorp Nomad Server",
				Secret:              true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryHashiCorpNomadConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"NOMAD_TOKEN":   fieldname.Token,
	"NOMAD_ADDRESS": fieldname.Address,
}

// TODO: Check if the platform stores the Auth Token in a local config file, and if so,
// implement the function below to add support for importing it.
func TryHashiCorpNomadConfigFile() sdk.Importer {
	return importer.NoOp()
}

// TODO: Implement the config file schema
// type Config struct {
//	Token string
// }
