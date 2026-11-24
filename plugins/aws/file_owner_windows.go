//go:build windows

package aws

import "os"

// ownedByCurrentUser always returns true on Windows, where file ownership is not expressed as a
// numeric user ID. Shell plugins do not run on Windows; this only keeps the package compiling.
func ownedByCurrentUser(os.FileInfo) bool {
	return true
}
