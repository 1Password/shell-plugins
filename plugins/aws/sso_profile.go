package aws

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func SSOProfile() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.SSOProfile,
		DocsURL:       sdk.URL("https://docs.aws.amazon.com/cli/latest/userguide/sso-configure-profile-token.html"),
		ManagementURL: sdk.URL("https://console.aws.amazon.com/iamidentitycenter"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.SSOStartURL,
				MarkdownDescription: "The AWS access portal URL for your organization, e.g. `https://your-org.awsapps.com/start`. Found in the IAM Identity Center console under Settings.",
			},
			{
				Name:                fieldname.SSORegion,
				MarkdownDescription: "The AWS region where IAM Identity Center is configured (e.g. `us-east-1`). Shown next to the access portal URL in the IAM Identity Center console.",
			},
			{
				Name:                fieldname.SSOAccountID,
				MarkdownDescription: "The 12-digit AWS account ID to sign in to. Listed next to each account in the AWS access portal.",
				Composition: &schema.ValueComposition{
					Length: 12,
					Charset: schema.Charset{
						Digits: true,
					},
				},
			},
			{
				Name:                fieldname.SSORoleName,
				MarkdownDescription: "The role to assume in the account, e.g. `AdministratorAccess`. Listed under each account in the AWS access portal.",
			},
			{
				Name:                fieldname.SSOSession,
				MarkdownDescription: "Optional. The shared `sso-session` name from your `~/.aws/config` file. Leave blank if you use the legacy per-profile format.",
				Optional:            true,
			},
			{
				Name:                fieldname.DefaultRegion,
				MarkdownDescription: "Optional. The AWS region to use for API calls (e.g. `us-east-1`).",
				Optional:            true,
			},
		},
		DefaultProvisioner: NewSSOProvisioner(""),
		Importer: importer.TryAll(
			TrySSOConfigFile(),
		),
		// The SSO bearer token lives in `~/.aws/sso/cache/<sha1>.json` (written by `aws sso login`),
		// so the vault item only stores configuration. Opt out of the "must have at least one secret
		// field" validator instead of relaxing it globally.
		AllowsExternalSecretCache: true,
	}
}
