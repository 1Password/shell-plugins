package aws_cdk

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"

	"github.com/1Password/shell-plugins/plugins/aws"
)

func AccessKey() schema.CredentialType {
	return schema.CredentialType{
		Name:               credname.AccessKey,
		DocsURL:            sdk.URL("https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html"),
		ManagementURL:      sdk.URL("https://console.aws.amazon.com/iam"),
		Fields:             aws.AWSCredentialFieldSchema,
		DefaultProvisioner: aws.AWSProvisioner(),
		Importer:           aws.TryAWSImporters(),
	}
}
