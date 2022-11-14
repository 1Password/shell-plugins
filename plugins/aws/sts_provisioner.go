package aws

import (
	"context"
	"fmt"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type STSProvisioner struct {
	TOTPCode  string
	MFASerial string
}

func (p STSProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(in.ItemFields[FieldNameDefaultRegion]),
		Credentials: credentials.NewStaticCredentials(in.ItemFields[fieldname.AccessKeyID], in.ItemFields[fieldname.SecretAccessKey], ""),
	})
	if err != nil {
		out.AddError(fmt.Errorf("could not start aws STS session: %s", err))
		return
	}
	stsProvider := sts.New(sess)
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(900), // minimum expiration time - 15 minutes
		SerialNumber:    aws.String(p.MFASerial),
		TokenCode:       aws.String(p.TOTPCode),
	}

	result, err := stsProvider.GetSessionToken(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == sts.ErrCodeRegionDisabledException {
				out.AddError(fmt.Errorf(sts.ErrCodeRegionDisabledException+": %s", aerr.Error()))
			}
		} else {
			out.AddError(aerr)
		}

		return
	}
	out.AddEnvVar("AWS_ACCESS_KEY_ID", *result.Credentials.AccessKeyId)
	out.AddEnvVar("AWS_SECRET_ACCESS_KEY", *result.Credentials.SecretAccessKey)
	out.AddEnvVar("AWS_SESSION_TOKEN", *result.Credentials.SessionToken)
	out.AddEnvVar("AWS_DEFAULT_REGION", in.ItemFields[FieldNameDefaultRegion])
}

func (p STSProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p STSProvisioner) Description() string {
	return fmt.Sprintf("Provision environment variables with the temporary credentials AWS_ACCESS_KEY_ID, AWS_ACCESS_KEY_ID, AWS_ACCESS_KEY_ID")
}
