package pypi

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TryPyPIRCFile() sdk.Importer {
	return importer.TryFile("~/.pypirc", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		cfg, err := contents.ToINI()
		if err != nil {
			out.AddError(err)
			return
		}

		// Try [pypi] section first, then [server-login]
		for _, section := range []string{"pypi", "server-login"} {
			if !cfg.HasSection(section) {
				continue
			}
			s := cfg.Section(section)

			password := s.Key("password").String()
			if password != "" {
				out.AddCandidate(sdk.ImportCandidate{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: password,
					},
				})
				return
			}
		}
	})
}
