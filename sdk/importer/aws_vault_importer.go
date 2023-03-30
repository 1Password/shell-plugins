package importer

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
)

func TryVault(result func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportAttempt)) sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		attempt := out.NewAttempt(sdk.ImportSource{AWSVault: true})

		result(ctx, in, attempt)
	}
}
