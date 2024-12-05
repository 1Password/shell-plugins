package importer

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
)

// NoOp can be used as an importer stub while developing plugins.
func NoOp() sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		attempt := out.NewAttempt(SourceEnvVars(""))
		attempt.AddCandidate(sdk.ImportCandidate{})
	}
}
