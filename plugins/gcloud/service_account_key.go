package gcloud

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func ServiceAccountKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.ServiceAccountKey,
		DocsURL:       sdk.URL("https://cloud.google.com/iam/docs/keys-create-delete"),
		ManagementURL: sdk.URL("https://console.cloud.google.com/iam-admin/serviceaccounts"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Credential,
				MarkdownDescription: "The JSON credential content for a service account key or authorized user credential.",
				Secret:              true,
			},
			{
				Name:                fieldname.ProjectID,
				MarkdownDescription: "The default GCP project to use.",
				Optional:            true,
			},
			{
				Name:                fieldname.Account,
				MarkdownDescription: "The account email associated with the credential.",
				Optional:            true,
			},
		},
		DefaultProvisioner: GCPProvisioner{},
		Importer: importer.TryAll(
			TryGCloudApplicationDefaultCredentialsFile(),
			TryGoogleApplicationCredentialsEnvVar(),
		),
	}
}

// GCPProvisioner writes the credential JSON to a temporary file and sets
// GOOGLE_APPLICATION_CREDENTIALS to point at it. This approach works with
// all GCP tools (gcloud, gsutil, bq, client libraries, Terraform).
type GCPProvisioner struct{}

func (p GCPProvisioner) Description() string {
	return "Provision GCP credential as a temporary JSON file and set GOOGLE_APPLICATION_CREDENTIALS"
}

func (p GCPProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	credJSON := in.ItemFields[fieldname.Credential]
	var cred gcpCredentialFile
	if err := json.Unmarshal([]byte(credJSON), &cred); err != nil {
		out.AddError(errInvalidJSON)
		return
	}

	outPath := filepath.Join(in.TempDir, "gcloud-credentials.json")
	out.AddSecretFile(outPath, []byte(credJSON))
	out.AddEnvVar("GOOGLE_APPLICATION_CREDENTIALS", outPath)

	if projectID, ok := in.ItemFields[fieldname.ProjectID]; ok && projectID != "" {
		out.AddEnvVar("CLOUDSDK_CORE_PROJECT", projectID)
	} else if cred.ProjectID != "" {
		out.AddEnvVar("CLOUDSDK_CORE_PROJECT", cred.ProjectID)
	}

	if account, ok := in.ItemFields[fieldname.Account]; ok && account != "" {
		out.AddEnvVar("CLOUDSDK_CORE_ACCOUNT", account)
	} else if cred.ClientEmail != "" {
		out.AddEnvVar("CLOUDSDK_CORE_ACCOUNT", cred.ClientEmail)
	}
}

func (p GCPProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Temp files are automatically cleaned up by the SDK.
}

var errInvalidJSON = &jsonError{}

type jsonError struct{}

func (e *jsonError) Error() string {
	return "credential field does not contain valid JSON"
}

// TryGCloudApplicationDefaultCredentialsFile imports credentials from the
// gcloud application default credentials file at ~/.config/gcloud/application_default_credentials.json.
func TryGCloudApplicationDefaultCredentialsFile() sdk.Importer {
	return importer.TryFile("~/.config/gcloud/application_default_credentials.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var cred gcpCredentialFile
		if err := contents.ToJSON(&cred); err != nil {
			out.AddError(err)
			return
		}

		if cred.Type == "" {
			return
		}

		candidate := sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.Credential: contents.ToString(),
			},
		}

		if cred.ProjectID != "" {
			candidate.Fields[fieldname.ProjectID] = cred.ProjectID
		}

		if cred.ClientEmail != "" {
			candidate.NameHint = importer.SanitizeNameHint(cred.ClientEmail)
		}

		out.AddCandidate(candidate)
	})
}

// TryGoogleApplicationCredentialsEnvVar imports credentials from the file
// pointed to by the GOOGLE_APPLICATION_CREDENTIALS environment variable.
func TryGoogleApplicationCredentialsEnvVar() sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		attempt := out.NewAttempt(importer.SourceEnvVars("GOOGLE_APPLICATION_CREDENTIALS"))

		filePath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if filePath == "" {
			return
		}

		contents, err := os.ReadFile(filePath)
		if err != nil {
			attempt.AddError(err)
			return
		}

		var cred gcpCredentialFile
		if err := json.Unmarshal(contents, &cred); err != nil {
			attempt.AddError(err)
			return
		}

		if cred.Type == "" {
			return
		}

		candidate := sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.Credential: string(contents),
			},
		}

		if cred.ProjectID != "" {
			candidate.Fields[fieldname.ProjectID] = cred.ProjectID
		}

		if cred.ClientEmail != "" {
			candidate.NameHint = importer.SanitizeNameHint(cred.ClientEmail)
		}

		attempt.AddCandidate(candidate)
	}
}

// gcpCredentialFile represents the minimal structure of a GCP credential JSON file,
// supporting both service_account and authorized_user types.
type gcpCredentialFile struct {
	Type        string `json:"type"`
	ProjectID   string `json:"project_id,omitempty"`
	ClientEmail string `json:"client_email,omitempty"`
}
