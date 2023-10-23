package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func Credentials() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.Credentials,
		DocsURL:       sdk.URL("https://docs.docker.com/engine/reference/commandline/login"),
		ManagementURL: sdk.URL("https://hub.docker.com/settings/security"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Username,
				MarkdownDescription: "Username used in Docker registries.",
				Secret:              false,
			},
			{
				Name:                fieldname.Secret,
				AlternativeNames:    []string{fieldname.Password.String(), fieldname.AccessToken.String()},
				MarkdownDescription: "Password or access token used to authenticate to a Docker registry.",
				Secret:              true,
			},
			{
				Name:                fieldname.Host,
				MarkdownDescription: "URL of the Docker registry server.",
				Optional:            true, // Defaults to Docker Hub registry
				Secret:              false,
			},
		},
		DefaultProvisioner: dockerProvisioner{},
		Importer: importer.TryAll(
			TryDockerConfigFile(),
		)}
}

func TryDockerConfigFile() sdk.Importer {
	// TODO use programmatic file finding to get HomeDir/.../config.json
	return importer.TryFile("~/.docker/config.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config AuthConfig
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		// Credit: Fixed with code from:
		// https://github.com/1Password/shell-plugins/pull/301/files#diff-9d1e3341794f3990528aebd795097bfc281f6e2acf3ec92f5a0bb02e9740ba9e
		for url, auth := range config.Auths {
			authBytes, err := base64.StdEncoding.DecodeString(auth.Auth)
			if err != nil {
				out.AddError(err)
				return
			}
			hostUrl := url
			credentials := string(authBytes)
			username, secret := parseCredentials(credentials)
			if username != "" {
				out.AddCandidate(sdk.ImportCandidate{
					Fields: map[sdk.FieldName]string{
						fieldname.Username: username,
						fieldname.Secret:   secret,
						fieldname.Host:     hostUrl,
					},
				})
			}
		}
	})
}

type dockerProvisioner struct {
}

func (p dockerProvisioner) Description() string {
	return "Docker login credentials provisioner"
}

func (p dockerProvisioner) Provision(ctx context.Context, input sdk.ProvisionInput, output *sdk.ProvisionOutput) {

	configFilePath := filepath.Join(input.HomeDir, ".docker", "config.json")

	configFileContents, err := json.Marshal(map[string]string{"credsStore": "1password"})
	if err != nil {
		output.AddError(err)
		return
	}

	executablePath := filepath.Join(input.TempDir, "docker-credential-1password")

	err = os.WriteFile(executablePath, []byte(bashScript), 0777)
	if err != nil {
		output.AddError(err)
		return
	}

	//output.AddFile(executablePath, sdk.OutputFile{Contents: []byte(bashScript)})
	err = os.MkdirAll(filepath.Join(input.HomeDir, ".docker"), os.ModePerm)
	if err != nil {
		output.AddError(fmt.Errorf("Config dir creation: %w", err))
		return
	}

	//output.AddNonSecretFile(configFilePath, configFileContents)
	err = os.WriteFile(configFilePath, configFileContents, 0666)
	if err != nil {
		output.AddError(fmt.Errorf("Error writing config file: %w", err))
		return
	}

	if registry := input.ItemFields[fieldname.Host]; registry != "" {
		output.AddEnvVar("DOCKER_REGISTRY", registry)
	}
	if username := input.ItemFields[fieldname.Username]; username != "" {
		output.AddEnvVar("DOCKER_CREDS_USR", input.ItemFields[fieldname.Username])
	}
	if secret := input.ItemFields[fieldname.Secret]; secret != "" {
		output.AddEnvVar("DOCKER_CREDS_PSW", input.ItemFields[fieldname.Secret])
	}

	path := os.Getenv("PATH")
	output.AddEnvVar("PATH", path+":"+input.TempDir)
}

func (p dockerProvisioner) Deprovision(ctx context.Context, input sdk.DeprovisionInput, output *sdk.DeprovisionOutput) {
	// TODO remove files that were manually added
}

type AuthConfig struct {
	Auths map[string]struct {
		Auth string `json:"auth"`
	} `json:"auths"`
}

// Credit: copied from https://github.com/1Password/shell-plugins/pull/301/files#diff-9d1e3341794f3990528aebd795097bfc281f6e2acf3ec92f5a0bb02e9740ba9e
func parseCredentials(credentials string) (username string, password string) {
	parts := strings.SplitN(credentials, ":", 2)
	if len(parts) == 2 {
		username = parts[0]
		password = parts[1]
	}
	return username, password
}

// Modified from https://gist.github.com/jasonk/480d87b49e4c8caf51932f184ff764b2
const bashScript = `#!/bin/bash
# docker-credential-1password
# 2018 - Jason Kohles

REG="${DOCKER_REGISTRY#https://}"
REG="${REG%%/*}"

die() {
  echo "$@" 1>&2
  exit 1
}

if [ -z "$REG" ]; then die "DOCKER_REGISTRY not set in environment"; fi
case "$1" in
  get)
    read HOST
    HOST="${HOST#https://}"
    HOST="${HOST%%/*}"
    if [ "$HOST" = "$REG" ]; then
      printf '{"ServerURL":"%s","Username":"%q","Secret":"%q"}\n' \
        "$HOST" "$DOCKER_CREDS_USR" "$DOCKER_CREDS_PSW"
    else
      die "No credentials available for $HOST"
    fi  
    ;;  
  *)  
    die "Unsupported operation"
    ;;  
esac
`
