package fieldname

import "github.com/1Password/shell-plugins/sdk"

// Credential field names.
const (
	APIHost          = sdk.FieldName("API Host")
	APIKey           = sdk.FieldName("API Key")
	APIKeyID         = sdk.FieldName("API Key ID")
	APISecret        = sdk.FieldName("API Secret")
	AccessKeyID      = sdk.FieldName("Access Key ID")
	AccessToken      = sdk.FieldName("Access Token")
	Account          = sdk.FieldName("Account")
	AccountID        = sdk.FieldName("Account ID")
	AccountSID       = sdk.FieldName("Account SID")
	Address          = sdk.FieldName("Address")
	AppKey           = sdk.FieldName("App Key")
	AppSecret        = sdk.FieldName("App Secret")
	AppToken         = sdk.FieldName("App Token")
	AuthToken        = sdk.FieldName("Auth Token")
	Authtoken        = sdk.FieldName("Authtoken")
	Cert             = sdk.FieldName("Cert")
	Certificate      = sdk.FieldName("Certificate")
	ClientSecret     = sdk.FieldName("Client Secret")
	ClientToken      = sdk.FieldName("Client Token")
	ConnectionString = sdk.FieldName("Connection String")
	Credential       = sdk.FieldName("Credential")
	Credentials      = sdk.FieldName("Credentials")
	Database         = sdk.FieldName("Database")
	DefaultRegion    = sdk.FieldName("Default Region")
	Email            = sdk.FieldName("Email")
	Endpoint         = sdk.FieldName("Endpoint")
	Host             = sdk.FieldName("Host")
	HostAddress      = sdk.FieldName("Host Address")
	Key              = sdk.FieldName("Key")
	MFASerial        = sdk.FieldName("MFA Serial")
	Mode             = sdk.FieldName("Mode")
	Namespace        = sdk.FieldName("Namespace")
	OneTimePassword  = sdk.FieldName("One-Time Password")
	OrgID            = sdk.FieldName("Org ID")
	OrgURL           = sdk.FieldName("Org URL")
	Organization     = sdk.FieldName("Organization")
	Password         = sdk.FieldName("Password")
	Port             = sdk.FieldName("Port")
	PublicKey        = sdk.FieldName("Public Key")
	PrivateKey       = sdk.FieldName("Private Key")
	ProjectID        = sdk.FieldName("Project ID")
	Project          = sdk.FieldName("Project")
	Region           = sdk.FieldName("Region")
	Secret           = sdk.FieldName("Secret")
	SecretAccessKey  = sdk.FieldName("Secret Access Key")
	Subdomain        = sdk.FieldName("Subdomain")
	Token            = sdk.FieldName("Token")
	URL              = sdk.FieldName("URL")
	User             = sdk.FieldName("User")
	Username         = sdk.FieldName("Username")
	Website          = sdk.FieldName("Website")
)

func ListAll() []sdk.FieldName {
	return []sdk.FieldName{
		APIHost,
		APIKey,
		APIKeyID,
		APISecret,
		AccessKeyID,
		AccessToken,
		Account,
		AccountID,
		AccountSID,
		Address,
		AppKey,
		AppSecret,
		AppToken,
		AuthToken,
		Authtoken,
		Cert,
		Certificate,
		ClientSecret,
		ClientToken,
		ConnectionString,
		Credential,
		Credentials,
		Database,
		DefaultRegion,
		Endpoint,
		Host,
		HostAddress,
		Key,
		MFASerial,
		Mode,
		Namespace,
		OneTimePassword,
		OrgID,
		OrgURL,
		Organization,
		Password,
		Port,
		PublicKey,
		PrivateKey,
		ProjectID,
		Project,
		Region,
		Secret,
		SecretAccessKey,
		Token,
		URL,
		User,
		Username,
		Website,
	}
}
