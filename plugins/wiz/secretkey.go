// Package wiz is a 1Password Shell Plugin designed for handling provisioning
// and deprovisioning of credentials for the WizCLI utility used to interact
// with Wiz.
//
// The WizCLI's `auth` subcommand takes a `--id` (56 character alphanumeric
// string) and a `--secret` (64 character alphanumeric string) argument, or
// alternatively looks for the `WIZ_CLIENT_ID` and `WIZ_CLIENT_SECRET`
// environment variables, which it uses to request a time-limited access token
// (Bearer token) from `https://auth.app.wiz.io/oauth/token`. The access token
// is written to `$WIZ_DIR/auth.json` (defaults to `$HOME/wiz/auth.json`),
// along with an expiration date, a data center identifier, tenant id and the
// client id.
//
// When a subcommand of the wizcli, that requires authentication (such as
// `wizcli iac scan --path .` or `wizcli dir scan --path .`) is invoked, it
// looks for the `$WIZ_DIR/auth.json` file for the access token to include.
//
// This plugin implements the functionality of the auth subcommand and writes
// the `$WIZ_DIR/auth.json` file in the expected format. After the wizcli has
// completed, the `$WIZ_DIR/auth.json` file is deleted.
//
// The plugin does not support importing, as the client secret is never written
// to disk anywhere. Only access token (and a refresh token) is written, and
// they will expire fairly quickly.
//
// Note that the plugin doesn't invalidate the access token (or refresh token)
// created, because Wiz doesn't seem to expose an API for that.
//
// The plugin also requires the installation of the data center ID of ones Wiz
// tenant and the tenant ID. This information can be found under
// https://app.wiz.io/user/tenant. The information can alternatively be fetched
// via the Wiz API, but this uses GraphQL and would introduce a new dependency
// to this project.
//
// Finally, the authentication could also be done by having the plugin invoke
// the WizCLI's `auth` subcommand with the needed secret and client ID. Calling
// the authentication API was chosen as an alternative because calling the
// binary could be messy.
//
// A future iteration of this plugin could choose to use the $WIZ_DIR
// environment variable to point to a RAMdisk, to ensure the authentication
// files are never written to persistent storage.
package wiz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func SecretKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.SecretKey,
		DocsURL:       sdk.URL("https://docs.wiz.io/wiz-docs/docs/set-up-wiz-cli"),
		ManagementURL: sdk.URL("https://app.wiz.io/settings/service-accounts"),
		Fields: []schema.CredentialField{
			{
				Name:                "ClientSecret",
				MarkdownDescription: "Client secret used to authenticate to Wiz.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 64,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                "ClientID",
				MarkdownDescription: "Client ID of Wiz Service Account",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Length: 53,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                "DataCenter",
				MarkdownDescription: "The Wiz datacenter in use",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Length: 4,
					Charset: schema.Charset{
						Uppercase: false,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                "TenantID",
				MarkdownDescription: "The ID of the Wiz tenant",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Length: 36,
					Charset: schema.Charset{
						Uppercase: false,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: WizProvisioner{},
		Importer:           importer.NoOp()}
}

type WizAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type WizAuthFile struct {
	ClientID    string `json:"client_id"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   string `json:"expires_at"`
	DataCenter  string `json:"data_center"`
	TenantID    string `json:"tenant_id"`
}

type WizProvisioner struct{}

func (wp WizProvisioner) Description() string {
	return "The Wiz provisioner takes care of creating an authentication token from the client ID and client secret, and to the token from the local disk after the wiz cli exits"
}
func (wp WizProvisioner) Provision(ctx context.Context, input sdk.ProvisionInput, output *sdk.ProvisionOutput) {
	r, err := WizAuth(input.ItemFields["ClientID"], input.ItemFields["ClientSecret"])
	if err != nil {
		output.AddError(err)
		return
	}

	dur, err := time.ParseDuration(fmt.Sprintf("%ds", r.ExpiresIn))
	if err != nil {
		output.AddError(err)
		return
	}

	b, err := json.Marshal(WizAuthFile{
		ClientID:    input.ItemFields["ClientID"],
		AccessToken: r.AccessToken,
		TokenType:   r.TokenType,
		ExpiresAt:   time.Now().Add(dur).Format(time.RFC3339),
		DataCenter:  input.ItemFields["DataCenter"],
		TenantID:    input.ItemFields["TenantID"],
	})
	if err != nil {
		output.AddError(err)
		return
	}

	wizDir, ok := os.LookupEnv("WIZ_DIR")
	if !ok {
		wizDir = input.FromHomeDir("/.wiz")
	}

	output.AddSecretFile(path.Join(wizDir, "/auth.json"), b)
}
func (wp WizProvisioner) Deprovision(ctx context.Context, input sdk.DeprovisionInput, output *sdk.DeprovisionOutput) {
	// No deprovisioning steps, as credentials file is deleted as a result of it
	// being in output.Files in the provisioning step.
}

func WizAuth(clientID, clientSecret string) (WizAuthResponse, error) {
	params := url.Values{}
	params.Add("grant_type", `client_credentials`)
	params.Add("client_id", clientID)
	params.Add("client_secret", clientSecret)
	params.Add("audience", `wiz-api`)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://auth.app.wiz.io/oauth/token", body)
	if err != nil {
		return WizAuthResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return WizAuthResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return WizAuthResponse{}, fmt.Errorf("unable to retrieve token: %s", resp.Status)
	}

	var r WizAuthResponse
	dec := json.NewDecoder(resp.Body)
	for {
		if err := dec.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			return WizAuthResponse{}, err
		}
	}

	return r, nil
}
