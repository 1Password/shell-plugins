package ibmcloud

import (
	"context"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

// CLIProvisioner is responsible for provisioning the IBM Cloud CLI with the appropriate credentials
type CLIProvisioner struct {
}

func (p CLIProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	// First, set the API key as an environment variable
	out.AddEnvVar("IBMCLOUD_API_KEY", in.ItemFields[fieldname.APIKey])

	// Check if we need to add a resource group flag to the login command
	targetGroup := ""
	found := false
	// Try exact match first
	if tg, ok := in.ItemFields[sdk.FieldName("resource group")]; ok && tg != "" {
		targetGroup = tg
		found = true
	}

	// Try with trimmed whitespace
	if !found {
		for fieldName, fieldValue := range in.ItemFields {
			fieldNameStr := string(fieldName)
			trimmedName := strings.TrimSpace(fieldNameStr)
			if trimmedName == "resource group" && fieldValue != "" {
				targetGroup = fieldValue
				found = true
				break
			}
		}
	}

	if found && targetGroup != "" {
		// Only modify the command line if it's an ibmcloud login command
		if len(out.CommandLine) > 1 && out.CommandLine[1] == "login" {
			// Check if the command already has a -g or --target-group flag
			hasTargetGroupFlag := false
			for _, arg := range out.CommandLine {
				if arg == "-g" || arg == "--target-group" || strings.HasPrefix(arg, "-g=") || strings.HasPrefix(arg, "--target-group=") {
					hasTargetGroupFlag = true
					break
				}
			}

			// If no resource group flag is present, add it
			if !hasTargetGroupFlag {
				newCommandLine := append([]string{}, out.CommandLine[0], out.CommandLine[1])
				newCommandLine = append(newCommandLine, "-g", targetGroup)
				newCommandLine = append(newCommandLine, out.CommandLine[2:]...)
				out.CommandLine = newCommandLine
			}
		}
	}
}

func (p CLIProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p CLIProvisioner) Description() string {
	return "Provision environment variables with IBM Cloud API Key and add resource group flag to login command if specified"
}
