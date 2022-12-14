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

const maxNameHintLength = 24

// SanitizeNameHint can be used to sanitize the name hint before passing it to the import candidate to
// improve the suggested item title.
func SanitizeNameHint(nameHint string) string {
	// Omit the name hint if it's "default", which doesn't add much value in the item name
	if nameHint == "default" {
		return ""
	}

	// Avoid name hints that are too long
	if len(nameHint) > maxNameHintLength {
		nameHint = nameHint[:maxNameHintLength-1] + "…"
	}

	return nameHint
}
