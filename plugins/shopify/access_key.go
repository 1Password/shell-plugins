package shopify

import (
	"context"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func AccessKey() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.CLIToken,
		DocsURL: sdk.URL("https://shopify.dev/docs/themes/tools/cli"),
		//ManagementURL: sdk.URL("https://admin.shopify.com/store/{YOUR_STORE_ID}/apps/theme-kit-access"), // Can't support unless URL supports variable URLs
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Shopify Themes.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 39,
					Prefix: "shptka_",
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: shopifyThemeProvisioner{},
		// Can't Implement Importer as the Shopify environment variables are stored by project directory instead of a fixed location
		// See: https://shopify.dev/docs/themes/tools/cli/environments
	}
}

type shopifyThemeProvisioner struct{}

func (v shopifyThemeProvisioner) Description() string {
	return "Shopify Theme CLI password provisioner"
}

func (v shopifyThemeProvisioner) Provision(ctx context.Context, input sdk.ProvisionInput, output *sdk.ProvisionOutput) {
	output.AddArgs("--password", input.ItemFields[fieldname.Token])
}

func (v shopifyThemeProvisioner) Deprovision(ctx context.Context, input sdk.DeprovisionInput, output *sdk.DeprovisionOutput) {
	// No Operator
}
