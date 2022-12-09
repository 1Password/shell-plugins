package fieldname

import "github.com/1Password/shell-plugins/sdk"

// Credential field names.
const (
	APIHost         = sdk.FieldName("API Host")
	APIKey          = sdk.FieldName("API Key")
	APIKeyID        = sdk.FieldName("API Key ID")
	APISecret       = sdk.FieldName("API Secret")
	AccessKeyID     = sdk.FieldName("Access Key ID")
	Account         = sdk.FieldName("Account")
	AccountID       = sdk.FieldName("Account ID")
	AccountSID      = sdk.FieldName("Account SID")
	Address         = sdk.FieldName("Address")
	AppKey          = sdk.FieldName("App Key")
	AppSecret       = sdk.FieldName("App Secret")
	AppToken        = sdk.FieldName("App Token")
	AuthToken       = sdk.FieldName("Auth Token")
	Cert            = sdk.FieldName("Cert")
	Certificate     = sdk.FieldName("Certificate")
	Credential      = sdk.FieldName("Credential")
	Credentials     = sdk.FieldName("Credentials")
	Database        = sdk.FieldName("Database")
	DefaultRegion   = sdk.FieldName("Default Region")
	Host            = sdk.FieldName("Host")
	HostAddress     = sdk.FieldName("Host Address")
	Key             = sdk.FieldName("Key")
	MFASerial       = sdk.FieldName("MFA Serial")
	Mode            = sdk.FieldName("Mode")
	Namespace       = sdk.FieldName("Namespace")
	OneTimePassword = sdk.FieldName("One-Time Password")
	OrgURL          = sdk.FieldName("Org URL")
	Organization    = sdk.FieldName("Organization")
	Password        = sdk.FieldName("Password")
	Port            = sdk.FieldName("Port")
	PrivateKey      = sdk.FieldName("Private Key")
	Region          = sdk.FieldName("Region")
	Secret          = sdk.FieldName("Secret")
	SecretAccessKey = sdk.FieldName("Secret Access Key")
	Token           = sdk.FieldName("Token")
	URL             = sdk.FieldName("URL")
	User            = sdk.FieldName("User")
	Username        = sdk.FieldName("Username")
	Website         = sdk.FieldName("Website")
	Role            = sdk.FieldName("Role")
)

func ListAll() []sdk.FieldName {
	return []sdk.FieldName{
		APIHost,
		APIKey,
		APIKeyID,
		APISecret,
		AccessKeyID,
		Account,
		AccountID,
		AccountSID,
		Address,
		AppKey,
		AppSecret,
		AppToken,
		AuthToken,
		Cert,
		Certificate,
		Credential,
		Credentials,
		Database,
		DefaultRegion,
		Host,
		HostAddress,
		Key,
		MFASerial,
		Mode,
		Namespace,
		OneTimePassword,
		OrgURL,
		Organization,
		Password,
		Port,
		PrivateKey,
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
