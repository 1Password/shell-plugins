package cvelib

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func APIKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIKey,
		DocsURL:       sdk.URL("https://github.com/RedHatProductSecurity/cvelib"),
		ManagementURL: sdk.URL("https://vulnogram.github.io/cve5/#cvePortal"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.User,
				MarkdownDescription: "User to authenticate to CVE Services API (CVE user).",
			},
			{
				Name:                fieldname.Organization,
				MarkdownDescription: "Organization to authenticate to CVE Services API (CNA short name).",
			},
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to CVE Services API (CNA API key).",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 36,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"CVE_USER":    fieldname.User,
	"CVE_ORG":     fieldname.Organization,
	"CVE_API_KEY": fieldname.APIKey,
}
