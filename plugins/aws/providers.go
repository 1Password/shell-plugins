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
	roleCacheKey := getRoleCacheKey(awsConfig.RoleARN)
	if in.Cache.Has(roleCacheKey) {
		return NewStsCacheProvider(roleCacheKey, in.Cache)
	}

	if awsConfig.HasMfaSerial() && awsConfig.MfaToken != "" {
		return &assumeRoleProvider{
			AssumeRoleProvider: confighelpers.AssumeRoleProvider{
				StsClient:         getSTSClient(awsConfig.Region, NewMFASessionTokenProvider(awsConfig, in, out)),
				RoleARN:           awsConfig.RoleARN,
				RoleSessionName:   awsConfig.RoleSessionName,
				ExternalID:        awsConfig.ExternalID,
				Duration:          900 * time.Second, // minimum duration of 15 minutes
				Tags:              awsConfig.SessionTags,
				TransitiveTagKeys: awsConfig.TransitiveSessionTags,
				SourceIdentity:    awsConfig.SourceIdentity,
				Mfa:               &confighelpers.Mfa{},
			},
			stsCacheWriter: NewStsCacheWriter(roleCacheKey, out.Cache),
		}
	}

	return &assumeRoleProvider{
		AssumeRoleProvider: confighelpers.AssumeRoleProvider{
			StsClient:         getSTSClient(awsConfig.Region, NewMasterCredentialsProvider(in.ItemFields)),
			RoleARN:           awsConfig.RoleARN,
			RoleSessionName:   awsConfig.RoleSessionName,
			ExternalID:        awsConfig.ExternalID,
			Duration:          900 * time.Second, // minimum duration of 15 minutes
			Tags:              awsConfig.SessionTags,
			TransitiveTagKeys: awsConfig.TransitiveSessionTags,
			SourceIdentity:    awsConfig.SourceIdentity,
			Mfa:               &confighelpers.Mfa{},
		},
		stsCacheWriter: NewStsCacheWriter(roleCacheKey, out.Cache),
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
	if in.Cache.Has(mfaCacheKey) {
		return NewStsCacheProvider(mfaCacheKey, in.Cache)
	}

	return &mfaSessionTokenProvider{
		SessionTokenProvider: confighelpers.SessionTokenProvider{
			StsClient: getSTSClient(awsConfig.Region, NewMasterCredentialsProvider(in.ItemFields)),
			Duration:  900 * time.Second, // minimum duration of 15 minutes,
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

type masterAwsCredentialsProvider struct {
	itemFields map[sdk.FieldName]string
}

func (p masterAwsCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
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

func NewMasterCredentialsProvider(itemFields map[sdk.FieldName]string) aws.CredentialsProvider {
	return masterAwsCredentialsProvider{itemFields: itemFields}
}

func getSTSClient(region string, credsProvider aws.CredentialsProvider) *sts.Client {
	clientConfig := aws.Config{
		Region:      region,
		Credentials: credsProvider,
	}
	return sts.NewFromConfig(clientConfig)
}
