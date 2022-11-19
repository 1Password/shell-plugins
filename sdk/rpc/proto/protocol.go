package proto

import (
	"fmt"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

const (
	Version          uint = 1
	MagicCookieKey        = "OP_PLUGIN_MAGIC_COOKIE"
	MagicCookieValue      = "ThisIsNotForSecurityPurposesButToImproveUserExperience"
)

// ExecutableID uniquely identifies an executable within a schema.Plugin by its slice index.
type ExecutableID int

func (e ExecutableID) String() string {
	return fmt.Sprintf("plugin.Executables[%d]", e)
}

// CredentialID uniquely identifies a credential within a plugin.
type CredentialID int

func (c CredentialID) String() string {
	return fmt.Sprintf("plugin.Credentials[%d]", c)
}

// ProvisionerID uniquely identifies a provisioner within a plugin.
type ProvisionerID struct {
	Credential CredentialID
	Executable *ExecutableID
}

func (p ProvisionerID) String() string {
	if p.Executable == nil {
		return fmt.Sprintf("plugin.Credentials[%d].DefaultProvisioner", p.Credential)
	}
	return fmt.Sprintf("plugin.Credentials[%d].Provisioner[%d]", p.Credential, *p.Executable)
}

// GetPluginResponse augments schema.Plugin with information about which credentials have the (optional) Importer set
// and which executables have the (optional) NeedsAuth field set. This is stored separately because these fields are
// all set to nil before sending the schema.Plugin over RPC.
type GetPluginResponse struct {
	schema.Plugin
	// CredentialHasImporter contains a true value for all credentials that have their Importer field set.
	CredentialHasImporter map[CredentialID]bool
	// ExecutableHasNeedAuth contains a true value for all executables that have their NeedsAuth field set.
	ExecutableHasNeedAuth map[ExecutableID]bool
	// CredentialIDByName maps credential names to IDs.
	CredentialIDByName map[string]CredentialID
}

// ImportCredentialRequest augments sdk.ImportInput with a CredentialID so Import() can be called over RPC.
type ImportCredentialRequest struct {
	CredentialID
	sdk.ImportInput
}

// ProvisionCredentialRequest augments sdk.ProvisionInput with a CredentialID so Provision() can be called over RPC.
type ProvisionCredentialRequest struct {
	ProvisionerID
	sdk.ProvisionInput
}

// DeprovisionCredentialRequest augments sdk.DeprovisionInput with a CredentialID so Deprovision() can be called over RPC.
type DeprovisionCredentialRequest struct {
	ProvisionerID
	sdk.DeprovisionInput
}

// ExecutableNeedsAuthRequest augments sdk.NeedsAuthenticationInput with the ID of an executable so NeedsAuth() can be
// called over RPC. ExecutableID resembles the slice index of the executable in schema.Plugin.
type ExecutableNeedsAuthRequest struct {
	ExecutableID
	sdk.NeedsAuthenticationInput
}
