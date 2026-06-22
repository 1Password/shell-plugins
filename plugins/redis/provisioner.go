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

// Redis CLI flags that, when already supplied by the user, signal that we
// should not provision the corresponding field from the 1Password item.
var (
	hostFlags     = []string{"-h"}
	portFlags     = []string{"-p"}
	userFlags     = []string{"--user"}
	passwordFlags = []string{"-a", "--pass"}
)

func (p redisArgsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	suppliedFlags := flagSet(out.CommandLine)

	// The password is passed via an environment variable so it never appears in
	// the process's argument list. Skip it if the user already authenticated on
	// the command line.
	if value, ok := in.ItemFields[fieldname.Password]; ok && !containsAny(suppliedFlags, passwordFlags) {
		out.AddEnvVar("REDISCLI_AUTH", value)
	}

	// Collect the flags to inject first, then prepend them in a single pass.
	// Mutating out.CommandLine while ranging over it risks index-out-of-range
	// panics and stale reads, so we never modify it during inspection.
	var injected []string
	if value, ok := in.ItemFields[fieldname.Host]; ok && !containsAny(suppliedFlags, hostFlags) {
		injected = append(injected, "-h", value)
	}
	if value, ok := in.ItemFields[fieldname.Port]; ok && !containsAny(suppliedFlags, portFlags) {
		injected = append(injected, "-p", value)
	}
	if value, ok := in.ItemFields[fieldname.Username]; ok && !containsAny(suppliedFlags, userFlags) {
		injected = append(injected, "--user", value)
	}

	if len(injected) > 0 {
		out.CommandLine = prependArgs(out.CommandLine, injected)
	}
}

func (p redisArgsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: credentials get wiped automatically when the process exits.
}

func (p redisArgsProvisioner) Description() string {
	return "Provision redis secrets as command-line arguments and the password as an environment variable."
}

// flagSet returns the set of arguments present on the command line, excluding
// the executable name at index 0.
func flagSet(commandLine []string) map[string]bool {
	set := make(map[string]bool, len(commandLine))
	for i, arg := range commandLine {
		if i == 0 {
			continue
		}
		set[arg] = true
	}
	return set
}

func containsAny(set map[string]bool, flags []string) bool {
	for _, f := range flags {
		if set[f] {
			return true
		}
	}
	return false
}

// prependArgs inserts args immediately after the executable name (index 0),
// leaving the rest of the user-supplied command line intact.
func prependArgs(commandLine []string, args []string) []string {
	if len(commandLine) == 0 {
		return append([]string{}, args...)
	}
	result := make([]string, 0, len(commandLine)+len(args))
	result = append(result, commandLine[0])
	result = append(result, args...)
	result = append(result, commandLine[1:]...)
	return result
}
