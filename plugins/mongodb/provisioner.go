package mongodb

import (
	"context"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

type mongodbShellArgsProvisioner struct {
}

func mongodbShellProvisioner() sdk.Provisioner {
	return mongodbShellArgsProvisioner{}
}

func (p mongodbShellArgsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	listOfPossibleUserInputArguments := []string{
		"--host",
		"--port",
		"-u", "--username",
		"-p", "--password",
	}

	var (
		commandLineContainsHostArgument     bool = false
		commandLineContainsPortArgument     bool = false
		commandLineContainsUsernameArgument bool = false
		commandLineContainsPasswordArgument bool = false
	)

	connectionStringProvisioned := false

	for _, arg := range out.CommandLine {
		if strings.HasPrefix(arg, "mongodb://") || strings.HasPrefix(arg, "mongodb+srv://") {
			connectionStringProvisioned = true
			break
		}
	}

	if value, ok := in.ItemFields[fieldname.ConnectionString]; ok && !connectionStringProvisioned {
		out.AddArgs(value)
		connectionStringProvisioned = true
	}

	for i, arg := range out.CommandLine {
		for _, userArg := range listOfPossibleUserInputArguments {
			if arg == userArg {
				// Get the executable "mongosh", the matched argument we are looking for and its value
				commandLine := []string{out.CommandLine[0], out.CommandLine[i], out.CommandLine[i+1]}
				// Remove the matched argument and its value as they will be added to the beginning of the command line
				out.CommandLine = append(out.CommandLine[:i], out.CommandLine[i+2:]...)
				// Add the executable "mongosh", the matched argument and its value to the beginning of the command line
				commandLine = append(commandLine, out.CommandLine[1:]...)
				out.CommandLine = commandLine
				// Controller to check if the user has already provided the argument
				switch userArg {
				case "--host":
					commandLineContainsHostArgument = true
				case "--port":
					commandLineContainsPortArgument = true
				case "-u", "--username":
					commandLineContainsUsernameArgument = true
				case "-p", "--password":
					commandLineContainsPasswordArgument = true
				default:
					break
				}
			}
		}
	}

	if value, ok := in.ItemFields[fieldname.Host]; ok && !commandLineContainsHostArgument && !connectionStringProvisioned {
		commandLine := []string{out.CommandLine[0], "--host", value}
		commandLine = append(commandLine, out.CommandLine[1:]...)
		out.CommandLine = commandLine
	}

	if value, ok := in.ItemFields[fieldname.Port]; ok && !commandLineContainsPortArgument && !connectionStringProvisioned {
		commandLine := []string{out.CommandLine[0], "--port", value}
		commandLine = append(commandLine, out.CommandLine[1:]...)
		out.CommandLine = commandLine
	}

	if value, ok := in.ItemFields[fieldname.Username]; ok && !commandLineContainsUsernameArgument {
		commandLine := []string{out.CommandLine[0], "--username", value}
		commandLine = append(commandLine, out.CommandLine[1:]...)
		out.CommandLine = commandLine
	}

	if value, ok := in.ItemFields[fieldname.Password]; ok && !commandLineContainsPasswordArgument {
		commandLine := []string{out.CommandLine[0], "--password", value}
		commandLine = append(commandLine, out.CommandLine[1:]...)
		out.CommandLine = commandLine
	}
}

func (p mongodbShellArgsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: credentials get wiped automatically when the process exits.
}

func (p mongodbShellArgsProvisioner) Description() string {
	return "Provision MongoDB Shell secrets as command-line arguments."
}
