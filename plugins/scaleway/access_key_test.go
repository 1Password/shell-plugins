package scaleway

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AccessKeyID:     "AABBCCDDEEFFGGHHIIJJ",
				fieldname.SecretAccessKey: "c3a22663-2770-4428-8166-c214643cd70b",
				fieldname.DefaultRegion:   "fr-par",
				fieldname.DefaultZone:     "fr-par-1",
				fieldname.ProjectID:       "d3a22663-2770-4428-8166-c214643cd70c",
				fieldname.OrgID:           "e3a22663-2770-4428-8166-c214643cd70d",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SCW_ACCESS_KEY":              "AABBCCDDEEFFGGHHIIJJ",
					"SCW_SECRET_KEY":              "c3a22663-2770-4428-8166-c214643cd70b",
					"SCW_DEFAULT_REGION":          "fr-par",
					"SCW_DEFAULT_ZONE":            "fr-par-1",
					"SCW_DEFAULT_PROJECT_ID":      "d3a22663-2770-4428-8166-c214643cd70c",
					"SCW_DEFAULT_ORGANIZATION_ID": "e3a22663-2770-4428-8166-c214643cd70d",
				},
			},
		},
	})
}

func TestAccessKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"SCW_ACCESS_KEY":              "yolo",
				"SCW_SECRET_KEY":              "c3a22663-2770-4428-8166-c214643cd70b",
				"SCW_DEFAULT_REGION":          "fr-par",
				"SCW_DEFAULT_ZONE":            "fr-par-1",
				"SCW_DEFAULT_PROJECT_ID":      "01696acf-7a78-4d94-a129-5e135d0377cc",
				"SCW_DEFAULT_ORGANIZATION_ID": "14800390-5df1-4a90-b38e-9b461bdcd108",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "yolo",
						fieldname.SecretAccessKey: "c3a22663-2770-4428-8166-c214643cd70b",
						fieldname.DefaultRegion:   "fr-par",
						fieldname.DefaultZone:     "fr-par-1",
						fieldname.ProjectID:       "01696acf-7a78-4d94-a129-5e135d0377cc",
						fieldname.OrgID:           "14800390-5df1-4a90-b38e-9b461bdcd108",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.config/scw/config.yaml": plugintest.LoadFixture(t, "config"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "yolo",
						fieldname.SecretAccessKey: "c3a22663-2770-4428-8166-c214643cd70b",
						fieldname.DefaultRegion:   "fr-par",
						fieldname.DefaultZone:     "fr-par-1",
						fieldname.ProjectID:       "01696acf-7a78-4d94-a129-5e135d0377cc",
						fieldname.OrgID:           "14800390-5df1-4a90-b38e-9b461bdcd108",
					},
				},
			},
		},
	})
}
