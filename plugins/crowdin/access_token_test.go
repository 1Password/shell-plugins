package crowdin

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token:       "pt2sk1g1nqfervdjk8av8r88wvjsprzzck4vws97n5r7fn8ad0l3pe45k09fecy152ra6c8sgexample",
				fieldname.ProjectID:   "123",
				fieldname.HostAddress: "https://testOrg.crowdin.com",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"CROWDIN_PERSONAL_TOKEN": "pt2sk1g1nqfervdjk8av8r88wvjsprzzck4vws97n5r7fn8ad0l3pe45k09fecy152ra6c8sgexample",
					"CROWDIN_PROJECT_ID":     "123",
					"CROWDIN_BASE_URL":       "https://testOrg.crowdin.com",
				},
			},
		},
	})
}

func TestAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"CROWDIN_PERSONAL_TOKEN": "pt2sk1g1nqfervdjk8av8r88wvjsprzzck4vws97n5r7fn8ad0l3pe45k09fecy152ra6c8sgexample",
				"CROWDIN_PROJECT_ID":     "123",
				"CROWDIN_BASE_URL":       "https://testOrg.crowdin.com",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:       "pt2sk1g1nqfervdjk8av8r88wvjsprzzck4vws97n5r7fn8ad0l3pe45k09fecy152ra6c8sgexample",
						fieldname.ProjectID:   "123",
						fieldname.HostAddress: "https://testOrg.crowdin.com",
					},
				},
			},
		},
	})
}
