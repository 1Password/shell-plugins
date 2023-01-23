package snowflake

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func Credentials() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.Credentials,
		DocsURL: sdk.URL("https://docs.snowflake.com/en/user-guide/snowsql.html"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Account,
				MarkdownDescription: "Snowflake account name.",
				Optional:            false,
				Secret:              false,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
			{
				Name:                fieldname.Username,
				MarkdownDescription: "Snowflake username.",
				Optional:            false,
				Secret:              false,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password used to authenticate to Snowflake account.",
				Optional:            false,
				Secret:              true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TrySnowflakeConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"SNOWSQL_ACCOUNT": fieldname.Account,
	"SNOWSQL_USER":    fieldname.Username,
	"SNOWSQL_PWD":     fieldname.Password,
}

func TrySnowflakeConfigFile() sdk.Importer {
	return importer.TryFile("~/.snowsql/config", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		credentialsFile, err := contents.ToINI()
		if err != nil {
			out.AddError(err)
			return
		}

		// The generated ~/.snowsql/config file, by default, includes an uncommented [connections.example] credential section below the actual
		// [connections] section. We save the accountname, username, and password listed under [connections] then stop parsing the file
		fields := make(map[sdk.FieldName]string)
		for _, section := range credentialsFile.Sections() {
			if section.HasKey("accountname") && section.Key("accountname").Value() != "" {
				fields[fieldname.Account] = section.Key("accountname").Value()
			}

			if section.HasKey("username") && section.Key("username").Value() != "" {
				fields[fieldname.Username] = section.Key("username").Value()
			}

			if section.HasKey("password") && section.Key("password").Value() != "" {
				fields[fieldname.Password] = section.Key("password").Value()
			}

			// Only add candidates with all required credential fields
			if fields[fieldname.Account] != "" && fields[fieldname.Username] != "" && fields[fieldname.Password] != "" {
				out.AddCandidate(sdk.ImportCandidate{
					Fields: fields,
				})
				break
			}
		}
	})
}
