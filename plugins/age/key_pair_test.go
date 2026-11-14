package age

import (
	"testing"

	"github.com/1Password/shell-plugins/plugins/age/provisioner"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAsymmetricKeyPairProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, KeyPair().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"defaults-to-encryption-mode": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.PrivateKey: "AGE-SECRET-KEY-10000000000000000000000000000000000000000000000000000000000",
				fieldname.PublicKey:  "age10000000000000000000000000000000000000000000000000000000000",
			},
			CommandLine: []string{"age", "-o", "/tmp/encrypted.txt", "/tmp/unencrypted.txt"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"age", "-R", "/tmp/age.public.txt", "-o", "/tmp/encrypted.txt", "/tmp/unencrypted.txt"},
				Files: map[string]sdk.OutputFile{
					"/tmp/age.public.txt": {
						Contents: []byte(plugintest.LoadFixture(t, "age.public.txt")),
					},
				},
			},
		},
		"explicit-encryption-mode-short": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.PrivateKey: "AGE-SECRET-KEY-10000000000000000000000000000000000000000000000000000000000",
				fieldname.PublicKey:  "age10000000000000000000000000000000000000000000000000000000000",
			},
			CommandLine: []string{"age", "-e", "-o", "/tmp/encrypted.txt", "/tmp/unencrypted.txt"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"age", "-R", "/tmp/age.public.txt", "-e", "-o", "/tmp/encrypted.txt", "/tmp/unencrypted.txt"},
				Files: map[string]sdk.OutputFile{
					"/tmp/age.public.txt": {
						Contents: []byte(plugintest.LoadFixture(t, "age.public.txt")),
					},
				},
			},
		},
		"explicit-encryption-mode-long": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.PrivateKey: "AGE-SECRET-KEY-10000000000000000000000000000000000000000000000000000000000",
				fieldname.PublicKey:  "age10000000000000000000000000000000000000000000000000000000000",
			},
			CommandLine: []string{"age", "--encrypt", "-o", "/tmp/encrypted.txt", "/tmp/unencrypted.txt"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"age", "-R", "/tmp/age.public.txt", "--encrypt", "-o", "/tmp/encrypted.txt", "/tmp/unencrypted.txt"},
				Files: map[string]sdk.OutputFile{
					"/tmp/age.public.txt": {
						Contents: []byte(plugintest.LoadFixture(t, "age.public.txt")),
					},
				},
			},
		},
		"decryption-mode-short": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.PrivateKey: "AGE-SECRET-KEY-10000000000000000000000000000000000000000000000000000000000",
				fieldname.PublicKey:  "age10000000000000000000000000000000000000000000000000000000000",
			},
			CommandLine: []string{"age", "-d", "-o", "/tmp/encrypted.txt", "/tmp/unencrypted.txt"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"age", "-i", "/tmp/age.private.txt", "-d", "-o", "/tmp/encrypted.txt", "/tmp/unencrypted.txt"},
				Files: map[string]sdk.OutputFile{
					"/tmp/age.private.txt": {
						Contents: []byte(plugintest.LoadFixture(t, "age.private.txt")),
					},
				},
			},
		},
		"decryption-mode-long": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.PrivateKey: "AGE-SECRET-KEY-10000000000000000000000000000000000000000000000000000000000",
				fieldname.PublicKey:  "age10000000000000000000000000000000000000000000000000000000000",
			},
			CommandLine: []string{"age", "--decrypt", "-o", "/tmp/encrypted.txt", "/tmp/unencrypted.txt"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"age", "-i", "/tmp/age.private.txt", "--decrypt", "-o", "/tmp/encrypted.txt", "/tmp/unencrypted.txt"},
				Files: map[string]sdk.OutputFile{
					"/tmp/age.private.txt": {
						Contents: []byte(plugintest.LoadFixture(t, "age.private.txt")),
					},
				},
			},
		},
		"decryption-identity-flag-provided": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.PrivateKey: "AGE-SECRET-KEY-10000000000000000000000000000000000000000000000000000000000",
				fieldname.PublicKey:  "age10000000000000000000000000000000000000000000000000000000000",
			},
			CommandLine: []string{"age", "-i", "/tmp/age.user_provided_identity_flag.txt", "--decrypt", "-o", "/tmp/encrypted.txt", "/tmp/unencrypted.txt"},
			ExpectedOutput: sdk.ProvisionOutput{
				Diagnostics: sdk.Diagnostics{
					Errors: []sdk.Error{{Message: provisioner.ErrConflictingIdentityFlag.Error()}},
				},
				CommandLine: []string{"age", "-i", "/tmp/age.private.txt", "-i", "/tmp/age.user_provided_identity_flag.txt", "--decrypt", "-o", "/tmp/encrypted.txt", "/tmp/unencrypted.txt"},
				Files: map[string]sdk.OutputFile{
					"/tmp/age.private.txt": {
						Contents: []byte(plugintest.LoadFixture(t, "age.private.txt")),
					},
				},
			},
		},
	})
}
