package server

import (
	"context"
	"net"
	"net/rpc"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/rpc/proto"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAccessKey  = sdk.CredentialName("Access Key")
	testSSOProfile = sdk.CredentialName("SSO Profile")
)

// markingProvisioner sets an environment variable named after itself, so that a test can tell which
// provisioner handled a call.
type markingProvisioner struct {
	name string
}

func (p markingProvisioner) Description() string { return p.name }

func (p markingProvisioner) Provision(_ context.Context, _ sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	out.AddEnvVar("PROVISIONED_BY", p.name)
}

func (p markingProvisioner) Deprovision(context.Context, sdk.DeprovisionInput, *sdk.DeprovisionOutput) {
}

// markingImporter returns an importer that reports a single candidate named after the given name, so
// that a test can tell which importer handled a call.
func markingImporter(name string) sdk.Importer {
	return func(_ context.Context, _ sdk.ImportInput, out *sdk.ImportOutput) {
		out.Attempts = append(out.Attempts, &sdk.ImportAttempt{
			Candidates: []sdk.ImportCandidate{{NameHint: name}},
		})
	}
}

// testPluginWithTwoCredentialTypes returns a plugin with two credential types that are both used by
// its only executable. The first usage relies on the default provisioner of its credential type, the
// second usage brings its own.
func testPluginWithTwoCredentialTypes() schema.Plugin {
	return schema.Plugin{
		Name:     "example",
		Platform: schema.PlatformInfo{Name: "Example"},
		Credentials: []schema.CredentialType{
			{
				Name:               testAccessKey,
				DefaultProvisioner: markingProvisioner{name: "access key default"},
				Importer:           markingImporter("access key importer"),
			},
			{
				Name:               testSSOProfile,
				DefaultProvisioner: markingProvisioner{name: "sso profile default"},
				Importer:           markingImporter("sso profile importer"),
			},
		},
		Executables: []schema.Executable{
			{
				Name: "Example CLI",
				Runs: []string{"example"},
				Uses: []schema.CredentialUsage{
					{Name: testAccessKey, Optional: true},
					{Name: testSSOProfile, Optional: true, Provisioner: markingProvisioner{name: "sso profile usage"}},
				},
			},
		},
	}
}

// connectToTestPlugin serves the given plugin over RPC, in the same way as go-plugin does, and returns
// a client that is connected to it.
func connectToTestPlugin(t *testing.T, plugin schema.Plugin) *rpc.Client {
	t.Helper()

	rpcServer := rpc.NewServer()
	require.NoError(t, rpcServer.RegisterName("Plugin", newServer(plugin)))

	serverConn, clientConn := net.Pipe()
	go rpcServer.ServeConn(serverConn)

	client := rpc.NewClient(clientConn)
	t.Cleanup(func() {
		_ = client.Close()
	})
	return client
}

func newTestProvisionRequest(id proto.ProvisionerID) proto.ProvisionCredentialRequest {
	return proto.ProvisionCredentialRequest{
		ProvisionerID: id,
		ProvisionOutput: sdk.ProvisionOutput{
			Environment: map[string]string{},
			CommandLine: []string{"example"},
			Files:       map[string]sdk.OutputFile{},
		},
	}
}

func TestGetPluginWithTwoCredentialTypesOverRPC(t *testing.T) {
	client := connectToTestPlugin(t, testPluginWithTwoCredentialTypes())

	var resp proto.GetPluginResponse
	require.NoError(t, client.Call("Plugin.GetPlugin", 0, &resp))

	require.Len(t, resp.Credentials, 2)
	assert.Equal(t, testAccessKey, resp.Credentials[0].Name)
	assert.Equal(t, testSSOProfile, resp.Credentials[1].Name)
	assert.Equal(t, map[proto.CredentialID]bool{0: true, 1: true}, resp.CredentialHasImporter)

	require.Len(t, resp.Executables, 1)
	uses := resp.Executables[0].Uses
	require.Len(t, uses, 2)
	assert.Equal(t, testAccessKey, uses[0].Name)
	assert.True(t, uses[0].Optional)
	assert.Equal(t, testSSOProfile, uses[1].Name)
	assert.True(t, uses[1].Optional)
	assert.Equal(t, map[proto.CredentialUsageID]bool{
		{Executable: 0, Usage: 0}: false,
		{Executable: 0, Usage: 1}: true,
	}, resp.CredentialUsageHasProvisioner)
}

func TestProvisionWithTwoCredentialTypesOverRPC(t *testing.T) {
	client := connectToTestPlugin(t, testPluginWithTwoCredentialTypes())

	cases := map[string]struct {
		id          proto.ProvisionerID
		expected    string
		expectedErr string
	}{
		"default provisioner of the first credential type": {
			id:       proto.ProvisionerID{IsDefaultProvisioner: true, Credential: 0},
			expected: "access key default",
		},
		"default provisioner of the second credential type": {
			id:       proto.ProvisionerID{IsDefaultProvisioner: true, Credential: 1},
			expected: "sso profile default",
		},
		"provisioner of the second credential usage": {
			id:       proto.ProvisionerID{CredentialUsage: proto.CredentialUsageID{Executable: 0, Usage: 1}},
			expected: "sso profile usage",
		},
		"first credential usage has no provisioner of its own": {
			id:          proto.ProvisionerID{CredentialUsage: proto.CredentialUsageID{Executable: 0, Usage: 0}},
			expectedErr: "field not set for plugin.Executables[0].Uses[0].Provisioner.Provisioner",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var description string
			err := client.Call("Plugin.CredentialProvisionerDescription", tc.id, &description)
			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, description)

			var out sdk.ProvisionOutput
			require.NoError(t, client.Call("Plugin.CredentialProvisionerProvision", newTestProvisionRequest(tc.id), &out))
			assert.Equal(t, map[string]string{"PROVISIONED_BY": tc.expected}, out.Environment)
			assert.Equal(t, []string{"example"}, out.CommandLine)
		})
	}
}

func TestImportWithTwoCredentialTypesOverRPC(t *testing.T) {
	client := connectToTestPlugin(t, testPluginWithTwoCredentialTypes())

	cases := map[string]struct {
		id       proto.CredentialID
		expected string
	}{
		"importer of the first credential type": {
			id:       0,
			expected: "access key importer",
		},
		"importer of the second credential type": {
			id:       1,
			expected: "sso profile importer",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var out sdk.ImportOutput
			require.NoError(t, client.Call("Plugin.CredentialImport", proto.ImportCredentialRequest{CredentialID: tc.id}, &out))

			require.Len(t, out.Attempts, 1)
			require.Len(t, out.Attempts[0].Candidates, 1)
			assert.Equal(t, tc.expected, out.Attempts[0].Candidates[0].NameHint)
		})
	}
}
