package backupdrill

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

// The backup bucket uses the standard S3 field names. The Supabase Storage keys are a second
// S3 key pair, so they get their own prefixed names to fit in the same item.
const (
	fieldDatabaseURL                    = sdk.FieldName("Database URL")
	fieldBucket                         = sdk.FieldName("Bucket")
	fieldSupabaseStorageEndpoint        = sdk.FieldName("Supabase Storage Endpoint")
	fieldSupabaseStorageRegion          = sdk.FieldName("Supabase Storage Region")
	fieldSupabaseStorageAccessKeyID     = sdk.FieldName("Supabase Storage Access Key ID")
	fieldSupabaseStorageSecretAccessKey = sdk.FieldName("Supabase Storage Secret Access Key")
)

// Credentials holds everything backupdrill reads from the environment to authenticate: the source
// database, the backup bucket, and optionally the Supabase Storage S3 keys. A plugin can only have
// one credential type, so they share one item.
//
// Only the bucket keys are required. `drill`, `restore` and `exporter` work without the database
// URL, so a machine that just runs drills doesn't need to hold the database password. Settings that
// aren't secrets (key prefix, project name, schemas) stay in backupdrill.config.json or flags.
func Credentials() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.Credentials,
		DocsURL:       sdk.URL("https://backupdrill.com/docs/cli#configure"),
		ManagementURL: sdk.URL("https://backupdrill.com/docs/buckets"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldDatabaseURL,
				MarkdownDescription: "Connection string for the Supabase Postgres database to back up, such as the Session pooler string. Used by `backup` and `estimate`.",
				Secret:              true,
				Optional:            true,
			},
			{
				Name:                fieldBucket,
				MarkdownDescription: "Name of the S3-compatible bucket that backups are written to and read from.",
				Optional:            true,
			},
			{
				Name:                fieldname.AccessKeyID,
				MarkdownDescription: "Access key ID for the backup bucket.",
			},
			{
				Name:                fieldname.SecretAccessKey,
				MarkdownDescription: "Secret access key for the backup bucket.",
				Secret:              true,
			},
			{
				Name:                fieldname.Endpoint,
				MarkdownDescription: "S3 endpoint of the backup bucket. Needed for Cloudflare R2, Backblaze B2 and other S3-compatible providers; leave empty for AWS S3.",
				Optional:            true,
			},
			{
				Name:                fieldname.Region,
				MarkdownDescription: "Region of the backup bucket. Defaults to `auto`.",
				Optional:            true,
			},
			{
				Name:                fieldSupabaseStorageEndpoint,
				MarkdownDescription: "S3 endpoint of the project's Supabase Storage, e.g. `https://<ref>.storage.supabase.co/storage/v1/s3`. Set it together with both Supabase Storage keys to back up Storage files, or leave all three empty to back up the database only.",
				Optional:            true,
			},
			{
				Name:                fieldSupabaseStorageRegion,
				MarkdownDescription: "Region of the Supabase project, used to sign Storage requests. Defaults to `auto`.",
				Optional:            true,
			},
			{
				Name:                fieldSupabaseStorageAccessKeyID,
				MarkdownDescription: "Supabase Storage S3 access key ID.",
				Optional:            true,
			},
			{
				Name:                fieldSupabaseStorageSecretAccessKey,
				MarkdownDescription: "Supabase Storage S3 secret access key.",
				Secret:              true,
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		// backupdrill also reads a JSON config file, but it lives in each project directory (or
		// wherever --config points), so there's no fixed path for a file importer to scan.
		Importer: importer.TryEnvVarPair(defaultEnvVarMapping),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"BACKUPDRILL_DATABASE_URL":                       fieldDatabaseURL,
	"BACKUPDRILL_S3_BUCKET":                          fieldBucket,
	"BACKUPDRILL_S3_ACCESS_KEY_ID":                   fieldname.AccessKeyID,
	"BACKUPDRILL_S3_SECRET_ACCESS_KEY":               fieldname.SecretAccessKey,
	"BACKUPDRILL_S3_ENDPOINT":                        fieldname.Endpoint,
	"BACKUPDRILL_S3_REGION":                          fieldname.Region,
	"BACKUPDRILL_SUPABASE_STORAGE_ENDPOINT":          fieldSupabaseStorageEndpoint,
	"BACKUPDRILL_SUPABASE_STORAGE_REGION":            fieldSupabaseStorageRegion,
	"BACKUPDRILL_SUPABASE_STORAGE_ACCESS_KEY_ID":     fieldSupabaseStorageAccessKeyID,
	"BACKUPDRILL_SUPABASE_STORAGE_SECRET_ACCESS_KEY": fieldSupabaseStorageSecretAccessKey,
}
