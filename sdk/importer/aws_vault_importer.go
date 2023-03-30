package importer

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
)

func TryVault(path string, result func(ctx context.Context, contents FileContents, in sdk.ImportInput, out *sdk.ImportAttempt)) sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		abspath := path
		if strings.HasPrefix(path, "~/") {
			abspath = filepath.Join(in.HomeDir, strings.TrimPrefix(path, "~/"))
		} else if strings.HasPrefix(path, "/") {
			abspath = filepath.Join(in.RootDir, path)
		}

		attempt := out.NewAttempt(sdk.ImportSource{AWSVault: true})
		contents, err := os.ReadFile(abspath)
		if os.IsNotExist(err) {
			return
		} else if err != nil {
			attempt.AddError(err)
			return
		}

		result(ctx, contents, in, attempt)
	}
}
