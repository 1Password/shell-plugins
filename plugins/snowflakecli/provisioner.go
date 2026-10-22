package snowflakecli

import (
	"context"
	"regexp"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

type snowTemporaryConnectionProvisioner struct{}

// snowProvisioner provisions the credential as SNOWFLAKE_* environment variables and adds -x (--temporary-connection)
// to the command line, so that snow builds the connection from those environment variables only.
//
// -x is needed because snow refuses to run without a configured connection, and connection values from the user's
// own config files take precedence over the generic SNOWFLAKE_* environment variables. A temporary config file can't
// be used instead: provisioned files can only be read once, while snow reads its config file multiple times.
func snowProvisioner() sdk.Provisioner {
	return snowTemporaryConnectionProvisioner{}
}

var envVarsByField = map[sdk.FieldName]string{
	fieldname.Account:         "SNOWFLAKE_ACCOUNT",
	fieldname.Username:        "SNOWFLAKE_USER",
	fieldPrivateKeyPassphrase: "PRIVATE_KEY_PASSPHRASE",
	fieldRole:                 "SNOWFLAKE_ROLE",
	fieldWarehouse:            "SNOWFLAKE_WAREHOUSE",
	fieldname.Database:        "SNOWFLAKE_DATABASE",
	fieldSchema:               "SNOWFLAKE_SCHEMA",
}

func (p snowTemporaryConnectionProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	for field, envVarName := range envVarsByField {
		if value, ok := in.ItemFields[field]; ok && value != "" {
			out.AddEnvVar(envVarName, value)
		}
	}
	out.AddEnvVar("SNOWFLAKE_AUTHENTICATOR", "SNOWFLAKE_JWT")
	out.AddEnvVar("SNOWFLAKE_PRIVATE_KEY_RAW", normalizePrivateKey(in.ItemFields[fieldname.PrivateKey]))

	out.CommandLine = addTemporaryConnectionFlag(out.CommandLine)
}

func (p snowTemporaryConnectionProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p snowTemporaryConnectionProvisioner) Description() string {
	return "Provision Snowflake CLI key pair credentials as environment variables for a temporary connection."
}

// addTemporaryConnectionFlag appends -x to the command line. If the user ended the options with "--", -x is inserted
// before it, so that it doesn't get treated as a positional argument.
func addTemporaryConnectionFlag(commandLine []string) []string {
	result := make([]string, 0, len(commandLine)+1)
	for i, arg := range commandLine {
		if i > 0 && arg == "--" {
			result = append(result, "-x")
			return append(result, commandLine[i:]...)
		}
		result = append(result, arg)
	}
	return append(result, "-x")
}

var pemPattern = regexp.MustCompile(`^(-----BEGIN [A-Z ]+-----)(.*?)(-----END [A-Z ]+-----)$`)

// normalizePrivateKey restores the line breaks of a PEM private key in case they got lost when the key was stored,
// e.g. replaced by literal "\n" sequences or by spaces.
func normalizePrivateKey(key string) string {
	key = strings.TrimSpace(key)
	if strings.Contains(key, "\n") {
		return key + "\n"
	}

	if strings.Contains(key, `\n`) {
		return strings.TrimSpace(strings.ReplaceAll(key, `\n`, "\n")) + "\n"
	}

	if matches := pemPattern.FindStringSubmatch(key); matches != nil {
		body := strings.Join(strings.Fields(matches[2]), "\n")
		return matches[1] + "\n" + body + "\n" + matches[3] + "\n"
	}

	return key
}
