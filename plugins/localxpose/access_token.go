package localxpose

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"

	"github.com/99designs/keyring"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

const (
	serviceName    = "https://loclx.io"
	keyringKeyName = "token"
	tokenFilePath  = "~/.localxpose/.access"
)

func AccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AccessToken,
		DocsURL:       sdk.URL("https://localxpose.io/docs"),
		ManagementURL: sdk.URL("https://localxpose.io/dashboard/access"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.AccessToken,
				MarkdownDescription: "Token used to authenticate to LocalXpose.",
				Secret:              true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryOSKeyring(),
			TryAccessTokenFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"LX_ACCESS_TOKEN": fieldname.AccessToken,
}

func TryAccessTokenFile() sdk.Importer {
	return importer.TryFile(tokenFilePath, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		if contents.ToString() == "" {
			return
		}
		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.AccessToken: contents.ToString(),
			},
		})
	})
}

func TryOSKeyring() sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		availableBackends := keyring.AvailableBackends()
		if len(availableBackends) == 0 {
			return
		}
		for _, backendType := range availableBackends {
			attempt := out.NewAttempt(importer.SourceOther(string(backendType), ""))
			openKeyring, err := keyring.Open(keyring.Config{
				AllowedBackends:          []keyring.BackendType{backendType},
				KeychainTrustApplication: true,
				ServiceName:              serviceName,
			})
			if err != nil {
				attempt.AddError(err)
				return
			}
			key, err := openKeyring.Get(keyringKeyName)
			if err != nil {
				attempt.AddError(err)
				continue
			}

			data, err := decode(key.Data)
			if err != nil {
				attempt.AddError(err)
				data = key.Data
			}

			attempt.AddCandidate(sdk.ImportCandidate{
				Fields: map[sdk.FieldName]string{
					fieldname.AccessToken: string(data),
				},
			})
		}
	}
}

// decode handles keys stored with previous versions of github.com/zalando/go-keyring
func decode(keyData []byte) ([]byte, error) {
	const (
		prefixBase64Encoded = "go-keyring-base64:"
		prefixHexEncoded    = "go-keyring-encoded:"
	)
	switch {
	case bytes.HasPrefix(keyData, []byte(prefixBase64Encoded)):
		return base64.StdEncoding.DecodeString(string(keyData[len(prefixBase64Encoded):]))
	case bytes.HasPrefix(keyData, []byte(prefixHexEncoded)):
		return hex.DecodeString(string(keyData[len(prefixHexEncoded):]))
	default:
		return keyData, nil
	}
}
