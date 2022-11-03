package importer

import (
	"context"
	"os"

	"github.com/1Password/shell-plugins/sdk"
)

func TryAllEnvVars(fieldName string, possibleEnvVarNames ...string) sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		matches := ScanEnvironment(fieldName, possibleEnvVarNames...)

		for _, match := range matches {
			out.AddCandidate(sdk.ImportCandidate{
				Fields: []sdk.ImportCandidateField{match},
			})
		}
	}
}

func TryEnvVarPairVariations(pairPossibilities map[string][]string) sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		var detectedFields []sdk.ImportCandidateField
		for fieldName, possibleEnvVarNames := range pairPossibilities {
			matches := ScanEnvironment(fieldName, possibleEnvVarNames...)
			if len(matches) > 0 {
				detectedFields = append(detectedFields, matches[0])
			}
		}
		if len(detectedFields) != 0 {
			out.AddCandidate(sdk.ImportCandidate{
				Fields: detectedFields,
			})
		}
	}
}

func TryEnvVarPair(pairPossibilities map[string]string) sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		var detectedFields []sdk.ImportCandidateField
		for fieldName, possibleEnvVarName := range pairPossibilities {
			matches := ScanEnvironment(fieldName, possibleEnvVarName)
			if len(matches) > 0 {
				detectedFields = append(detectedFields, matches[0])
			}
		}

		if len(detectedFields) != 0 {
			out.AddCandidate(sdk.ImportCandidate{
				Fields: detectedFields,
			})
		}
	}
}

func ScanEnvironment(field string, possibleEnvVarNames ...string) []sdk.ImportCandidateField {
	var matches []sdk.ImportCandidateField
	for _, envVarName := range possibleEnvVarNames {
		if value := os.Getenv(envVarName); value != "" {
			matches = append(matches, sdk.ImportCandidateField{
				Field:  field,
				Value:  value,
				Source: SourceEnvName(envVarName),
			})
		}
	}
	return matches
}
