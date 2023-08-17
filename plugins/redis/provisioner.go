package redis

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

type redisArgsProvisioner struct {
}

func redisProvisioner() sdk.Provisioner {
	return redisArgsProvisioner{}
}

func (p redisArgsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	listOfPossibleUserInputArguments := []string{
		"--user",
		"-h",
		"-p",
		"--pass",
		"-a",
	}

	var (
		commandLineContainsUserArgument     bool = false
		commandLineContainsHostArgument     bool = false
		commandLineContainsPortArgument     bool = false
		commandLineContainsPasswordArgument bool = false
	)

	for i, arg := range out.CommandLine {
		for _, userArg := range listOfPossibleUserInputArguments {
			if arg == userArg {
				// Get the executable "redis-cli", the matched argument we are looking for and its value
				commandLine := []string{out.CommandLine[0], out.CommandLine[i], out.CommandLine[i+1]}
				// Remove the matched argument and its value as they will be added to the beginning of the command line
				out.CommandLine = append(out.CommandLine[:i], out.CommandLine[i+2:]...)
				// Add the executable "redis-cli", the matched argument and its value to the beginning of the command line
				commandLine = append(commandLine, out.CommandLine[1:]...)
				out.CommandLine = commandLine
				// Controller to check if the user has already provided the argument
				switch userArg {
				case "--user":
					commandLineContainsUserArgument = true
				case "-h":
					commandLineContainsHostArgument = true
				case "-p":
					commandLineContainsPortArgument = true
				case "--pass", "-a":
					commandLineContainsPasswordArgument = true
				default:
					break
				}
			}
		}
	}

	if value, ok := in.ItemFields[fieldname.Password]; ok && !commandLineContainsPasswordArgument {
		out.AddEnvVar("REDISCLI_AUTH", value)
	}

	if value, ok := in.ItemFields[fieldname.Host]; ok && !commandLineContainsHostArgument {
		commandLine := []string{out.CommandLine[0], "-h", value}
		commandLine = append(commandLine, out.CommandLine[1:]...)
		out.CommandLine = commandLine
	}

	if value, ok := in.ItemFields[fieldname.Username]; ok && !commandLineContainsUserArgument {
		commandLine := []string{out.CommandLine[0], "--user", value}
		commandLine = append(commandLine, out.CommandLine[1:]...)
		out.CommandLine = commandLine
	}

	if value, ok := in.ItemFields[fieldname.Port]; ok && !commandLineContainsPortArgument {
		commandLine := []string{out.CommandLine[0], "-p", value}
		commandLine = append(commandLine, out.CommandLine[1:]...)
		out.CommandLine = commandLine
	}
}

func (p redisArgsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: credentials get wiped automatically when the process exits.
}

func (p redisArgsProvisioner) Description() string {
	return "Provision redis secrets as command-line arguments."
}
