package snowflakecli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

var envVarMapping = map[string]sdk.FieldName{
	"SNOWFLAKE_ACCOUNT":         fieldname.Account,
	"SNOWFLAKE_USER":            fieldname.Username,
	"SNOWFLAKE_PRIVATE_KEY_RAW": fieldname.PrivateKey,
	"PRIVATE_KEY_PASSPHRASE":    fieldPrivateKeyPassphrase,
	"SNOWFLAKE_ROLE":            fieldRole,
	"SNOWFLAKE_WAREHOUSE":       fieldWarehouse,
	"SNOWFLAKE_DATABASE":        fieldname.Database,
	"SNOWFLAKE_SCHEMA":          fieldSchema,
}

// privateKeyFileEnvVars are the environment variables pointing to a private key file, in order of precedence.
var privateKeyFileEnvVars = []string{"SNOWFLAKE_PRIVATE_KEY_FILE", "SNOWFLAKE_PRIVATE_KEY_PATH"}

// TrySnowflakeEnvVars imports the credential from the environment variables Snowflake CLI supports. The private key
// is read from the file that SNOWFLAKE_PRIVATE_KEY_FILE or SNOWFLAKE_PRIVATE_KEY_PATH points to, which also allows
// importing a private key file without having to paste its contents.
func TrySnowflakeEnvVars() sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		envVarNames := make([]string, 0, len(envVarMapping)+len(privateKeyFileEnvVars))
		fields := make(map[sdk.FieldName]string)
		for envVarName, fieldName := range envVarMapping {
			envVarNames = append(envVarNames, envVarName)
			if value := os.Getenv(envVarName); value != "" {
				fields[fieldName] = value
			}
		}
		envVarNames = append(envVarNames, privateKeyFileEnvVars...)
		sort.Strings(envVarNames)

		attempt := out.NewAttempt(importer.SourceEnvVars(envVarNames...))

		for _, envVarName := range privateKeyFileEnvVars {
			path := os.Getenv(envVarName)
			if path == "" {
				continue
			}

			contents, err := os.ReadFile(resolvePath(path, in))
			if err != nil {
				attempt.AddError(fmt.Errorf("reading private key file from %s: %w", envVarName, err))
				return
			}
			// Like Snowflake CLI, prefer the private key file over SNOWFLAKE_PRIVATE_KEY_RAW.
			fields[fieldname.PrivateKey] = strings.TrimSpace(string(contents))
			break
		}

		if len(fields) > 0 {
			attempt.AddCandidate(sdk.ImportCandidate{
				Fields: fields,
			})
		}
	}
}

// TrySnowflakeConfigFiles imports key pair connections from the config.toml and connections.toml files in all
// locations Snowflake CLI looks for them.
func TrySnowflakeConfigFiles() sdk.Importer {
	return importer.TryAll(
		TryConfigTOML("~/.snowflake/config.toml"),
		TryConnectionsTOML("~/.snowflake/connections.toml"),
		importer.MacOnly(
			TryConfigTOML("~/Library/Application Support/snowflake/config.toml"),
			TryConnectionsTOML("~/Library/Application Support/snowflake/connections.toml"),
		),
		importer.LinuxOnly(
			TryConfigTOML("~/.config/snowflake/config.toml"),
			TryConnectionsTOML("~/.config/snowflake/connections.toml"),
		),
	)
}

// TryConfigTOML imports connections defined as [connections.<name>] tables in a config.toml file.
func TryConfigTOML(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config struct {
			Connections map[string]map[string]any `toml:"connections"`
		}
		if err := contents.ToTOML(&config); err != nil {
			out.AddError(err)
			return
		}

		addConnectionCandidates(config.Connections, in, out)
	})
}

// TryConnectionsTOML imports connections defined as top-level [<name>] tables in a connections.toml file.
func TryConnectionsTOML(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config map[string]map[string]any
		if err := contents.ToTOML(&config); err != nil {
			out.AddError(err)
			return
		}

		addConnectionCandidates(config, in, out)
	})
}

func addConnectionCandidates(connections map[string]map[string]any, in sdk.ImportInput, out *sdk.ImportAttempt) {
	// Sort the names so that the order of the candidates is deterministic.
	names := make([]string, 0, len(connections))
	for name := range connections {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		fields, err := connectionFields(connections[name], in)
		if err != nil {
			out.AddError(fmt.Errorf("connection %q: %w", name, err))
			continue
		}
		if fields == nil {
			continue
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields:   fields,
			NameHint: importer.SanitizeNameHint(name),
		})
	}
}

// connectionFields maps a connection to the fields of the credential. It returns nil if the connection does not
// use key pair authentication.
func connectionFields(connection map[string]any, in sdk.ImportInput) (map[sdk.FieldName]string, error) {
	account := stringValue(connection, "account")
	user := stringValue(connection, "user")
	if account == "" || user == "" || !strings.EqualFold(stringValue(connection, "authenticator"), "SNOWFLAKE_JWT") {
		return nil, nil
	}

	privateKey := ""
	if path := stringValue(connection, "private_key_file", "private_key_path"); path != "" {
		contents, err := os.ReadFile(resolvePath(path, in))
		if err != nil {
			return nil, fmt.Errorf("reading private key file: %w", err)
		}
		privateKey = string(contents)
	} else {
		privateKey = stringValue(connection, "private_key_raw")
	}
	if privateKey == "" {
		return nil, nil
	}

	fields := map[sdk.FieldName]string{
		fieldname.Account:    account,
		fieldname.Username:   user,
		fieldname.PrivateKey: strings.TrimSpace(privateKey),
	}

	optionalFields := map[sdk.FieldName]string{
		fieldPrivateKeyPassphrase: stringValue(connection, "private_key_passphrase", "private_key_file_pwd"),
		fieldRole:                 stringValue(connection, "role"),
		fieldWarehouse:            stringValue(connection, "warehouse"),
		fieldname.Database:        stringValue(connection, "database"),
		fieldSchema:               stringValue(connection, "schema"),
	}
	for name, value := range optionalFields {
		if value != "" {
			fields[name] = value
		}
	}

	return fields, nil
}

// stringValue returns the first non-empty string value found for the given keys.
func stringValue(connection map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := connection[key].(string); ok && value != "" {
			return value
		}
	}
	return ""
}

func resolvePath(path string, in sdk.ImportInput) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(in.HomeDir, strings.TrimPrefix(path, "~/"))
	}
	if filepath.IsAbs(path) {
		return filepath.Join(in.RootDir, path)
	}
	return path
}
