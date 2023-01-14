package argocd

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAuthTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AuthToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AuthToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcmdvY2QiLCJzdWIiOiJhZG1pbjphcGlLZXkiLCJuYmYiOjE2NzM3MjE0MDMsImlhdCI6MTY3MzcyMTQwMywianRpIjoiNTI5ODcwMTEtNGRiNy00ZmIxLWE4Y2MtMTk5NGViYTRjZDU0In0.ApWbS8sQbtp1l_ILaRO94izPv9AML2vnv1F3EXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"ARGOCD_AUTH_TOKEN": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcmdvY2QiLCJzdWIiOiJhZG1pbjphcGlLZXkiLCJuYmYiOjE2NzM3MjE0MDMsImlhdCI6MTY3MzcyMTQwMywianRpIjoiNTI5ODcwMTEtNGRiNy00ZmIxLWE4Y2MtMTk5NGViYTRjZDU0In0.ApWbS8sQbtp1l_ILaRO94izPv9AML2vnv1F3EXAMPLE",
				},
			},
		},
	})
}

func TestAuthTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AuthToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"ARGOCD_AUTH_TOKEN": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcmdvY2QiLCJzdWIiOiJhZG1pbjphcGlLZXkiLCJuYmYiOjE2NzM3MjE0MDMsImlhdCI6MTY3MzcyMTQwMywianRpIjoiNTI5ODcwMTEtNGRiNy00ZmIxLWE4Y2MtMTk5NGViYTRjZDU0In0.ApWbS8sQbtp1l_ILaRO94izPv9AML2vnv1F3EXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AuthToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcmdvY2QiLCJzdWIiOiJhZG1pbjphcGlLZXkiLCJuYmYiOjE2NzM3MjE0MDMsImlhdCI6MTY3MzcyMTQwMywianRpIjoiNTI5ODcwMTEtNGRiNy00ZmIxLWE4Y2MtMTk5NGViYTRjZDU0In0.ApWbS8sQbtp1l_ILaRO94izPv9AML2vnv1F3EXAMPLE",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.config/argocd/config": plugintest.LoadFixture(t, "config"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					NameHint: "test-context",
					Fields: map[sdk.FieldName]string{
						fieldname.AuthToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcmdvY2QiLCJzdWIiOiJhZG1pbjphcGlLZXkiLCJuYmYiOjE2NzM3MjE0MDMsImlhdCI6MTY3MzcyMTQwMywianRpIjoiNTI5ODcwMTEtNGRiNy00ZmIxLWE4Y2MtMTk5NGViYTRjZDU0In0.ApWbS8sQbtp1l_ILaRO94izPv9AML2vnv1F3EXAMPLE",
						fieldname.Address:   "argocd.test.domain",
					},
				},
			},
		},
	})
}
