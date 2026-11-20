package kubectl

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func Credentials() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.Credentials,
		DocsURL:       sdk.URL("https://kubernetes.io/docs/reference/access-authn-authz/authentication/"),
		ManagementURL: nil,
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Address,
				MarkdownDescription: "Kubernetes API server URL, as in the cluster's `server` field (for example https://prod.example.com:6443).",
			},
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Bearer token used to authenticate to the Kubernetes API server.",
				Secret:              true,
				Optional:            true,
			},
			{
				Name:                fieldname.Certificate,
				MarkdownDescription: "Client certificate used to authenticate to the Kubernetes API server, as PEM or as the base64-encoded PEM found in `client-certificate-data`.",
				Optional:            true,
			},
			{
				Name:                fieldname.PrivateKey,
				MarkdownDescription: "Private key of the client certificate, as PEM or as the base64-encoded PEM found in `client-key-data`.",
				Secret:              true,
				Optional:            true,
			},
			{
				Name:                fieldname.CertificateAuthority,
				MarkdownDescription: "Certificate authority used to verify the Kubernetes API server, as PEM or as the base64-encoded PEM found in `certificate-authority-data`. Used only when no kubeconfig exists on the machine.",
				Optional:            true,
			},
		},
		DefaultProvisioner: Provisioner(),
		Importer:           TryKubeconfigFiles(),
	}
}
