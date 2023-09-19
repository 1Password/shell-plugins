package mongodb

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func DatabaseCredentials() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.DatabaseCredentials,
		DocsURL: sdk.URL("https://www.mongodb.com/docs/mongodb-shell/connect/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.ConnectionString,
				MarkdownDescription: "Connection String for the MongoDB database.",
				Secret:              false,
				Optional:            true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
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
				Optional:            false,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.NoOp(),
		Importer:           importer.NoOp(),
	}
}
