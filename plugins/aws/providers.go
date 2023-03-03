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

	//return stscreds.NewAssumeRoleProvider(client, awsConfig.RoleARN, func(options *stscreds.AssumeRoleOptions) {
	//	options.SerialNumber = aws.String(awsConfig.MfaSerial)
	//	options.TokenProvider = func() (string, error) {
	//		return awsConfig.MfaToken, nil
	//	}
	//	options.RoleARN = awsConfig.RoleARN
	//	options.RoleSessionName = awsConfig.RoleSessionName // gets substituted by randomly generated string if empty
	//	options.Duration = awsConfig.AssumeRoleDuration     // gets substituted by 15min if empty
	//	options.ExternalID = aws.String(awsConfig.ExternalID)
	//	options.SourceIdentity = aws.String(awsConfig.SourceIdentity)
	//	if len(awsConfig.SessionTags) > 0 {
	//		tags := make([]types.Tag, 0)
	//		for key, value := range awsConfig.SessionTags {
	//			tag := types.Tag{
	//				Key:   aws.String(key),
	//				Value: aws.String(value),
	//			}
	//			tags = append(tags, tag)
	//		}
	//		options.Tags = tags
	//	}
	//	options.TransitiveTagKeys = awsConfig.TransitiveSessionTags
	//})
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
