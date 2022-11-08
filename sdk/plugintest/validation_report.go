package plugintest

import (
	"github.com/fatih/color"

	"github.com/1Password/shell-plugins/sdk/schema"
)

func MakeValidationReport(plugin schema.Plugin) string {
	color.Cyan("Validation report")
	return ""
}
