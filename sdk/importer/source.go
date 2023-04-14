package importer

import (
	"github.com/1Password/shell-plugins/sdk"
)

func SourceEnvVars(envVars ...string) sdk.ImportSource {
	return sdk.ImportSource{Env: envVars}
}

func SourceEnvName(envVarName string) sdk.ImportSource {
	return sdk.ImportSource{Env: []string{envVarName}}
}

func SourceFile(filename string) sdk.ImportSource {
	return sdk.ImportSource{Files: []string{filename}}
}

func SourceOther(sourceType string, sourceValue string) sdk.ImportSource {
	return sdk.ImportSource{Other: sdk.CustomSource{Type: sourceType, Value: []string{sourceValue}}}
}
