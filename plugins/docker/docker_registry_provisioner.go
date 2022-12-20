package docker

import (
	"encoding/json"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/provision"
)

// ConfigFileProvisioner returns a provisioner for Docker that creates a `config.json` tempfile, based on
// the output of the passed in `configOpts` and sets the output path of the config file as the `DOCKER_CONFIG`
// environment variable which the Docker CLI will use to authenticate to a registry.
func ConfigFileProvisioner(configOpts ...DockerConfigOpt) sdk.Provisioner {
	fileFunc := func(in sdk.ProvisionInput) ([]byte, error) {
		config := &DockerConfig{
			Auths: make(map[string]DockerConfigAuth),
		}

		for _, opt := range configOpts {
			err := opt(config, in)
			if err != nil {
				return nil, err
			}
		}

		return json.Marshal(config)
	}

	return provision.TempFile(fileFunc, provision.SetPathAsEnvVar("DOCKER_CONFIG"))
}

type DockerConfigOpt func(config *DockerConfig, in sdk.ProvisionInput) error

func WithUsernameField(field sdk.FieldName) DockerConfigOpt {
	return func(config *DockerConfig, in sdk.ProvisionInput) error {
		if username, ok := in.ItemFields[field]; ok {
			for _, auth := range config.Auths {
				auth.Username = username
			}
		}

		return nil
	}
}

func WithStaticUsername(usernameValue string) DockerConfigOpt {
	return func(config *DockerConfig, in sdk.ProvisionInput) error {
		for _, auth := range config.Auths {
			auth.Username = usernameValue
		}

		return nil
	}
}

func WithPasswordField(field sdk.FieldName) DockerConfigOpt {
	return func(config *DockerConfig, in sdk.ProvisionInput) error {
		if password, ok := in.ItemFields[field]; ok {
			for _, auth := range config.Auths {
				auth.Password = password
			}
		}

		return nil
	}
}

func WithPasswordFunc(passwordFunc func(in sdk.ProvisionInput) (string, error)) DockerConfigOpt {
	return func(config *DockerConfig, in sdk.ProvisionInput) error {
		password, err := passwordFunc(in)
		if err != nil {
			return err
		}

		for _, auth := range config.Auths {
			auth.Password = password
		}

		return nil
	}
}

func WithStaticURL(url string) DockerConfigOpt {
	return func(config *DockerConfig, in sdk.ProvisionInput) error {
		config.Auths[url] = DockerConfigAuth{}
		return nil
	}
}

func WithURLField(field sdk.FieldName) DockerConfigOpt {
	return func(config *DockerConfig, in sdk.ProvisionInput) error {
		if url, ok := in.ItemFields[field]; ok {
			config.Auths[url] = DockerConfigAuth{}
		}
		return nil
	}
}

type DockerConfig struct {
	Auths       map[string]DockerConfigAuth `json:"auths"`
	CredHelpers map[string]string           `json:"credHelpers"`
}

type DockerConfigAuth struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Auth     string `json:"auth,omitempty"`
}
