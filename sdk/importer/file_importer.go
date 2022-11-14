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

func TryFile(path string, result func(ctx context.Context, contents FileContents, in sdk.ImportInput, out *sdk.ImportOutput)) sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		if strings.HasPrefix(path, "~/") {
			path = filepath.Join(in.HomeDir, strings.TrimPrefix(path, "~/"))
		}

		contents, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			return
		} else if err != nil {
			out.AddError(err)
			return
		}

		result(ctx, contents, in, out)
	}
}

type FileContents []byte

func (fc FileContents) ToJSON(result interface{}) error {
	err := json.Unmarshal(fc, result)
	if err != nil {
		return err
	}

	return nil
}

func (fc FileContents) ToYAML(result interface{}) error {
	err := yaml.Unmarshal(fc, result)
	if err != nil {
		return err
	}

	return nil
}

func (fc FileContents) ToTOML(result interface{}) error {
	err := toml.Unmarshal(fc, result)
	if err != nil {
		return err
	}

	return nil
}

func (fc FileContents) ToXML(result interface{}) error {
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
