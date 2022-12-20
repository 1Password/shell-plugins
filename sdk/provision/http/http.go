package http

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/1Password/shell-plugins/sdk"
)

const (
	HeaderNameAuthorization = "Authorization"
)

type ValueFunc func(ctx context.Context, in sdk.HTTPProvisionInput) (string, error)

func FieldValue(fieldName sdk.FieldName) ValueFunc {
	return func(ctx context.Context, in sdk.HTTPProvisionInput) (string, error) {
		return in.ItemFields[fieldName], nil
	}
}

func BearerToken(tokenFunc ValueFunc) sdk.HTTPProvisioner {
	return HTTPHeadersProvisioner(map[string]ValueFunc{
		HeaderNameAuthorization: func(ctx context.Context, in sdk.HTTPProvisionInput) (string, error) {
			token, err := tokenFunc(ctx, in)
			if err != nil {
				return "", err
			}

			return "Bearer " + token, nil
		},
	})
}

func BasicAuth(usernameFunc ValueFunc, passwordFunc ValueFunc) sdk.HTTPProvisioner {
	return HTTPHeadersProvisioner(map[string]ValueFunc{
		HeaderNameAuthorization: func(ctx context.Context, in sdk.HTTPProvisionInput) (string, error) {
			username, err := usernameFunc(ctx, in)
			if err != nil {
				return "", err
			}

			password, err := passwordFunc(ctx, in)
			if err != nil {
				return "", err
			}

			base64AuthStr := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
			return "Basic " + base64AuthStr, nil
		},
	})
}
