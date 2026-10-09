package oci

import (
	"context"
	"os"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/1Password/shell-plugins/plugins/oci/schema/ocinames"
)

func AccessKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AccessKey,
		DocsURL:       sdk.URL("https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/clienvironmentvariables.htm"),
		ManagementURL: sdk.URL("console.oracle.com"),
		Fields: []schema.CredentialField{
			{
				Name: fieldname.SecretAccessKey,
				MarkdownDescription: "The full content of the private key enclosed in single quotes. Important: The key pair must be in PEM format.",
				Secret: true,
			},
			{
				Name: ocinames.FingerPrint,
				MarkdownDescription: "The fingerprint for the public key that was added to this user. To get the value, see [Required Keys and OCIDs](https://docs.oracle.com/en-us/iaas/Content/API/Concepts/apisigningkey.htm#Required_Keys_and_OCIDs).",
				Secret: false,
			},{
				Name: fieldname.User,
				MarkdownDescription: "The OCID of the user calling the API. To get the value",
				Secret: false,
			}, {
				Name: fieldname.Region,
				MarkdownDescription: "An Oracle Cloud Infrastructure region",
				Secret: false,
				Optional: true,
			}, {
				Name: ocinames.Tenancy,
				MarkdownDescription: "The OCID of your tenancy",
				Secret: false,
				Optional: false,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryOracleCloudConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"OCI_CLI_KEY_CONTENT": fieldname.SecretAccessKey,
	"OCI_CLI_FINGERPRINT": ocinames.FingerPrint,
	"OCI_CLI_TENANCY": ocinames.Tenancy,
	"OCI_CLI_USER": fieldname.User,
	"OCI_CLI_REGION": fieldname.Region,
}

// TODO: Check if the platform stores the Access Key in a local config file, and if so,
// implement the function below to add support for importing it.
func TryOracleCloudConfigFile() sdk.Importer {
	return importer.TryFile("~/oci/config", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToTOML(&config); err != nil {
			out.AddError(err)
			return
		}
	

		defaultConfig := config.DEFAULT
		if defaultConfig.keyfile == "" {
			return
		}
		
		bytes,_ := os.ReadFile(defaultConfig.keyfile);
		
		key_contents := string(bytes[:])

		 out.AddCandidate(sdk.ImportCandidate{
			
			 Fields: map[sdk.FieldName]string{
				 fieldname.SecretAccessKey: key_contents,
				 ocinames.FingerPrint: defaultConfig.fingerprint,
				 ocinames.Tenancy: defaultConfig.tenancy,
				 fieldname.User: defaultConfig.user,
				 fieldname.Region: defaultConfig.region,
		 	},
		 })
	})
}

 type (
 	Config struct {
		DEFAULT ConfigContent
 	}
	ConfigContent struct {
		fingerprint string
		tenancy string
		region string
		user string
		keyfile string
	} 
 )
