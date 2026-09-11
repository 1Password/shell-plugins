package redli

import (
	"context"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

const (
	fnUri      = sdk.FieldName("Uri")
	fnProtocol = sdk.FieldName("Protocol")
	fnHost     = sdk.FieldName("Server")
	fnPort     = sdk.FieldName("Port")
	fnUsername = sdk.FieldName("Username")
	fnPassword = sdk.FieldName("Password")
	fnDb       = sdk.FieldName("Database")
)

func DatabaseCredentials() schema.CredentialType {
	return schema.CredentialType{
		Name: credname.DatabaseCredentials,
		Fields: []schema.CredentialField{
			{
				Name:                fnUri,
				MarkdownDescription: "URI to connect to host",
				Secret:              true,
				Optional:            true,
			},
			{
				Name:                fnProtocol,
				MarkdownDescription: "Protocol to use when connecting",
				Secret:              false,
				Optional:            true,
			},
			{
				Name:                fnHost,
				MarkdownDescription: "Host to connect to",
				Secret:              false,
				Optional:            true,
			},
			{
				Name:                fnPort,
				MarkdownDescription: "Port to connect to",
				Secret:              false,
				Optional:            true,
			},
			{
				Name:                fnDb,
				MarkdownDescription: "DB to connect to",
				Secret:              false,
				Optional:            true,
			},
			{
				Name:                fnUsername,
				MarkdownDescription: "Username to use when connecting to host",
				Secret:              true,
				Optional:            true,
			},
			{
				Name:                fnPassword,
				MarkdownDescription: "Password to use when connecting to host",
				Secret:              true,
				Optional:            true,
			},
		},
		DefaultProvisioner: commandLineProvisioner{},
	}
}

type commandLineProvisioner struct{}

func (p commandLineProvisioner) Description() string {
	return ""
}

func (p commandLineProvisioner) Provision(ctx context.Context, input sdk.ProvisionInput, output *sdk.ProvisionOutput) {
	foundArgs := make(map[string]bool)
	for _, arg := range output.CommandLine {
		if strings.HasPrefix(arg, "-") {
			foundArgs[arg] = true
		}
	}

	if foundArgs["-u"] {
		return
	}

	if uri, ok := input.ItemFields[fnUri]; ok {
		output.AddArgs("-u", uri)
		return
	}

	var (
		protocol = input.ItemFields[fnProtocol]
		host     = input.ItemFields[fnHost]
		port     = input.ItemFields[fnPort]
		username = input.ItemFields[fnUsername]
		password = input.ItemFields[fnPassword]
		db       = input.ItemFields[fnDb]
	)

	if _, ok := foundArgs["-h"]; !ok {
		if _, ok := foundArgs["--host"]; !ok {
			output.AddArgs("-h", host)
		}
	}

	if _, ok := foundArgs["-p"]; !ok {
		if _, ok := foundArgs["--port"]; !ok {
			output.AddArgs("-p", port)
		}
	}

	if _, ok := foundArgs["-r"]; !ok {
		if _, ok := foundArgs["--redisuser"]; !ok {
			output.AddArgs("-r", username)
		}
	}

	if _, ok := foundArgs["-a"]; !ok {
		if _, ok := foundArgs["--auth"]; !ok {
			output.AddArgs("-a", password)
		}
	}

	if _, ok := foundArgs["--tls"]; !ok {
		if protocol == "rediss" {
			output.AddArgs("--tls")
		}
	}

	if _, ok := foundArgs["-n"]; !ok {
		output.AddArgs("-n", db)
	}
}

func (p commandLineProvisioner) Deprovision(ctx context.Context, input sdk.DeprovisionInput, output *sdk.DeprovisionOutput) {
}
