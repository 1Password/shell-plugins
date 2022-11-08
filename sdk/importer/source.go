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
