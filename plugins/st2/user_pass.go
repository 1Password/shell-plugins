package st2

import (
	"context"
	"fmt"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func UserPass() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.UserPass,
		DocsURL:       sdk.URL("https://docs.stackstorm.com/authentication.html"),
		ManagementURL: nil,
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Username,
				MarkdownDescription: "Username used to authenticate to StackStorm.",
				Optional:            false,
				Secret:              false,
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password used to authenticate to StackStorm.",
				Optional:            false,
				Secret:              true,
			},
			{
				Name:                fieldname.Website,
				MarkdownDescription: "StackStorm base URL.",
				Optional:            false,
				Secret:              false,
			},

		},
		DefaultProvisioner: provision.TempFile(st2Config, provision.AtFixedPath("~/.st2/config"), provision.AddArgs("--config-file {{ .Path }}")),
		Importer: importer.TryAll(
			TryStackStormConfigFile("~/.st2/config"),
		),
	}
}

func st2Config(in sdk.ProvisionInput) ([]byte, error) {
	content := "[general]\n"
	if baseurl, ok := in.ItemFields[fieldname.Website]; ok {
		content += configFileEntry("base_url", baseurl)
		content += "[cli]\n"
		content += configFileEntry("cache_token", "False")
		content += "[api]\n"
		content += configFileEntry("url", baseurl + "/api/v1")
		content += "[auth]\n"
		content += configFileEntry("url", baseurl + "/auth/v1")
		content += "[stream]\n"
		content += configFileEntry("url", baseurl + "/stream/v1")
		content += "[credentials]\n"
		if username, ok := in.ItemFields[fieldname.Username]; ok {
			content += configFileEntry("username", username)
		}
		if password, ok := in.ItemFields[fieldname.Password]; ok {
			content += configFileEntry("password", password)
		}
		content += "\n"
	}

	return []byte(content), nil
}

func TryStackStormConfigFile(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		credentialsFile, err := contents.ToINI()
		if err != nil {
			out.AddError(err)
			return
		}

		fields := make(map[sdk.FieldName]string)
		for _, section := range credentialsFile.Sections() {
			if section.HasKey("username") && section.Key("username").Value() != "" {
				fields[fieldname.Username] = section.Key("username").Value()
			}
			if section.HasKey("password") && section.Key("password").Value() != "" {
				fields[fieldname.Password] = section.Key("password").Value()
			}
			if section.HasKey("base_url") && section.Key("base_url").Value() != "" {
				fields[fieldname.Website] = section.Key("base_url").Value()
			}
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: fields,
		})
	})
}

func configFileEntry(key string, value string) string {
	return fmt.Sprintf("%s = %s\n", key, value)
}
