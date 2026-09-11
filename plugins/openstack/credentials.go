package openstack

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

const (
	fieldUserDomainName             = sdk.FieldName("User Domain Name")
	fieldProjectDomainName          = sdk.FieldName("Project Domain Name")
	fieldProjectDomainID            = sdk.FieldName("Project Domain ID")
	fieldInterface                  = sdk.FieldName("Interface")
	fieldIdentityAPIVersion         = sdk.FieldName("Identity API Version")
	fieldCloud                      = sdk.FieldName("Cloud")
	fieldAuthType                   = sdk.FieldName("Auth Type")
	fieldApplicationCredentialID    = sdk.FieldName("Application Credential ID")
	fieldApplicationCredentialSecret = sdk.FieldName("Application Credential Secret")
)

func Credentials() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AccessKey,
		DocsURL:       sdk.URL("https://docs.openstack.org/python-openstackclient/latest/configuration/index.html"),
		ManagementURL: sdk.URL("https://docs.openstack.org/keystone/latest/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.URL,
				MarkdownDescription: "The Keystone identity service endpoint URL (OS_AUTH_URL).",
			},
			// Password-based auth fields
			{
				Name:                fieldname.Username,
				MarkdownDescription: "The username used to authenticate to OpenStack (OS_USERNAME). Required for password-based auth.",
				Optional:            true,
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "The password used to authenticate to OpenStack (OS_PASSWORD). Required for password-based auth.",
				Secret:              true,
				Optional:            true,
			},
			// Application credential auth fields
			{
				Name:                fieldApplicationCredentialID,
				MarkdownDescription: "The application credential ID (OS_APPLICATION_CREDENTIAL_ID). Required for v3applicationcredential auth.",
				Optional:            true,
			},
			{
				Name:                fieldApplicationCredentialSecret,
				MarkdownDescription: "The application credential secret (OS_APPLICATION_CREDENTIAL_SECRET). Required for v3applicationcredential auth.",
				Secret:              true,
				Optional:            true,
			},
			// Scoping fields
			{
				Name:                fieldname.Project,
				AlternativeNames:    []string{"Project Name"},
				MarkdownDescription: "The project name to scope the OpenStack session to (OS_PROJECT_NAME).",
				Optional:            true,
			},
			{
				Name:                fieldname.ProjectID,
				MarkdownDescription: "The project ID to scope the OpenStack session to (OS_PROJECT_ID).",
				Optional:            true,
			},
			{
				Name:                fieldname.Region,
				AlternativeNames:    []string{"Region Name"},
				MarkdownDescription: "The region of the OpenStack endpoint to use (OS_REGION_NAME).",
				Optional:            true,
			},
			{
				Name:                fieldUserDomainName,
				MarkdownDescription: "The domain name of the user (OS_USER_DOMAIN_NAME).",
				Optional:            true,
			},
			{
				Name:                fieldProjectDomainName,
				MarkdownDescription: "The domain name of the project (OS_PROJECT_DOMAIN_NAME).",
				Optional:            true,
			},
			{
				Name:                fieldProjectDomainID,
				MarkdownDescription: "The domain ID of the project (OS_PROJECT_DOMAIN_ID).",
				Optional:            true,
			},
			// Connection fields
			{
				Name:                fieldInterface,
				MarkdownDescription: "The endpoint interface to use (OS_INTERFACE). Defaults to \"public\".",
				Optional:            true,
			},
			{
				Name:                fieldIdentityAPIVersion,
				MarkdownDescription: "The Keystone API version to use (OS_IDENTITY_API_VERSION). Defaults to \"3\".",
				Optional:            true,
			},
			{
				Name:                fieldAuthType,
				MarkdownDescription: "The authentication type (OS_AUTH_TYPE). Auto-detected: \"password\" or \"v3applicationcredential\".",
				Optional:            true,
			},
			{
				Name:                fieldCloud,
				MarkdownDescription: "The cloud name from clouds.yaml to use (OS_CLOUD). Useful for distinguishing multiple environments.",
				Optional:            true,
			},
		},
		DefaultProvisioner: envVarsWithDefaults(defaultEnvVarMapping, map[string]string{
			"OS_INTERFACE":            "public",
			"OS_IDENTITY_API_VERSION": "3",
		}),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryOpenStackCloudRCFromCWD(),
			TryOpenStackCloudsYAMLFromCWD(),
			TryOpenStackCloudRC("~/openrc.sh"),
			TryOpenStackCloudRC("~/.config/openstack/openrc.sh"),
			TryOpenStackCloudsYAMLFromEnvVar(),
			TryOpenStackCloudsAndSecureYAML(
				"~/.config/openstack/clouds.yaml",
				"~/.config/openstack/secure.yaml",
			),
			TryOpenStackCloudsAndSecureYAML(
				"/etc/openstack/clouds.yaml",
				"/etc/openstack/secure.yaml",
			),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"OS_AUTH_URL":                    fieldname.URL,
	"OS_USERNAME":                    fieldname.Username,
	"OS_PASSWORD":                    fieldname.Password,
	"OS_APPLICATION_CREDENTIAL_ID":   fieldApplicationCredentialID,
	"OS_APPLICATION_CREDENTIAL_SECRET": fieldApplicationCredentialSecret,
	"OS_PROJECT_NAME":                fieldname.Project,
	"OS_PROJECT_ID":                  fieldname.ProjectID,
	"OS_REGION_NAME":                 fieldname.Region,
	"OS_USER_DOMAIN_NAME":            fieldUserDomainName,
	"OS_PROJECT_DOMAIN_NAME":         fieldProjectDomainName,
	"OS_PROJECT_DOMAIN_ID":           fieldProjectDomainID,
	"OS_INTERFACE":                   fieldInterface,
	"OS_IDENTITY_API_VERSION":        fieldIdentityAPIVersion,
	"OS_AUTH_TYPE":                   fieldAuthType,
	"OS_CLOUD":                       fieldCloud,
}

// envVarsWithDefaults wraps provision.EnvVars and injects default env var values
// for any fields not present in the 1Password item.
func envVarsWithDefaults(schema map[string]sdk.FieldName, defaults map[string]string) sdk.Provisioner {
	return provisionerWithDefaults{
		base:     provision.EnvVars(schema),
		defaults: defaults,
	}
}

type provisionerWithDefaults struct {
	base     sdk.Provisioner
	defaults map[string]string
}

func (p provisionerWithDefaults) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	p.base.Provision(ctx, in, out)

	// Apply static defaults for fields not set in the item.
	for envVar, defaultVal := range p.defaults {
		if _, set := out.Environment[envVar]; !set {
			out.AddEnvVar(envVar, defaultVal)
		}
	}

	// Auto-detect OS_AUTH_TYPE if not explicitly set in the item.
	if _, set := out.Environment["OS_AUTH_TYPE"]; !set {
		if _, hasAppCred := out.Environment["OS_APPLICATION_CREDENTIAL_ID"]; hasAppCred {
			out.AddEnvVar("OS_AUTH_TYPE", "v3applicationcredential")
		} else {
			out.AddEnvVar("OS_AUTH_TYPE", "password")
		}
	}

	// When using application credentials, remove password-based fields to avoid
	// confusing the OpenStack CLI with conflicting auth parameters.
	if out.Environment["OS_AUTH_TYPE"] == "v3applicationcredential" {
		delete(out.Environment, "OS_USERNAME")
		delete(out.Environment, "OS_PASSWORD")
	}
}

func (p provisionerWithDefaults) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	p.base.Deprovision(ctx, in, out)
}

func (p provisionerWithDefaults) Description() string {
	return p.base.Description()
}

