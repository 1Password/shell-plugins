package importer

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
)

func TryAll(importers ...sdk.Importer) sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		for _, imp := range importers {
			imp(ctx, in, out)
		}
	}
}
