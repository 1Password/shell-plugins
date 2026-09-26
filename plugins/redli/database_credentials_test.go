package redli

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
)

func TestDatabaseCredentialsProvisioner(t *testing.T) {
	const (
		insecureProtocol = "redis"
		secureProtocol   = "rediss"
		host             = "Host"
		port             = "Port"
		username         = "Username"
		password         = "Password"
		db               = "Db"
		uri              = username + ":" + password + "@" + host + ":" + port
	)
	plugintest.TestProvisioner(t, DatabaseCredentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"using uri": {
			ItemFields: map[sdk.FieldName]string{
				fnUri: uri,
			},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"-u", uri},
			},
		},
		"using args (secure)": {
			ItemFields: map[sdk.FieldName]string{
				fnProtocol: secureProtocol,
				fnHost:     host,
				fnPort:     port,
				fnUsername: username,
				fnPassword: password,
				fnDb:       db,
			},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{
					"-h", host,
					"-p", port,
					"-r", username,
					"-a", password,
					"--tls",
					"-n", db,
				},
			},
		},
		"using args (not secure)": {
			ItemFields: map[sdk.FieldName]string{
				fnProtocol: insecureProtocol,
				fnHost:     host,
				fnPort:     port,
				fnUsername: username,
				fnPassword: password,
				fnDb:       db,
			},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{
					"-h", host,
					"-p", port,
					"-r", username,
					"-a", password,
					"-n", db,
				},
			},
		},
	})
}
