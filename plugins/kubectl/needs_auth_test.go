package kubectl

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk/plugintest"
)

func TestKubectlNeedsAuth(t *testing.T) {
	yes := func(args ...string) plugintest.NeedsAuthCase {
		return plugintest.NeedsAuthCase{Args: args, ExpectedNeedsAuth: true}
	}
	no := func(args ...string) plugintest.NeedsAuthCase {
		return plugintest.NeedsAuthCase{Args: args, ExpectedNeedsAuth: false}
	}
	plugintest.TestNeedsAuth(t, needsAuth, map[string]plugintest.NeedsAuthCase{
		"yes get pods":                     yes("get", "pods"),
		"yes context get pods":             yes("--context", "prod", "get", "pods"),
		"yes short namespace get pods":     yes("-n", "kube-system", "get", "pods"),
		"yes kubeconfig apply":             yes("--kubeconfig", "f", "apply", "-f", "x"),
		"yes version":                      yes("version"),
		"yes version client false":         yes("version", "--client=false"),
		"yes version client last wins":     yes("version", "--client=true", "--client=false"),
		"yes api-resources":                yes("api-resources"),
		"yes explain":                      yes("explain", "pods"),
		"yes cluster-info":                 yes("cluster-info"),
		"yes auth can-i":                   yes("auth", "can-i", "get", "pods"),
		"yes exec after terminator":        yes("exec", "pod", "--", "df", "-h"),
		"yes exec help after terminator":   yes("exec", "pod", "--", "sh", "-c", "help"),
		"yes unknown subcommand":           yes("foo"),
		"yes help consumed as value":       yes("--namespace", "--help", "version"),
		"yes terminator consumed as value": yes("--namespace", "--", "version"),
		"yes version help false":           yes("version", "--help=false"),
		"no explicit token":                no("get", "pods", "--token=t"),
		"no explicit user":                 no("--user", "other", "get", "pods"),
		"no explicit client key":           no("get", "pods", "--client_key=k"),
		"no empty":                         no(),
		"no long help":                     no("--help"),
		"no short help":                    no("-h"),
		"no get -h":                        no("get", "-h"),
		"no help":                          no("help"),
		"no help get":                      no("help", "get"),
		"no config view":                   no("config", "view"),
		"no config use-context":            no("--context", "x", "config", "use-context", "y"),
		"no completion":                    no("completion", "zsh"),
		"no options":                       no("options"),
		"no plugin list":                   no("plugin", "list"),
		"no kustomize":                     no("kustomize", "."),
		"no kuberc":                        no("kuberc", "view"),
		"no version client":                no("version", "--client"),
		"no version client true":           no("version", "--client=true"),
		"no version client space true":     no("version", "--client", "true"),
		"no __complete":                    no("__complete", "get", ""),
		"no __completeNoDesc":              no("__completeNoDesc", "get", ""),
		"no verbosity alone":               no("-v=5"),
	})
}
