package aws_cdk

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "aws_cdk",
		Platform: schema.PlatformInfo{
			Name:     "AWS CDK",
			Homepage: sdk.URL("https://aws_cdk.com"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			AccessKey(),
		},
		Executables: []schema.Executable{
			AWSCDKCLI(),
		},
	}
}
