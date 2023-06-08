package atlas

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

const (
	index uint = 1
)

func DatabaseCredentials() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.DatabaseCredentials,
		DocsURL: sdk.URL("https://www.mongodb.com/docs/mongodb-shell/connect/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Host,
				MarkdownDescription: "Host address for the MongoDB database.",
				Secret:              false,
				Optional:            true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
						Specific:  []rune{'.', '-'},
					},
				},
			},
			{
				Name:                fieldname.Port,
				MarkdownDescription: "Port for the MongoDB database.",
				Secret:              false,
				Optional:            true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Digits: true,
					},
				},
			},
			{
				Name:                fieldname.Username,
				MarkdownDescription: "Username for authenticating to the MongoDB database.",
				Secret:              false,
				Optional:            true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
						Specific:  []rune{'-', '_'},
					},
				},
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password for authenticating to the MongoDB database.",
				Secret:              true,
				Optional:            true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: nil,
		Importer:           nil,
	}
}

var argsToProvision = map[string]sdk.FieldName{
	"--host":     fieldname.Host,
	"--port":     fieldname.Port,
	"--username": fieldname.Username,
	"--password": fieldname.Password,
}

var indexToProvisionAt = map[string]uint{
	"--host":     index,
	"--port":     index,
	"--username": index,
	"--password": index,
}
