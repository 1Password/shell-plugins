package wireguard

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

func AccessConfig() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.AccessConfig,
		DocsURL: sdk.URL("https://www.wireguard.com/quickstart/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.PrivateKey,
				MarkdownDescription: "Private Key of your client.",
				Secret:              true,
			},
			{
				Name:                fieldname.Address,
				MarkdownDescription: "VPN IP address of your client.",
				Secret:              false,
			},
			{
				Name:                fieldname.PublicKey,
				MarkdownDescription: "Public Key of the Peer.",
				Secret:              false,
			},
			{
				Name:                fieldname.Endpoint,
				MarkdownDescription: "Endpoint address of the Peer.",
				Secret:              false,
			},
			{
				Name:                fieldname.AllowedIPs,
				MarkdownDescription: "Allowed IPs of the Peer.",
				Secret:              false,
			},
		},
		DefaultProvisioner: provision.TempFile(
			wgConfig,
			provision.Filename("wg0.conf"),
			provision.AddArgs("{{ .Path }}"),
		),
		Importer: importer.TryAll(
			TryWireguardVPNConfigFile("/etc/wireguard/wg0.conf"),
		)}
}

func wgConfig(in sdk.ProvisionInput) ([]byte, error) {
	content := "[Interface]\n"

	if privateKey, ok := in.ItemFields[fieldname.PrivateKey]; ok {
		content += configFileEntry("PrivateKey", privateKey)
	}

	if address, ok := in.ItemFields[fieldname.Address]; ok {
		content += configFileEntry("Address", address)
	}

	content = content + "\n"
	content = content + "[Peer]\n"

	if publicKey, ok := in.ItemFields[fieldname.PublicKey]; ok {
		content += configFileEntry("PublicKey", publicKey)
	}

	if endpoint, ok := in.ItemFields[fieldname.Endpoint]; ok {
		content += configFileEntry("Endpoint", endpoint)
	}

	if allowedIPs, ok := in.ItemFields[fieldname.AllowedIPs]; ok {
		content += configFileEntry("AllowedIPs", allowedIPs)
	}

	return []byte(content), nil
}

func configFileEntry(key string, value string) string {
	return fmt.Sprintf("%s = %s\n", key, value)
}

func TryWireguardVPNConfigFile(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		configFile, err := contents.ToINI()
		if err != nil {
			out.AddError(err)
			return
		}

		fields := make(map[sdk.FieldName]string)
		for _, section := range configFile.Sections() {
			if section.HasKey("PrivateKey") && section.Key("PrivateKey").Value() != "" {
				fields[fieldname.PrivateKey] = section.Key("PrivateKey").Value()
			}
			if section.HasKey("Address") && section.Key("Address").Value() != "" {
				fields[fieldname.Address] = section.Key("Address").Value()
			}
			if section.HasKey("PublicKey") && section.Key("PublicKey").Value() != "" {
				fields[fieldname.PublicKey] = section.Key("PublicKey").Value()
			}
			if section.HasKey("Endpoint") && section.Key("Endpoint").Value() != "" {
				fields[fieldname.Endpoint] = section.Key("Endpoint").Value()
			}
			if section.HasKey("AllowedIPs") && section.Key("AllowedIPs").Value() != "" {
				fields[fieldname.AllowedIPs] = section.Key("AllowedIPs").Value()
			}
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: fields,
		})
	})
}
