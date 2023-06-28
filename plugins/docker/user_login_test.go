// package docker

// import (
// 	"testing"

// 	"github.com/1Password/shell-plugins/sdk"
// 	"github.com/1Password/shell-plugins/sdk/plugintest"
// 	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
// )

// func TestUserLoginProvisioner(t *testing.T) {
// 	plugintest.TestProvisioner(t, UserLogin().DefaultProvisioner, map[string]plugintest.ProvisionCase{
// 		"default": {
// 			ItemFields: map[sdk.FieldName]string{
// 				fieldname.Username: "4LjPnJ2u4Yo02KRfP7ffF1Tf2eoDZxnvNEXAMPLE",
// 				fieldname.Password: "12345",
// 			},
// 			ExpectedOutput: sdk.ProvisionOutput{
// 				Environment: map[string]string{
// 					"DOCKER_USERNAME": "4LjPnJ2u4Yo02KRfP7ffF1Tf2eoDZxnvNEXAMPLE",
// 					"DOCKER_PASSWORD": "12345",
// 				},
// 			},
// 		},
// 	})
// }

// func TestUserLoginImporter(t *testing.T) {
// 	plugintest.TestImporter(t, UserLogin().Importer, map[string]plugintest.ImportCase{
// 		"environment": {
// 			Environment: map[string]string{
// 				"DOCKER_USERNAME": "4LjPnJ2u4Yo02KRfP7ffF1Tf2eoDZxnvNEXAMPLE",
// 				"DOCKER_PASSWORD": "12345",
// 			},
// 			ExpectedCandidates: []sdk.ImportCandidate{
// 				{
// 					Fields: map[sdk.FieldName]string{
// 						fieldname.Username: "4LjPnJ2u4Yo02KRfP7ffF1Tf2eoDZxnvNEXAMPLE",
// 						fieldname.Password: "12345",
// 					},
// 				},
// 			},
// 		},

// 		"config file": {
// 			Files: map[string]string{
// 				"~/.docker/config.json": plugintest.LoadFixture(t, "config.json"),
// 			},
// 			ExpectedCandidates: []sdk.ImportCandidate{
// 				{
// 					Fields: map[sdk.FieldName]string{
// 						fieldname.Username: "4LjPnJ2u4Yo02KRfP7ffF1Tf2eoDZxnvNEXAMPLE",
// 						fieldname.Password: "12345",
// 					},
// 				},
// 			},
// 		},
// 	})
// }
