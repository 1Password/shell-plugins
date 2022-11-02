package needsauth

import "github.com/1Password/shell-plugins/sdk"

// For returns a NeedsAuthentication rule that iterates over other NeedsAuthentication rules
// until there's one that opts out of the authentication requirement.
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

// ForCommands returns a NeedsAuthentication rule to require authentication for
// certain (sub)command, e.g. ["account"] or ["account", "list"], ["account", "delete"].
func ForCommands(commands ...[]string) sdk.NeedsAuthentication {
	return func(in sdk.NeedsAuthenticationInput) bool {
		for _, command := range commands {
			if len(command) > len(in.CommandArgs) {
				continue
			}

			for i := range command {
				if command[i] != in.CommandArgs[i] {
					break
				}
				if i == len(command)-1 {
					return true
				}
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

// NotForArgs returns a NeedsAuthentication rule to not require authentication when
// certain command-line args or flags are present.
func NotForArgs(argsToSkip ...string) sdk.NeedsAuthentication {
	return func(in sdk.NeedsAuthenticationInput) bool {
		for _, commandArg := range in.CommandArgs {
			for _, ignoreArg := range argsToSkip {
				if commandArg == ignoreArg {
					return false
				}
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
