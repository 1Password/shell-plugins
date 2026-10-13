package aws

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The first two profiles are the ~/.aws/config from issue #559.
const providerChainConfig = `[profile 1pass-source-profile]
region = eu-west-1

[profile my-1password-profile]
source_profile = 1pass-source-profile
role_arn       = arn:aws:iam::222222222222:role/MyRole
region         = eu-west-1

[profile direct-role]
role_arn       = arn:aws:iam::222222222222:role/MyRole
region         = eu-west-1

[profile chained-role]
source_profile = my-1password-profile
role_arn       = arn:aws:iam::333333333333:role/ChainedRole
region         = eu-west-1

[profile role-with-local-mfa]
source_profile = 1pass-source-profile
role_arn       = arn:aws:iam::222222222222:role/MyRole
mfa_serial     = arn:aws:iam::111111111111:mfa/some.user
region         = eu-west-1
`

// stsCallChain lists the STS calls a provider makes, outermost first, ending
// with what signs the innermost call.
func stsCallChain(p aws.CredentialsProvider) []string {
	var chain []string
	for p != nil {
		switch v := p.(type) {
		case *assumeRoleProvider:
			chain = append(chain, "AssumeRole")
			p = v.StsClient.Options().Credentials
		case *mfaSessionTokenProvider:
			hop := "GetSessionToken"
			if v.Mfa != nil && v.GetMfaSerial() != "" {
				hop += "(MFA)"
			}
			chain = append(chain, hop)
			p = v.StsClient.Options().Credentials
		case accessKeysProvider:
			chain = append(chain, "long-lived-keys")
			p = nil
		default:
			chain = append(chain, fmt.Sprintf("%T", p))
			p = nil
		}
	}
	return chain
}

func TestSTSProviderChain(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "awsConfig")
	require.NoError(t, os.WriteFile(configPath, []byte(providerChainConfig), 0600))
	t.Setenv("AWS_CONFIG_FILE", configPath)
	t.Setenv("AWS_PROFILE", "")
	t.Setenv("AWS_DEFAULT_REGION", "")

	keysAndOTP := map[sdk.FieldName]string{
		fieldname.AccessKeyID:     "AKIAHPIZFMD5EEXAMPLE",
		fieldname.SecretAccessKey: "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
		fieldname.OneTimePassword: "123456",
	}
	withItemSerial := map[sdk.FieldName]string{fieldname.MFASerial: "arn:aws:iam::111111111111:mfa/some.user"}
	for k, v := range keysAndOTP {
		withItemSerial[k] = v
	}

	for _, c := range []struct {
		name, profile string
		itemFields    map[sdk.FieldName]string
		want          []string
	}{
		{"source_profile with role", "my-1password-profile", withItemSerial, []string{"AssumeRole", "GetSessionToken(MFA)", "long-lived-keys"}},
		{"role without source_profile", "direct-role", withItemSerial, []string{"AssumeRole", "GetSessionToken(MFA)", "long-lived-keys"}},
		{"role chained through role", "chained-role", withItemSerial, []string{"AssumeRole", "AssumeRole", "GetSessionToken(MFA)", "long-lived-keys"}},
		// mfa_serial set only on the role profile must still reach the chain root.
		{"mfa_serial only on role profile", "role-with-local-mfa", keysAndOTP, []string{"AssumeRole", "GetSessionToken(MFA)", "long-lived-keys"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			awsConfig, err := getAWSAuthConfigurationForProfile(c.profile)
			require.NoError(t, err)

			// Cold cache: nothing to short-circuit the STS calls.
			factory := &CacheProviderFactory{
				InCache:    sdk.CacheState{},
				OutCache:   sdk.CacheOperations{Puts: map[string]sdk.CacheEntry{}},
				ItemFields: c.itemFields,
			}
			provider, err := GetTemporaryCredentialsProviderForProfile(awsConfig, factory, c.itemFields)
			require.NoError(t, err)

			assert.Equal(t, c.want, stsCallChain(provider))
		})
	}
}
