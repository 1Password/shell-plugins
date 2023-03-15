package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	confighelpers "github.com/99designs/aws-vault/v7/vault"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type assumeRoleProvider struct {
	confighelpers.AssumeRoleProvider
	stsCacheWriter
}

func (p assumeRoleProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	credentials, err := p.AssumeRoleProvider.Retrieve(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}

	err = p.stsCacheWriter.persist(credentials)
	if err != nil {
		return aws.Credentials{}, err
	}

	return credentials, nil
}

func NewAssumeRoleProvider(awsConfig *confighelpers.Config, in sdk.ProvisionInput, out *sdk.ProvisionOutput) aws.CredentialsProvider {
	roleCacheKey := getRoleCacheKey(awsConfig.RoleARN, in.ItemFields[fieldname.AccessKeyID])
	if in.Cache.Has(roleCacheKey) {
		return NewStsCacheProvider(roleCacheKey, in.Cache)
	}

	cacheWriter := NewStsCacheWriter(roleCacheKey, out.Cache)

	if awsConfig.HasMfaSerial() && awsConfig.MfaToken != "" {
		return initAssumeRoleProvider(awsConfig, getSTSClient(awsConfig.Region, NewMFASessionTokenProvider(awsConfig, in, out)), cacheWriter)
	}

	return initAssumeRoleProvider(awsConfig, getSTSClient(awsConfig.Region, NewAccessKeysProvider(in.ItemFields)), cacheWriter)
}

func initAssumeRoleProvider(awsConfig *confighelpers.Config, stsClient *sts.Client, cacheWriter stsCacheWriter) *assumeRoleProvider {
	if awsConfig.AssumeRoleDuration == 0 {
		awsConfig.AssumeRoleDuration = 900 * time.Second // default to minimum duration of 15 minutes for security
	}
	return &assumeRoleProvider{
		AssumeRoleProvider: confighelpers.AssumeRoleProvider{
			StsClient:         stsClient,
			RoleARN:           awsConfig.RoleARN,
			RoleSessionName:   awsConfig.RoleSessionName,
			ExternalID:        awsConfig.ExternalID,
			Duration:          awsConfig.AssumeRoleDuration,
			Tags:              awsConfig.SessionTags,
			TransitiveTagKeys: awsConfig.TransitiveSessionTags,
			SourceIdentity:    awsConfig.SourceIdentity,
			Mfa:               &confighelpers.Mfa{},
		},
		stsCacheWriter: cacheWriter,
	}
}

type mfaSessionTokenProvider struct {
	confighelpers.SessionTokenProvider
	stsCacheWriter
}

func (p mfaSessionTokenProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	credentials, err := p.SessionTokenProvider.Retrieve(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}

	err = p.stsCacheWriter.persist(credentials)
	if err != nil {
		return aws.Credentials{}, err
	}

	return credentials, nil
}

func NewMFASessionTokenProvider(awsConfig *confighelpers.Config, in sdk.ProvisionInput, out *sdk.ProvisionOutput) aws.CredentialsProvider {
	mfaCacheKey := getMfaCacheKey(in.ItemFields[fieldname.AccessKeyID])
	if in.Cache.Has(mfaCacheKey) {
		return NewStsCacheProvider(mfaCacheKey, in.Cache)
	}

	if awsConfig.NonChainedGetSessionTokenDuration == 0 {
		awsConfig.NonChainedGetSessionTokenDuration = 900 * time.Second // default to minimum duration of 15 minutes for security
	}

	return &mfaSessionTokenProvider{
		SessionTokenProvider: confighelpers.SessionTokenProvider{
			StsClient: getSTSClient(awsConfig.Region, NewAccessKeysProvider(in.ItemFields)),
			Duration:  awsConfig.NonChainedGetSessionTokenDuration,
			Mfa:       confighelpers.NewMfa(awsConfig),
		},
		stsCacheWriter: NewStsCacheWriter(mfaCacheKey, out.Cache),
	}
}

func NewStsCacheProvider(key string, cache sdk.CacheState) aws.CredentialsProvider {
	return stsCacheProvider{
		awsCacheKey: key,
		cache:       cache,
	}
}

type accessKeysProvider struct {
	itemFields map[sdk.FieldName]string
}

func (p accessKeysProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	secret, hasSecretKey := p.itemFields[fieldname.SecretAccessKey]
	keyId, hasKeyId := p.itemFields[fieldname.AccessKeyID]

	if !hasKeyId || !hasSecretKey {
		return aws.Credentials{}, fmt.Errorf("no long lived access key pair found. Please add your Access Key Id and Secret Access Key to your 1Password item's fields")
	}

	return aws.Credentials{
		AccessKeyID:     keyId,
		SecretAccessKey: secret,
	}, nil
}

func NewAccessKeysProvider(itemFields map[sdk.FieldName]string) aws.CredentialsProvider {
	return accessKeysProvider{itemFields: itemFields}
}

func getSTSClient(region string, credsProvider aws.CredentialsProvider) *sts.Client {
	clientConfig := aws.Config{
		Region:      region,
		Credentials: credsProvider,
	}
	return sts.NewFromConfig(clientConfig)
}
