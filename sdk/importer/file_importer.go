package importer

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/BurntSushi/toml"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

func TryFile(path string, result func(ctx context.Context, contents FileContents, in sdk.ImportInput, out *sdk.ImportAttempt)) sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		abspath := path
		if strings.HasPrefix(path, "~/") {
			abspath = filepath.Join(in.HomeDir, strings.TrimPrefix(path, "~/"))
		} else if strings.HasPrefix(path, "/") {
			abspath = filepath.Join(in.RootDir, path)
		}

		attempt := out.NewAttempt(SourceFile(path))
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

type FileContents []byte

func (fc FileContents) ToString() string {
	return string(fc)
}

func (fc FileContents) ToJSON(result any) error {
	err := json.Unmarshal(fc, result)
	if err != nil {
		return err
	}

	return nil
}

func (fc FileContents) ToYAML(result any) error {
	err := yaml.Unmarshal(fc, result)
	if err != nil {
		return err
	}

	return nil
}

func (fc FileContents) ToTOML(result any) error {
	err := toml.Unmarshal(fc, result)
	if err != nil {
		return err
	}

	return nil
}

func (fc FileContents) ToXML(result any) error {
	err := xml.Unmarshal(fc, result)
	if err != nil {
		return err
	}

	return nil
}

func (fc FileContents) ToINI() (*ini.File, error) {
	result, err := ini.Load([]byte(fc))
	if err != nil {
		return nil, err
	}

	return result, nil
}
