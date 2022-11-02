package sdk

// NeedsAuthentication provides a hook to check whether authentication are required for certain command args.
type NeedsAuthentication func(in NeedsAuthenticationInput) (needsAuthentication bool)

type NeedsAuthenticationInput struct {
	CredentialType string
	CommandArgs    []string
}
