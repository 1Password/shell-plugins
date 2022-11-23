package plugintest

import (
	"context"
	"fmt"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/stretchr/testify/assert"
)

// TestProvisioner will invoke the specified provisioner with the item fields specified in each test case, comparing
// the provisioner output with the specified expected output.
func TestProvisioner(t *testing.T, provisioner sdk.Provisioner, cases map[string]ProvisionCase) {
	t.Helper()

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if c.ExpectedOutput.Environment == nil {
				c.ExpectedOutput.Environment = make(map[string]string)
			}

			if c.ExpectedOutput.Files == nil {
				c.ExpectedOutput.Files = make(map[string]sdk.OutputFile)
			}

			ctx := context.Background()

			in := sdk.ProvisionInput{
				ItemFields: c.ItemFields,
				HomeDir:    "~",
				TempDir:    "/tmp",
			}

			out := sdk.ProvisionOutput{
				Environment: make(map[string]string),
				Files:       make(map[string]sdk.OutputFile),
				CommandLine: c.CommandLine,
			}

			provisioner.Provision(ctx, in, &out)

			description := fmt.Sprintf("Provision: %s", name)
			assert.Equal(t, c.ExpectedOutput, out, description)
		})
	}
}

type ProvisionCase struct {
	// ItemFields can be used to populate the item fields to pass to the provisioner.
	ItemFields map[string]string

	// CommandLine can be used to populate the command line to pass to the provisioner.
	CommandLine []string

	// ExpectedOutput can be used to set the exact expected provision output, which contains the
	// environment, files, and command line.
	ExpectedOutput sdk.ProvisionOutput
}
