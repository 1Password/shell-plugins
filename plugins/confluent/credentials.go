package confluent

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func CloudCredentials() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.Credentials,
		DocsURL:       sdk.URL("https://docs.confluent.io/confluent-cli/current/command-reference/confluent_login.html"),
		ManagementURL: sdk.URL("https://confluent.cloud"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Username,
				MarkdownDescription: "Email used to authenticate to Confluent Cloud.",
				Secret:              false,
				Optional:            false,
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password used to authenticate to Confluent Cloud.",
				Secret:              true,
				Optional:            false,
			},
			{
				Name:                fieldname.Organization,
				MarkdownDescription: "Organization to use with Confluent Cloud.",
				Secret:              false,
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultCloudEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultCloudEnvVarMapping),
		)}
}

var defaultCloudEnvVarMapping = map[string]sdk.FieldName{
	"CONFLUENT_CLOUD_EMAIL":           fieldname.Username,
	"CONFLUENT_CLOUD_PASSWORD":        fieldname.Password,
	"CONFLUENT_CLOUD_ORGANIZATION_ID": fieldname.Organization,
}

//TODO: Revist once the Shell Plugins ecosystem adds support for multiple credential types per plugin
// func PlatformCredentials() schema.CredentialType {
// 	return schema.CredentialType{
// 		Name:    credname.Credentials,
// 		DocsURL: sdk.URL("https://docs.confluent.io/confluent-cli/current/command-reference/confluent_login.html"),
// 		Fields: []schema.CredentialField{
// 			{
// 				Name:                fieldname.Username,
// 				MarkdownDescription: "Username used to authenticate to Confluent Platform.",
// 				Secret:              false,
// 				Optional:            false,
// 			},
// 			{
// 				Name:                fieldname.Password,
// 				MarkdownDescription: "Password used to authenticate to Confluent Platform.",
// 				Secret:              true,
// 				Optional:            false,
// 			},
// 			{
// 				Name:                fieldname.URL,
// 				MarkdownDescription: "Metadata Service (MDS) URL used to authenticate to Confluent Platform.",
// 				Secret:              false,
// 				Optional:            false,
// 			},
// 			{
// 				Name:                fieldname.Certificate,
// 				MarkdownDescription: "Self-signed certificate chain in PEM format.",
// 				Secret:              true,
// 				Optional:            false,
// 			},
// 		},
// 		DefaultProvisioner: provision.EnvVars(defaultPlatformEnvVarMapping),
// 		Importer: importer.TryAll(
// 			importer.TryEnvVarPair(defaultPlatformEnvVarMapping),
// 		)}
// }

// var defaultPlatformEnvVarMapping = map[string]sdk.FieldName{
// 	"CONFLUENT_PLATFORM_USERNAME":     fieldname.Username,
// 	"CONFLUENT_PLATFORM_PASSWORD":     fieldname.Password,
// 	"CONFLUENT_PLATFORM_MDS_URL":      fieldname.URL,
// 	"CONFLUENT_PLATFORM_CA_CERT_PATH": fieldname.Certificate,
// }
