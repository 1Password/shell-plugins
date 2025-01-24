package motherduck

import (
	"testing"
	
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)
	
func TestAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.Token: "TERAkHVPg65C6UGDw42llLlgPtZhBbafxnpqs74fjyuKnDpSt7TZNODw3catnivaruR09REDcNIwystkLMlRw5foxRjvytBFmkk0t0x9iHqY0MBY40Ltbcdw8fvt3OzsCgxmbh89v0XIWrRiwCfALA1dbqWDLaatAZWOLQhJmYcggQR6YBVoKM9H7XBrBjDtP7YJOoU2Z7rc7KWgTTqS9vyCtLx7GDSBitWQLvUYuWzvgh94qk1Wt16oua34jzDtosd59ahNlvA1vEqPtkYqC5mNbDbWqcunwelka4tI4uuEfyojeXowBzkv6izjT48J3usTPIIqTFMYgJnMwUtV6n8UgeuLumEsKd86HVLywapqO37zfNrlrVLzjHSv0rGA2NjDgBAueK2clqEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"MOTHERDUCK_TOKEN": "TERAkHVPg65C6UGDw42llLlgPtZhBbafxnpqs74fjyuKnDpSt7TZNODw3catnivaruR09REDcNIwystkLMlRw5foxRjvytBFmkk0t0x9iHqY0MBY40Ltbcdw8fvt3OzsCgxmbh89v0XIWrRiwCfALA1dbqWDLaatAZWOLQhJmYcggQR6YBVoKM9H7XBrBjDtP7YJOoU2Z7rc7KWgTTqS9vyCtLx7GDSBitWQLvUYuWzvgh94qk1Wt16oua34jzDtosd59ahNlvA1vEqPtkYqC5mNbDbWqcunwelka4tI4uuEfyojeXowBzkv6izjT48J3usTPIIqTFMYgJnMwUtV6n8UgeuLumEsKd86HVLywapqO37zfNrlrVLzjHSv0rGA2NjDgBAueK2clqEXAMPLE",
				},
			},
		},
	})
}

func TestAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"MOTHERDUCK_TOKEN": "TERAkHVPg65C6UGDw42llLlgPtZhBbafxnpqs74fjyuKnDpSt7TZNODw3catnivaruR09REDcNIwystkLMlRw5foxRjvytBFmkk0t0x9iHqY0MBY40Ltbcdw8fvt3OzsCgxmbh89v0XIWrRiwCfALA1dbqWDLaatAZWOLQhJmYcggQR6YBVoKM9H7XBrBjDtP7YJOoU2Z7rc7KWgTTqS9vyCtLx7GDSBitWQLvUYuWzvgh94qk1Wt16oua34jzDtosd59ahNlvA1vEqPtkYqC5mNbDbWqcunwelka4tI4uuEfyojeXowBzkv6izjT48J3usTPIIqTFMYgJnMwUtV6n8UgeuLumEsKd86HVLywapqO37zfNrlrVLzjHSv0rGA2NjDgBAueK2clqEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "TERAkHVPg65C6UGDw42llLlgPtZhBbafxnpqs74fjyuKnDpSt7TZNODw3catnivaruR09REDcNIwystkLMlRw5foxRjvytBFmkk0t0x9iHqY0MBY40Ltbcdw8fvt3OzsCgxmbh89v0XIWrRiwCfALA1dbqWDLaatAZWOLQhJmYcggQR6YBVoKM9H7XBrBjDtP7YJOoU2Z7rc7KWgTTqS9vyCtLx7GDSBitWQLvUYuWzvgh94qk1Wt16oua34jzDtosd59ahNlvA1vEqPtkYqC5mNbDbWqcunwelka4tI4uuEfyojeXowBzkv6izjT48J3usTPIIqTFMYgJnMwUtV6n8UgeuLumEsKd86HVLywapqO37zfNrlrVLzjHSv0rGA2NjDgBAueK2clqEXAMPLE",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in motherduck/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
			// 	{
			// 		Fields: map[sdk.FieldName]string{
			// 			fieldname.Token: "TERAkHVPg65C6UGDw42llLlgPtZhBbafxnpqs74fjyuKnDpSt7TZNODw3catnivaruR09REDcNIwystkLMlRw5foxRjvytBFmkk0t0x9iHqY0MBY40Ltbcdw8fvt3OzsCgxmbh89v0XIWrRiwCfALA1dbqWDLaatAZWOLQhJmYcggQR6YBVoKM9H7XBrBjDtP7YJOoU2Z7rc7KWgTTqS9vyCtLx7GDSBitWQLvUYuWzvgh94qk1Wt16oua34jzDtosd59ahNlvA1vEqPtkYqC5mNbDbWqcunwelka4tI4uuEfyojeXowBzkv6izjT48J3usTPIIqTFMYgJnMwUtV6n8UgeuLumEsKd86HVLywapqO37zfNrlrVLzjHSv0rGA2NjDgBAueK2clqEXAMPLE",
			// 		},
			// 	},
			},
		},
	})
}
