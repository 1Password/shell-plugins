package needsauth

import (
	"github.com/1Password/shell-plugins/sdk"
)

// For returns a NeedsAuthentication that opts in to the authentication requirement, unless there's
// one rule specified that opts out to the authentication requirement.
func For(rules ...sdk.NeedsAuthentication) sdk.NeedsAuthentication {
	return func(in sdk.NeedsAuthenticationInput) bool {
		for _, rule := range rules {
			if !rule(in) {
				return false
			}
		}
		return true
	}
}

// OnlyFor returns a NeedsAuthentication rule that only opts in of the authentication requirement if there's
// at least one rule specified that opts in to the authentication requirement.
func OnlyFor(rules ...sdk.NeedsAuthentication) sdk.NeedsAuthentication {
	return func(in sdk.NeedsAuthenticationInput) bool {
		for _, rule := range rules {
			if rule(in) {
				return true
			}
		}
		return false
	}
}

// ForCommand returns a NeedsAuthentication rule to require authentication for
// certain (sub)command, e.g. ["account"] or ["account", "list"].
func ForCommand(command ...string) sdk.NeedsAuthentication {
	return func(in sdk.NeedsAuthenticationInput) bool {
		if len(command) > len(in.CommandArgs) {
			return false
		}

		for i := range command {
			if command[i] != in.CommandArgs[i] {
				return false
			}
			if i == len(command)-1 {
				return true
			}
		}

		return false
	}
}

// Always returns a NeedsAuthentication rule to always require authentication.
func Always() sdk.NeedsAuthentication {
	return func(in sdk.NeedsAuthenticationInput) bool {
		return true
	}
}

// NotForExactArgs returns a NeedsAuthentication rule to opt out of authentication when
// the command-line args are an exact match with the passed in args.
func NotForExactArgs(argsToSkip ...string) sdk.NeedsAuthentication {
	return func(in sdk.NeedsAuthenticationInput) bool {
		if len(in.CommandArgs) != len(argsToSkip) {
			return true
		}

		for i, commandArg := range in.CommandArgs {
			if commandArg != argsToSkip[i] {
				return true
			}
		}

		return false
	}
}

// NotWhenContainsArgs returns a NeedsAuthentication rule to not require authentication when
// the exact sequence of argsToSkip is present somewhere in the command-line args.
func NotWhenContainsArgs(argsToSkip ...string) sdk.NeedsAuthentication {
	return func(in sdk.NeedsAuthenticationInput) bool {
		if len(argsToSkip) == 0 {
			return true
		}

		if len(argsToSkip) > len(in.CommandArgs) {
			return true
		}

		for i := range in.CommandArgs {
			if i+len(argsToSkip) > len(in.CommandArgs) {
				return true
			}

			matches := true
			for i, argsToCompare := range in.CommandArgs[i : i+len(argsToSkip)] {
				if argsToCompare != argsToSkip[i] {
					matches = false
				}
			}

			// If the argsToSkip are found in the command-line args, return that the command
			// does not not require authentication
			if matches {
				return false
			}
		}
		return true
	}
}

func NotForHelp() sdk.NeedsAuthentication {
	return NotForArgs("-h", "--help", "-help", "help")
}

func NotForVersion() sdk.NeedsAuthentication {
	return NotForArgs("-v", "--version", "-version", "version")
}

func NotForHelpOrVersion() sdk.NeedsAuthentication {
	return For(NotForHelp(), NotForVersion())
}
