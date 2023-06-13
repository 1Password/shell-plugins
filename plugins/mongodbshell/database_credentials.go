package mongodbshell

import (
	"fmt"

	"github.com/1Password/shell-plugins/sdk"
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

var argsToProvision = []string{
	fmt.Sprintf("--host %s", fieldname.Host),
	fmt.Sprintf("--port %s", fieldname.Port),
	fmt.Sprintf("--username %s", fieldname.Username),
	fmt.Sprintf("--password %s", fieldname.Password),
}
