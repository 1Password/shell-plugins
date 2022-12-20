package aws

import (
	"context"

	"github.com/1Password/shell-plugins/plugins/docker"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

func ECRProvisioner() sdk.Provisioner {
	return docker.ConfigFileProvisioner(
		docker.WithURLField(fieldname.URL),
		docker.WithStaticUsername("AWS"),
		docker.WithPasswordFunc(func(in sdk.ProvisionInput) (string, error) {
			config := aws.NewConfig()
			config.Credentials = credentials.NewStaticCredentialsProvider(in.ItemFields[fieldname.AccessKeyID], in.ItemFields[fieldname.SecretAccessKey], "")

			ecrProvider := ecr.NewFromConfig(*config)
			input := &ecr.GetAuthorizationTokenInput{}

			result, err := ecrProvider.GetAuthorizationToken(context.TODO(), input)
			if err != nil {
				return "", err
			}

			if len(result.AuthorizationData) > 0 {
				ecrToken := *result.AuthorizationData[0].AuthorizationToken
				return ecrToken, nil
			}

			return "", nil
		}),
	)
}
