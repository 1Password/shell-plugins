package age

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/provision"
)

const (
	decryptShort = "-d"
	decryptLong  = "--decrypt"
	encryptShort = "-e"
	encryptLong  = "--encrypt"
)

const (
	Encrypt Operation = iota
	Decrypt
)

// Operation defines the type of action (encryption or decryption) to be performed.
type Operation int

// OperationHandler defines a function to handle specific operations.
type OperationHandler func(KeyFiles, *sdk.ProvisionOutput) (provision.ItemToFileContents, []string, string)

// String returns the string representation of an Operation.
func (op Operation) String() string {
	switch op {
	case Encrypt:
		return "encrypt"
	case Decrypt:
		return "decrypt"
	default:
		return "unknown"
	}
}

// detectOperation determines the operation (encrypt or decrypt) based on the provided command-line arguments.
// If no valid operation flags are detected, it defaults to encryption mode.
func detectOperation(args []string) Operation {
	for _, arg := range args {
		switch arg {
		case decryptShort, decryptLong:
			return Decrypt
		case encryptShort, encryptLong:
			return Encrypt
		}
	}
	return Encrypt
}
