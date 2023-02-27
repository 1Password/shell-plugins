package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	confighelpers "github.com/99designs/aws-vault/v7/vault"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

type assumeRoleProvisioner struct {
	mfaProvisioner
	RoleARN           string
	RoleSessionName   string
	DurationSeconds   time.Duration
	ExternalID        string
	SourceIdentity    string
	Tags              map[string]string
	TransitiveTagKeys []string
}

func newAssumeRoleProvisioner(mfaSerial string, totpCode string, awsConfig *confighelpers.Config) sdk.Provisioner {
	if awsConfig.RoleSessionName == "" {
		awsConfig.RoleSessionName = "1password-shell-plugin"
	}
	if awsConfig.AssumeRoleDuration == 0 {
		awsConfig.AssumeRoleDuration = 900 * time.Second // minimum expiration time - 15 minutes
	}
	return assumeRoleProvisioner{
		mfaProvisioner: mfaProvisioner{
			TOTPCode:  totpCode,
			MFASerial: mfaSerial,
			Region:    awsConfig.Region,
		},
		RoleARN:           awsConfig.RoleARN,
		RoleSessionName:   awsConfig.RoleSessionName,
		DurationSeconds:   awsConfig.AssumeRoleDuration,
		ExternalID:        awsConfig.ExternalID,
		SourceIdentity:    awsConfig.SourceIdentity,
		Tags:              awsConfig.SessionTags,
		TransitiveTagKeys: awsConfig.TransitiveSessionTags,
	}
}

func (p assumeRoleProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	config := aws.NewConfig()
	config.Credentials = credentials.NewStaticCredentialsProvider(in.ItemFields[fieldname.AccessKeyID], in.ItemFields[fieldname.SecretAccessKey], "")
	config.Region = p.Region
	// Region from 1Password has priority
	region, ok := in.ItemFields[fieldname.DefaultRegion]
	if ok {
		config.Region = region
	}

	stsProvider := sts.NewFromConfig(*config)

	if ok := tryUsingCachedCredentials(assumeRoleCacheKey, in, out); ok {
		return
	}

	input := &sts.AssumeRoleInput{
		DurationSeconds: aws.Int32(int32(p.DurationSeconds.Seconds())),
		RoleArn:         aws.String(p.RoleARN),
		RoleSessionName: aws.String(p.RoleSessionName),
	}
	if p.MFASerial != "" && p.TOTPCode != "" {
		input.SerialNumber = aws.String(p.MFASerial)
		input.TokenCode = aws.String(p.TOTPCode)
	}
	if p.SourceIdentity != "" {
		input.SourceIdentity = aws.String(p.SourceIdentity)
	}
	if p.ExternalID != "" {
		input.ExternalId = aws.String(p.ExternalID)
	}
	if len(p.Tags) > 0 {
		input.Tags = make([]types.Tag, 0)
		for key, value := range p.Tags {
			tag := types.Tag{
				Key:   aws.String(key),
				Value: aws.String(value),
			}
			input.Tags = append(input.Tags, tag)
		}
	}
	if len(p.TransitiveTagKeys) > 0 {
		input.TransitiveTagKeys = p.TransitiveTagKeys
	}

	resp, err := stsProvider.AssumeRole(ctx, input)
	if err != nil {
		out.AddError(err)
		return
	}

	out.AddEnvVar("AWS_ACCESS_KEY_ID", *resp.Credentials.AccessKeyId)
	out.AddEnvVar("AWS_SECRET_ACCESS_KEY", *resp.Credentials.SecretAccessKey)
	out.AddEnvVar("AWS_SESSION_TOKEN", *resp.Credentials.SessionToken)
	out.AddEnvVar("AWS_DEFAULT_REGION", region)
	err = out.Cache.Put(assumeRoleCacheKey, resp.Credentials, *resp.Credentials.Expiration)
	if err != nil {
		out.AddError(fmt.Errorf("failed to serialize aws sts credentials: %w", err))
	}
}

func (p assumeRoleProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here
}

func (p assumeRoleProvisioner) Description() string {
	return "Provision environment variables with temporary STS credentials AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN, when a role is specified."
}
