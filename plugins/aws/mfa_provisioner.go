package aws

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

const (
	mfaCacheKey        = "sts-mfa"
	assumeRoleCacheKey = "sts-assume-role"
)

type mfaProvisioner struct {
	TOTPCode        string
	MFASerial       string
	Region          string
	DurationSeconds time.Duration
}

func newMFAProvisioner(mfaSerial string, totpCode string, region string, durationSeconds time.Duration) sdk.Provisioner {
	if durationSeconds == 0 {
		durationSeconds = 900 * time.Second // minimum expiration time - 15 minutes
	}
	return mfaProvisioner{
		TOTPCode:        totpCode,
		MFASerial:       mfaSerial,
		Region:          region,
		DurationSeconds: durationSeconds,
	}
}

func (p mfaProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	config := aws.NewConfig()
	config.Credentials = credentials.NewStaticCredentialsProvider(in.ItemFields[fieldname.AccessKeyID], in.ItemFields[fieldname.SecretAccessKey], "")

	region, ok := in.ItemFields[fieldname.DefaultRegion]
	if !ok {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}
	if len(region) == 0 {
		out.AddError(errors.New("region is required for the AWS Shell Plugin MFA or Assume Role workflows: set 'default region' in 1Password or set the 'AWS_DEFAULT_REGION' environment variable yourself"))
		return
	}
	config.Region = region

	stsProvider := sts.NewFromConfig(*config)

	if useCache := tryUsingCachedCredentials(mfaCacheKey, in, out); useCache {
		return
	}
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int32(int32(p.DurationSeconds.Seconds())),
		SerialNumber:    aws.String(p.MFASerial),
		TokenCode:       aws.String(p.TOTPCode),
	}

	resp, err := stsProvider.GetSessionToken(ctx, input)
	if err != nil {
		out.AddError(err)
		return
	}
	out.AddEnvVar("AWS_ACCESS_KEY_ID", *resp.Credentials.AccessKeyId)
	out.AddEnvVar("AWS_SECRET_ACCESS_KEY", *resp.Credentials.SecretAccessKey)
	out.AddEnvVar("AWS_SESSION_TOKEN", *resp.Credentials.SessionToken)
	err = out.Cache.Put(mfaCacheKey, *resp.Credentials, *resp.Credentials.Expiration)
	if err != nil {
		out.AddError(fmt.Errorf("failed to serialize aws sts credentials: %w", err))
	}
	out.AddEnvVar("AWS_DEFAULT_REGION", region)

}

func tryUsingCachedCredentials(cacheKey string, in sdk.ProvisionInput, out *sdk.ProvisionOutput) bool {
	var cached types.Credentials
	if ok := in.Cache.Get(cacheKey, &cached); ok {

		out.AddEnvVar("AWS_ACCESS_KEY_ID", *cached.AccessKeyId)
		out.AddEnvVar("AWS_SECRET_ACCESS_KEY", *cached.SecretAccessKey)
		out.AddEnvVar("AWS_SESSION_TOKEN", *cached.SessionToken)

		if region, ok := in.ItemFields[fieldname.DefaultRegion]; ok {
			out.AddEnvVar("AWS_DEFAULT_REGION", region)
		}

		return true
	}
	return false
}

func (p mfaProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// nothing to do here: 1Password CLI removes env vars
}

func (p mfaProvisioner) Description() string {
	return "Provision environment variables with temporary STS credentials AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN"
}
