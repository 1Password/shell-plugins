package sops

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAgeSecretKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, AgeSecretKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"SOPS_AGE_KEY": "AGE-SECRET-KEY-1QFWENTHXAAPACFPMXHQCREP64GJE5YTHXLX0RPFSXRSPDJGCR0SSWYNX3D",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.PrivateKey: "AGE-SECRET-KEY-1QFWENTHXAAPACFPMXHQCREP64GJE5YTHXLX0RPFSXRSPDJGCR0SSWYNX3D",
					},
				},
			},
		},
	})
}

func TestAgeSecretKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AgeSecretKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.PrivateKey: "AGE-SECRET-KEY-1QFWENTHXAAPACFPMXHQCREP64GJE5YTHXLX0RPFSXRSPDJGCR0SSWYNX3D",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SOPS_AGE_KEY": "AGE-SECRET-KEY-1QFWENTHXAAPACFPMXHQCREP64GJE5YTHXLX0RPFSXRSPDJGCR0SSWYNX3D",
				},
			},
		},
	})
}
