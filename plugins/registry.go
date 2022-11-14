package plugins

import (
	"fmt"
	"strings"

	"github.com/1Password/shell-plugins/plugins/aws"
	"github.com/1Password/shell-plugins/plugins/circleci"
	"github.com/1Password/shell-plugins/plugins/datadog"
	"github.com/1Password/shell-plugins/plugins/digitalocean"
	"github.com/1Password/shell-plugins/plugins/gitlab"
	"github.com/1Password/shell-plugins/plugins/heroku"
	"github.com/1Password/shell-plugins/plugins/okta"
	"github.com/1Password/shell-plugins/plugins/postgresql"
	"github.com/1Password/shell-plugins/plugins/sentry"
	"github.com/1Password/shell-plugins/plugins/snyk"
	"github.com/1Password/shell-plugins/plugins/stripe"
	"github.com/1Password/shell-plugins/plugins/twilio"
	"github.com/1Password/shell-plugins/plugins/vault"

	"github.com/1Password/shell-plugins/plugins/github"

	"github.com/1Password/shell-plugins/sdk/schema"
)

var registry = []schema.Plugin{}

func List() []schema.Plugin {
	list := make([]schema.Plugin, len(registry))
	copy(list, registry)
	return list
}

func Get(pluginName string) (schema.Plugin, error) {
	for _, p := range registry {
		if p.Name == pluginName {
			return p, nil
		}
	}
	return schema.Plugin{}, fmt.Errorf("unknown plugin: %s", pluginName)
}

func GetByExecutable(executableQuery string) (schema.Plugin, schema.Executable, error) {
	for _, p := range registry {
		for _, e := range p.Executables {
			if strings.EqualFold(executableQuery, e.Command()) || strings.EqualFold(executableQuery, e.Name) {
				return p, e, nil
			}
		}
	}
	return schema.Plugin{}, schema.Executable{}, fmt.Errorf("unknown plugin: %s", executableQuery)
}

func Register(p schema.Plugin) {
	registry = append(registry, p)
}

func init() {
	Register(aws.New())
	Register(circleci.New())
	Register(datadog.New())
	Register(digitalocean.New())
	Register(github.New())
	Register(gitlab.New())
	Register(heroku.New())
	Register(okta.New())
	Register(postgresql.New())
	Register(sentry.New())
	Register(snyk.New())
	Register(stripe.New())
	Register(twilio.New())
	Register(vault.New())
}
