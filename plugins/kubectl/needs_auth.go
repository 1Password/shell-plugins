package kubectl

import "github.com/1Password/shell-plugins/sdk"

var localSubcommands = map[string]bool{
	"help":             true,
	"config":           true,
	"completion":       true,
	"options":          true,
	"plugin":           true,
	"kustomize":        true,
	"kuberc":           true,
	"__complete":       true,
	"__completeNoDesc": true,
}

func needsAuth(in sdk.NeedsAuthenticationInput) bool {
	p := parseArgs(in.CommandArgs)
	sub := p.subcommand()
	switch {
	case sub == "":
		return false
	case p.boolValue("help"):
		return false
	case localSubcommands[sub]:
		return false
	case sub == "version" && p.boolValue("client"):
		return false
	}
	for _, flag := range explicitAuthFlags {
		if p.has(flag) {
			return false
		}
	}
	return true
}
