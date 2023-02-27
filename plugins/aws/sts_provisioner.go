package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"strings"
)

const (
	mfaCacheKey               = "sts-mfa"
	assumeRoleCacheKey        = "sts-assume-role"
	assumeRoleWithMFACacheKey = "sts-assume-role-mfa"
)

type STSProvisioner struct {
	TOTPCode        string
	MFASerial       string
	ProfileWithRole *ini.Section
}

func (p STSProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
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

	if p.ProfileWithRole != nil && p.MFASerial != "" && p.TOTPCode != "" {
		if ok := p.tryUsingCachedCredentials(assumeRoleWithMFACacheKey, in, out); ok {
			return
		}

		roleArn, err := p.ProfileWithRole.GetKey("role_arn")
		if err != nil {
			return
		}

		input := &sts.AssumeRoleInput{
			DurationSeconds: aws.Int32(900), // minimum expiration time - 15 minutes
			SerialNumber:    aws.String(p.MFASerial),
			TokenCode:       aws.String(p.TOTPCode),
			RoleArn:         aws.String(roleArn.Value()),
			RoleSessionName: aws.String("1password-shell-plugin"),
		}

		resp, err := stsProvider.AssumeRole(ctx, input)
		if err != nil {
			out.AddError(err)
			return
		}
		provisionConfigFile(p.ProfileWithRole, in.TempDir, resp.Credentials, out)

		err = out.Cache.Put(assumeRoleWithMFACacheKey, resp.Credentials, *resp.Credentials.Expiration)
		if err != nil {
			out.AddError(fmt.Errorf("failed to serialize aws sts credentials: %w", err))
		}
	} else if p.ProfileWithRole == nil {
		if useCache := p.tryUsingCachedCredentials(mfaCacheKey, in, out); useCache {
			return
		}
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
		out.AddEnvVar("AWS_ACCESS_KEY_ID", *resp.Credentials.AccessKeyId)
		out.AddEnvVar("AWS_SECRET_ACCESS_KEY", *resp.Credentials.SecretAccessKey)
		out.AddEnvVar("AWS_SESSION_TOKEN", *resp.Credentials.SessionToken)
		err = out.Cache.Put(mfaCacheKey, *resp.Credentials, *resp.Credentials.Expiration)
		if err != nil {
			out.AddError(fmt.Errorf("failed to serialize aws sts credentials: %w", err))
		}
	} else {
		if ok := p.tryUsingCachedCredentials(assumeRoleCacheKey, in, out); ok {
			return
		}
		roleArn, err := p.ProfileWithRole.GetKey("role_arn")
		if err != nil {
			return
		}

		input := &sts.AssumeRoleInput{
			DurationSeconds: aws.Int32(900), // minimum expiration time - 15 minutes
			RoleArn:         aws.String(roleArn.Value()),
			RoleSessionName: aws.String("1password-shell-plugin"),
		}

		resp, err := stsProvider.AssumeRole(ctx, input)
		if err != nil {
			out.AddError(err)
			return
		}
		provisionConfigFile(p.ProfileWithRole, in.TempDir, resp.Credentials, out)
	}

	out.AddEnvVar("AWS_DEFAULT_REGION", region)
}

func (p STSProvisioner) tryUsingCachedCredentials(cacheKey string, in sdk.ProvisionInput, out *sdk.ProvisionOutput) bool {
	isAssumeRoleWorkflow := strings.Contains(cacheKey, assumeRoleCacheKey)
	var cached types.Credentials
	if ok := in.Cache.Get(cacheKey, &cached); ok {
		if isAssumeRoleWorkflow {
			provisionConfigFile(p.ProfileWithRole, in.TempDir, &cached, out)
		} else {
			out.AddEnvVar("AWS_ACCESS_KEY_ID", *cached.AccessKeyId)
			out.AddEnvVar("AWS_SECRET_ACCESS_KEY", *cached.SecretAccessKey)
			out.AddEnvVar("AWS_SESSION_TOKEN", *cached.SessionToken)
		}

		if region, ok := in.ItemFields[fieldname.DefaultRegion]; ok {
			out.AddEnvVar("AWS_DEFAULT_REGION", region)
		}

		return true
	}
	return false
}

func provisionConfigFile(profileWithRole *ini.Section, tempDir string, tempCredentials *types.Credentials, out *sdk.ProvisionOutput) {
	provisionedConfigFile := ini.Empty()
	provisionedSection, err := provisionedConfigFile.NewSection(profileWithRole.Name())
	if err != nil {
		out.AddError(err)
		return
	}
	for _, key := range profileWithRole.Keys() {
		if key.Name() != "role_arn" && key.Name() != "credential_process" {
			_, err = provisionedSection.NewKey(key.Name(), key.Value())
			if err != nil {
				out.AddError(err)
				return
			}
		}
	}
	credentialOutput := struct {
		Version         int
		AccessKeyId     string
		SecretAccessKey string
		SessionToken    string
	}{
		1,
		*tempCredentials.AccessKeyId,
		*tempCredentials.SecretAccessKey,
		*tempCredentials.SessionToken,
	}

	jsonCredentials, err := json.Marshal(credentialOutput)
	if err != nil {
		out.AddError(err)
		return
	}

	helperFIFOPath := filepath.Join(tempDir, "helperFIFO")
	out.AddSecretFile(helperFIFOPath, jsonCredentials)

	provisioningProcess := fmt.Sprintf("cat %s", helperFIFOPath)

	_, err = provisionedSection.NewKey("credential_process", provisioningProcess)
	if err != nil {
		out.AddError(err)
		return
	}

	var buf bytes.Buffer
	_, err = provisionedConfigFile.WriteTo(&buf)
	if err != nil {
		return
	}

	configPath := filepath.Join(tempDir, "config")
	err = os.WriteFile(configPath, buf.Bytes(), 0777)
	if err != nil {
		out.AddError(err)
		return
	}
	out.AddEnvVar("AWS_CONFIG_FILE", configPath)
}

func (p STSProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	configPath := filepath.Join(in.TempDir, "config")
	err := os.Remove(configPath)
	if err != nil {
		out.Diagnostics.Errors = append(out.Diagnostics.Errors, sdk.Error{Message: err.Error()})
		return
	}
}

func (p STSProvisioner) Description() string {
	return "Provision environment variables with temporary STS credentials AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN"
}
