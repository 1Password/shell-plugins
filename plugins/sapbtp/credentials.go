package sapbtp

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func Credentials() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.Credentials,
		DocsURL:       sdk.URL("https://help.sap.com/docs/BTP/65de2977205c403bbc107264b8eccf4b/cc1c676b43904066abb2a4838cbd0c37.html"),
		ManagementURL: sdk.URL("https://account.hana.ondemand.com/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Username,
				MarkdownDescription: " SAP BTP Username (email address)",
				Optional:            false,
				Secret:              false,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{Uppercase: true, Lowercase: true, Digits: true},
					Prefix:  "",
				},
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: " SAP BTP Password",
				Optional:            false,
				Secret:              true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{Uppercase: true, Lowercase: true, Digits: true},
					Prefix:  "",
				},
			},
		},
		DefaultProvisioner: BTPProvisioner(),
	}
}
