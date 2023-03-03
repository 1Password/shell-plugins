package aws

import (
	"context"
	"fmt"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	confighelpers "github.com/99designs/aws-vault/v7/vault"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func newAssumeRoleProvider(client *sts.Client, awsConfig *confighelpers.Config) aws.CredentialsProvider {
	return &confighelpers.AssumeRoleProvider{
		StsClient:         client,
		RoleARN:           awsConfig.RoleARN,
		RoleSessionName:   awsConfig.RoleSessionName,
		ExternalID:        awsConfig.ExternalID,
		Duration:          awsConfig.AssumeRoleDuration,
		Tags:              awsConfig.SessionTags,
		TransitiveTagKeys: awsConfig.TransitiveSessionTags,
		SourceIdentity:    awsConfig.SourceIdentity,
		Mfa:               confighelpers.NewMfa(awsConfig),
	}
}

func newSessionTokenProvider(client *sts.Client, awsConfig *confighelpers.Config) aws.CredentialsProvider {
	return &confighelpers.SessionTokenProvider{
		StsClient: client,
		Duration:  awsConfig.NonChainedGetSessionTokenDuration,
		Mfa:       confighelpers.NewMfa(awsConfig),
	}
}

func newStsCacheProvider(key string, cache sdk.CacheState) aws.CredentialsProvider {
	return StsCacheProvider{
		awsCacheKey: key,
		cache:       cache,
	}
}

type MasterAwsCredentialsProvider struct {
	itemFields map[sdk.FieldName]string
}

func (p MasterAwsCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	secret, hasSecretKey := p.itemFields[fieldname.SecretAccessKey]
	keyId, hasKeyId := p.itemFields[fieldname.AccessKeyID]

	if !hasKeyId || !hasSecretKey {
		return aws.Credentials{}, fmt.Errorf("no master credentials found. Please add your Access Key Id and Secret Access Key to your 1Password item's fields")
	}

	return aws.Credentials{
		AccessKeyID:     keyId,
		SecretAccessKey: secret,
	}, nil
}

func newMasterCredentialsProvider(itemFields map[sdk.FieldName]string) aws.CredentialsProvider {
	return MasterAwsCredentialsProvider{itemFields: itemFields}
}
