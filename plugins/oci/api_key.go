package oci

import (
	"context"
	"os"
	"path/filepath"
	"strings"

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
		DocsURL:       sdk.URL("https://docs.oracle.com/en-us/iaas/Content/API/Concepts/apisigningkey.htm"),
		ManagementURL: sdk.URL("https://cloud.oracle.com"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.User,
				MarkdownDescription: "OCID of the user used to authenticate to OCI.",
			},
			{
				Name:                tenantIDField,
				MarkdownDescription: "OCID of the tenant used to authenticate to OCI.",
			},
			{
				Name:                fingerprintField,
				MarkdownDescription: "Fingerprint of the API key used to authenticate to OCI.",
			},
			{
				Name:                fieldname.Region,
				MarkdownDescription: "OCI region to use.",
				Optional:            true,
			},
			{
				Name:                fieldname.PrivateKey,
				MarkdownDescription: "Private API key used to authenticate to OCI.",
				Secret:              true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryOCIConfigFile(),
		)}
}

var (
	tenantIDField    = sdk.FieldName("Tenant ID")
	fingerprintField = sdk.FieldName("Fingerprint")
)

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"OCI_CLI_USER":        fieldname.User,
	"OCI_CLI_TENANCY":     tenantIDField,
	"OCI_CLI_FINGERPRINT": fingerprintField,
	"OCI_CLI_REGION":      fieldname.Region,
	"OCI_CLI_KEY_CONTENT": fieldname.PrivateKey,
}

func TryOCIConfigFile() sdk.Importer {
	return importer.TryFile("~/.oci/config", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		config, err := contents.ToINI()
		if err != nil {
			out.AddError(err)
			return
		}

		for _, section := range config.Sections() {
			fields := map[sdk.FieldName]string{
				fieldname.User:   section.Key("user").String(),
				tenantIDField:    section.Key("tenancy").String(),
				fingerprintField: section.Key("fingerprint").String(),
			}

			if region := section.Key("region").String(); region != "" {
				fields[fieldname.Region] = region
			}

			keyFile := section.Key("key_file").String()
			if fields[fieldname.User] == "" || fields[tenantIDField] == "" || fields[fingerprintField] == "" || keyFile == "" {
				continue
			}

			if after, ok := strings.CutPrefix(keyFile, "~/"); ok {
				keyFile = in.FromHomeDir(after)
			} else if filepath.IsAbs(keyFile) {
				keyFile = in.FromRootDir(keyFile)
			}

			privateKey, err := os.ReadFile(keyFile)
			if err != nil {
				out.AddError(err)
				continue
			}
			fields[fieldname.PrivateKey] = string(privateKey)

			out.AddCandidate(sdk.ImportCandidate{
				NameHint: section.Name(),
				Fields:   fields,
			})
		}
	})
}
