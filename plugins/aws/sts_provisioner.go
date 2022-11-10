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

type StsProvisioner struct {
	TotpCode     string
	SerialNumber string
}

func (p StsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(in.ItemFields[FieldNameDefaultRegion]),
		Credentials: credentials.NewStaticCredentials(in.ItemFields[fieldname.AccessKeyID], in.ItemFields[fieldname.SecretAccessKey], ""),
	})
	if err != nil {
		out.Diagnostics.Errors = append(out.Diagnostics.Errors, sdk.Error{Message: fmt.Sprintf("Could not start aws STS session: %s", err.Error())})
		return
	}
	stsProvider := sts.New(sess)
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(900), // minimum expiration time - 15 minutes
		SerialNumber:    aws.String(p.SerialNumber),
		TokenCode:       aws.String(p.TotpCode),
	}

	result, err := stsProvider.GetSessionToken(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == sts.ErrCodeRegionDisabledException {
				out.Diagnostics.Errors = append(out.Diagnostics.Errors, sdk.Error{Message: fmt.Sprint(sts.ErrCodeRegionDisabledException, aerr.Error())})
			}
		} else {
			out.Diagnostics.Errors = append(out.Diagnostics.Errors, sdk.Error{Message: aerr.Error()})
		}

		return
	}
	out.Environment["AWS_ACCESS_KEY_ID"] = *result.Credentials.AccessKeyId
	out.Environment["AWS_ACCESS_KEY_ID"] = *result.Credentials.AccessKeyId
	out.Environment["AWS_ACCESS_KEY_ID"] = *result.Credentials.SessionToken
}

func (p StsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p StsProvisioner) Description() string {
	return fmt.Sprintf("Provision environment variables with the temporary credentials AWS_ACCESS_KEY_ID, AWS_ACCESS_KEY_ID, AWS_ACCESS_KEY_ID")
}
