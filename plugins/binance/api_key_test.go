package binance

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.APIKey:    "jThmEycY2J0RgJgNNrWQBq2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qEXAMPLE",
				fieldname.APISecret: "2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qjThmEycY2J0RgJgNNrWQBqEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"BINANCE_API_KEY":    "jThmEycY2J0RgJgNNrWQBq2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qEXAMPLE",
					"BINANCE_SECRET_KEY": "2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qjThmEycY2J0RgJgNNrWQBqEXAMPLE",
					"BINANCE_API_SECRET": "2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qjThmEycY2J0RgJgNNrWQBqEXAMPLE",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"binance-cli v2 environment": {
			Environment: map[string]string{
				"BINANCE_API_KEY":    "jThmEycY2J0RgJgNNrWQBq2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qEXAMPLE",
				"BINANCE_SECRET_KEY": "2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qjThmEycY2J0RgJgNNrWQBqEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey:    "jThmEycY2J0RgJgNNrWQBq2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qEXAMPLE",
						fieldname.APISecret: "2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qjThmEycY2J0RgJgNNrWQBqEXAMPLE",
					},
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "jThmEycY2J0RgJgNNrWQBq2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qEXAMPLE",
					},
				},
			},
		},
		"binance-cli v1 environment": {
			Environment: map[string]string{
				"BINANCE_API_KEY":    "jThmEycY2J0RgJgNNrWQBq2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qEXAMPLE",
				"BINANCE_API_SECRET": "2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qjThmEycY2J0RgJgNNrWQBqEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "jThmEycY2J0RgJgNNrWQBq2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qEXAMPLE",
					},
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey:    "jThmEycY2J0RgJgNNrWQBq2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qEXAMPLE",
						fieldname.APISecret: "2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qjThmEycY2J0RgJgNNrWQBqEXAMPLE",
					},
				},
			},
		},
		"both secret env vars with different values": {
			Environment: map[string]string{
				"BINANCE_API_KEY":    "jThmEycY2J0RgJgNNrWQBq2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qEXAMPLE",
				"BINANCE_SECRET_KEY": "2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qjThmEycY2J0RgJgNNrWQBqEXAMPLE",
				"BINANCE_API_SECRET": "Qx7Vb2NwLk9RtYp4HsMc6JdZf8GaUe3XiOn5TrWq1KyBv0PlCm2AzSj4DEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey:    "jThmEycY2J0RgJgNNrWQBq2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qEXAMPLE",
						fieldname.APISecret: "2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qjThmEycY2J0RgJgNNrWQBqEXAMPLE",
					},
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey:    "jThmEycY2J0RgJgNNrWQBq2raPzKvxCkcwPQFk8AuWUu5QxQSWaItIB1qEXAMPLE",
						fieldname.APISecret: "Qx7Vb2NwLk9RtYp4HsMc6JdZf8GaUe3XiOn5TrWq1KyBv0PlCm2AzSj4DEXAMPLE",
					},
				},
			},
		},
	})
}
