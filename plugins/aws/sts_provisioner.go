package aws

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

type STSProvisioner struct {
	TOTPCode  string
	MFASerial string
}

func (p STSProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	var cached types.Credentials
	if ok := in.Cache.Get("sts", &cached); ok {
		out.AddEnvVar("AWS_ACCESS_KEY_ID", *cached.AccessKeyId)
		out.AddEnvVar("AWS_SECRET_ACCESS_KEY", *cached.SecretAccessKey)
		out.AddEnvVar("AWS_SESSION_TOKEN", *cached.SessionToken)

		if region, ok := in.ItemFields[fieldname.DefaultRegion]; ok {
			out.AddEnvVar("AWS_DEFAULT_REGION", region)
		}

		return
	}

	config := aws.NewConfig()
	config.Credentials = credentials.NewStaticCredentialsProvider(in.ItemFields[fieldname.AccessKeyID], in.ItemFields[fieldname.SecretAccessKey], "")

	region, ok := in.ItemFields[fieldname.DefaultRegion]
	if !ok {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}
	if len(region) == 0 {
		out.AddError(errors.New("region is required for the AWS Shell Plugin MFA workflow: set 'default region' in 1Password or set the 'AWS_DEFAULT_REGION' environment variable yourself"))
		return
	}
	config.Region = region

	roleArn, ok := in.ItemFields[fieldname.Role]
	stsProvider := sts.NewFromConfig(*config)

	var awsTemporaryCredentials *types.Credentials

	if ok {
		input := &sts.AssumeRoleInput{
			SerialNumber:    aws.String(p.MFASerial),
			TokenCode:       aws.String(p.TOTPCode),
			RoleArn:         aws.String(roleArn),
			RoleSessionName: aws.String("1password-shell-plugin"),
		}

		resp, err := stsProvider.AssumeRole(ctx, input)
		if err != nil {
			out.AddError(err)
			return
		}

		awsTemporaryCredentials = resp.Credentials
	} else {
		input := &sts.GetSessionTokenInput{
			DurationSeconds: aws.Int32(900), // minimum expiration time - 15 minutes
			SerialNumber:    aws.String(p.MFASerial),
			TokenCode:       aws.String(p.TOTPCode),
		}

		resp, err := stsProvider.GetSessionToken(ctx, input)
		if err != nil {
			out.AddError(err)
			return
		}

		awsTemporaryCredentials = resp.Credentials
	}

	out.AddEnvVar("AWS_ACCESS_KEY_ID", *awsTemporaryCredentials.AccessKeyId)
	out.AddEnvVar("AWS_SECRET_ACCESS_KEY", *awsTemporaryCredentials.SecretAccessKey)
	out.AddEnvVar("AWS_SESSION_TOKEN", *awsTemporaryCredentials.SessionToken)
	out.AddEnvVar("AWS_DEFAULT_REGION", region)

	err := out.Cache.Put("sts", *awsTemporaryCredentials, *awsTemporaryCredentials.Expiration)
	if err != nil {
		out.AddError(fmt.Errorf("failed to serialize aws sts credentials: %w", err))
	}
}

func (p STSProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p STSProvisioner) Description() string {
	return "Provision environment variables with temporary STS credentials AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN"
}
