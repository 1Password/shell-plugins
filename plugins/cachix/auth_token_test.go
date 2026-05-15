package cachix

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAuthTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AuthToken().Importer, map[string]plugintest.ImportCase{
		"Cache auth token": {
			Environment: map[string]string{
				"CACHIX_AUTH_TOKEN": "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI2MGI1ZTJlNy03ZDZiLTRiZGYtYjhiMS1mZGU2NDgzMmY5YzgiLCJzY29wZXMiOiJjYWNoZSJ9.PXNGrCN7ovMgEK0haz9voQfeECCwzzD79mEwg9KjqVE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI2MGI1ZTJlNy03ZDZiLTRiZGYtYjhiMS1mZGU2NDgzMmY5YzgiLCJzY29wZXMiOiJjYWNoZSJ9.PXNGrCN7ovMgEK0haz9voQfeECCwzzD79mEwg9KjqVE",
					},
				},
			},
		},
		"Personal auth token": {
			Environment: map[string]string{
				"CACHIX_AUTH_TOKEN": "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI2MGI1ZTJlNy03ZDZiLTRiZGYtYjhiMS1mZGU2NDgzMmY5YzgiLCJzY29wZXMiOiJ0eCJ9.8u5huhMCEX8v58kp6oCU6ueJ8-EXzMNnpH5ERBKEabs",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI2MGI1ZTJlNy03ZDZiLTRiZGYtYjhiMS1mZGU2NDgzMmY5YzgiLCJzY29wZXMiOiJ0eCJ9.8u5huhMCEX8v58kp6oCU6ueJ8-EXzMNnpH5ERBKEabs",
					},
				},
			},
		},
		"Cache config file": {
			Files: map[string]string{
				"~/.config/cachix/cachix.dhall": plugintest.LoadFixture(t, "cachix.dhall"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI3OGRlOTA4Yi03MjhiLTRiMjUtODA1Yi1hNTRjODIxMWQ1ZjMiLCJzY29wZXMiOiJ0eCJ9.A0XjByVJtp2Di0Ui7M5KjiG1OinYW8PwVKRw5N4YETE",
					},
				},
			},
		},
	})
}

func TestAuthTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AuthToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI2MGI1ZTJlNy03ZDZiLTRiZGYtYjhiMS1mZGU2NDgzMmY5YzgiLCJzY29wZXMiOiJjYWNoZSJ9.PXNGrCN7ovMgEK0haz9voQfeECCwzzD79mEwg9KjqVE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"CACHIX_AUTH_TOKEN": "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI2MGI1ZTJlNy03ZDZiLTRiZGYtYjhiMS1mZGU2NDgzMmY5YzgiLCJzY29wZXMiOiJjYWNoZSJ9.PXNGrCN7ovMgEK0haz9voQfeECCwzzD79mEwg9KjqVE",
				},
			},
		},
	})
}
