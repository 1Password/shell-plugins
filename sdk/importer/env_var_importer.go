package importer

import (
	"context"
	"os"

	"github.com/1Password/shell-plugins/sdk"
)

// TryAllEnvVars tries the specified environment variables one by one and adds import candidates with
// the specified field for each environment variable that is set.
func TryAllEnvVars(fieldName sdk.FieldName, possibleEnvVarNames ...string) sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		for _, envVarName := range possibleEnvVarNames {
			attempt := out.NewAttempt(SourceEnvVars(envVarName))

			if value := os.Getenv(envVarName); value != "" {
				attempt.AddCandidate(sdk.ImportCandidate{
					Fields: map[sdk.FieldName]string{
						fieldName: value,
					},
				})
			}
		}
	}
}

// TryEnvVarPair tries the specified environment variables and adds an import candidate if at least
// one environment variable is set.
func TryEnvVarPair(pairPossibilities map[sdk.FieldName]string) sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		var envVarNames []string
		candidateFields := make(map[sdk.FieldName]string)

		for fieldName, possibleEnvVarName := range pairPossibilities {
			if value := os.Getenv(possibleEnvVarName); value != "" {
				candidateFields[fieldName] = value
			}

			envVarNames = append(envVarNames, possibleEnvVarName)
		}

		attempt := out.NewAttempt(SourceEnvVars(envVarNames...))
		if len(candidateFields) > 0 {
			attempt.AddCandidate(sdk.ImportCandidate{
				Fields: candidateFields,
			})
		}
	}
}
