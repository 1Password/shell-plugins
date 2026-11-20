package backupdrill

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

const (
	exampleDatabaseURL     = "postgresql://postgres.exampleref:EXAMPLEPASSWORD@aws-0-us-east-1.pooler.supabase.com:5432/postgres"
	exampleAccessKeyID     = "EXAMPLEACCESSKEYID000000000000000"
	exampleSecretAccessKey = "EXAMPLESECRETACCESSKEY0000000000000000000000000000000000000000000"
	exampleEndpoint        = "https://exampleaccount.r2.cloudflarestorage.com"
	exampleStorageEndpoint = "https://exampleref.storage.supabase.co/storage/v1/s3"
	exampleStorageKeyID    = "EXAMPLESTORAGEKEYID0000000000000"
	exampleStorageSecret   = "EXAMPLESTORAGESECRET00000000000000000000000000000000000000000000"
)

func TestCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, Credentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"all fields": {
			ItemFields: map[sdk.FieldName]string{
				fieldDatabaseURL:                    exampleDatabaseURL,
				fieldBucket:                         "my-backups",
				fieldname.AccessKeyID:               exampleAccessKeyID,
				fieldname.SecretAccessKey:           exampleSecretAccessKey,
				fieldname.Endpoint:                  exampleEndpoint,
				fieldname.Region:                    "auto",
				fieldSupabaseStorageEndpoint:        exampleStorageEndpoint,
				fieldSupabaseStorageRegion:          "us-east-1",
				fieldSupabaseStorageAccessKeyID:     exampleStorageKeyID,
				fieldSupabaseStorageSecretAccessKey: exampleStorageSecret,
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"BACKUPDRILL_DATABASE_URL":                       exampleDatabaseURL,
					"BACKUPDRILL_S3_BUCKET":                          "my-backups",
					"BACKUPDRILL_S3_ACCESS_KEY_ID":                   exampleAccessKeyID,
					"BACKUPDRILL_S3_SECRET_ACCESS_KEY":               exampleSecretAccessKey,
					"BACKUPDRILL_S3_ENDPOINT":                        exampleEndpoint,
					"BACKUPDRILL_S3_REGION":                          "auto",
					"BACKUPDRILL_SUPABASE_STORAGE_ENDPOINT":          exampleStorageEndpoint,
					"BACKUPDRILL_SUPABASE_STORAGE_REGION":            "us-east-1",
					"BACKUPDRILL_SUPABASE_STORAGE_ACCESS_KEY_ID":     exampleStorageKeyID,
					"BACKUPDRILL_SUPABASE_STORAGE_SECRET_ACCESS_KEY": exampleStorageSecret,
				},
			},
		},
		// Fields missing from the item must not be set at all, so values from
		// backupdrill.config.json still apply.
		"bucket keys only": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AccessKeyID:     exampleAccessKeyID,
				fieldname.SecretAccessKey: exampleSecretAccessKey,
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"BACKUPDRILL_S3_ACCESS_KEY_ID":     exampleAccessKeyID,
					"BACKUPDRILL_S3_SECRET_ACCESS_KEY": exampleSecretAccessKey,
				},
			},
		},
	})
}

func TestCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, Credentials().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"BACKUPDRILL_DATABASE_URL":                       exampleDatabaseURL,
				"BACKUPDRILL_S3_BUCKET":                          "my-backups",
				"BACKUPDRILL_S3_ACCESS_KEY_ID":                   exampleAccessKeyID,
				"BACKUPDRILL_S3_SECRET_ACCESS_KEY":               exampleSecretAccessKey,
				"BACKUPDRILL_S3_ENDPOINT":                        exampleEndpoint,
				"BACKUPDRILL_S3_REGION":                          "auto",
				"BACKUPDRILL_SUPABASE_STORAGE_ENDPOINT":          exampleStorageEndpoint,
				"BACKUPDRILL_SUPABASE_STORAGE_REGION":            "us-east-1",
				"BACKUPDRILL_SUPABASE_STORAGE_ACCESS_KEY_ID":     exampleStorageKeyID,
				"BACKUPDRILL_SUPABASE_STORAGE_SECRET_ACCESS_KEY": exampleStorageSecret,
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldDatabaseURL:                    exampleDatabaseURL,
						fieldBucket:                         "my-backups",
						fieldname.AccessKeyID:               exampleAccessKeyID,
						fieldname.SecretAccessKey:           exampleSecretAccessKey,
						fieldname.Endpoint:                  exampleEndpoint,
						fieldname.Region:                    "auto",
						fieldSupabaseStorageEndpoint:        exampleStorageEndpoint,
						fieldSupabaseStorageRegion:          "us-east-1",
						fieldSupabaseStorageAccessKeyID:     exampleStorageKeyID,
						fieldSupabaseStorageSecretAccessKey: exampleStorageSecret,
					},
				},
			},
		},
		"environment without Supabase Storage": {
			Environment: map[string]string{
				"BACKUPDRILL_DATABASE_URL":         exampleDatabaseURL,
				"BACKUPDRILL_S3_BUCKET":            "my-backups",
				"BACKUPDRILL_S3_ACCESS_KEY_ID":     exampleAccessKeyID,
				"BACKUPDRILL_S3_SECRET_ACCESS_KEY": exampleSecretAccessKey,
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldDatabaseURL:          exampleDatabaseURL,
						fieldBucket:               "my-backups",
						fieldname.AccessKeyID:     exampleAccessKeyID,
						fieldname.SecretAccessKey: exampleSecretAccessKey,
					},
				},
			},
		},
	})
}
