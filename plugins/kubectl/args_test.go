package kubectl

import (
	"reflect"
	"testing"
)

func sameStrings(a, b []string) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	return reflect.DeepEqual(a, b)
}

func TestParseArgsTerminatorAndConsumedValues(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		positionals []string
		rest        []string
		values      map[string]string
		absent      []string
		subcommand  string
	}{
		{
			name:        "terminator consumed as namespace value",
			args:        []string{"--namespace", "--", "version", "--server=https://x"},
			positionals: []string{"version"},
			values:      map[string]string{"namespace": "--", "server": "https://x"},
			subcommand:  "version",
		},
		{
			name:        "flag-looking token consumed as namespace value",
			args:        []string{"--namespace", "--server=https://x", "version"},
			positionals: []string{"version"},
			values:      map[string]string{"namespace": "--server=https://x"},
			absent:      []string{"server"},
			subcommand:  "version",
		},
		{
			name:        "terminator at option boundary",
			args:        []string{"exec", "pod", "--", "df", "-h"},
			positionals: []string{"exec", "pod"},
			rest:        []string{"df", "-h"},
			absent:      []string{"help", "h"},
			subcommand:  "exec",
		},
		{
			name:        "flags after terminator are not parsed",
			args:        []string{"exec", "--context", "a", "pod", "--", "kubectl", "--context", "b"},
			positionals: []string{"exec", "pod"},
			rest:        []string{"kubectl", "--context", "b"},
			values:      map[string]string{"context": "a"},
			subcommand:  "exec",
		},
		{
			name:        "global flags before subcommand are skipped",
			args:        []string{"--context", "x", "-n", "ns", "config", "view"},
			positionals: []string{"config", "view"},
			subcommand:  "config",
			values:      map[string]string{"context": "x", "namespace": "ns"},
		},
		{
			name:        "unknown long flag never consumes next token",
			args:        []string{"--frobnicate", "get", "pods"},
			positionals: []string{"get", "pods"},
			values:      map[string]string{"frobnicate": "true"},
			subcommand:  "get",
		},
		{
			name:        "single dash is a positional",
			args:        []string{"apply", "-f", "-"},
			positionals: []string{"apply", "-"},
			values:      map[string]string{"f": "true"},
			subcommand:  "apply",
		},
		{
			name:       "no args",
			args:       nil,
			subcommand: "",
		},
		{
			name:       "only terminator",
			args:       []string{"--"},
			subcommand: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parseArgs(tt.args)
			if !sameStrings(p.positionals, tt.positionals) {
				t.Errorf("positionals = %q, want %q", p.positionals, tt.positionals)
			}
			if !sameStrings(p.rest, tt.rest) {
				t.Errorf("rest = %q, want %q", p.rest, tt.rest)
			}
			for name, want := range tt.values {
				if got, ok := p.value(name); !ok || got != want {
					t.Errorf("value(%q) = %q, %v; want %q, true", name, got, ok, want)
				}
			}
			for _, name := range tt.absent {
				if p.has(name) {
					t.Errorf("has(%q) = true, want false", name)
				}
			}
			if got := p.subcommand(); got != tt.subcommand {
				t.Errorf("subcommand() = %q, want %q", got, tt.subcommand)
			}
		})
	}
}

func TestParseArgsFlagForms(t *testing.T) {
	tests := []struct {
		name string
		args []string
		flag string
		want string
	}{
		{"last long value wins", []string{"--context", "a", "--context=b"}, "context", "b"},
		{"last long value wins reversed", []string{"--context=b", "--context", "a"}, "context", "a"},
		{"short attached", []string{"-sX"}, "server", "X"},
		{"short separate", []string{"-s", "X"}, "server", "X"},
		{"short equals", []string{"-s=X"}, "server", "X"},
		{"long separate", []string{"--server", "X"}, "server", "X"},
		{"long equals", []string{"--server=X"}, "server", "X"},
		{"underscore normalised with equals", []string{"--client_certificate=p"}, "client-certificate", "p"},
		{"underscore normalised separate", []string{"--client_certificate", "p"}, "client-certificate", "p"},
		{"namespace short", []string{"-n", "ns"}, "namespace", "ns"},
		{"namespace short attached", []string{"-nns"}, "namespace", "ns"},
		{"verbosity short equals", []string{"-v=5"}, "v", "5"},
		{"verbosity short attached", []string{"-v5"}, "v", "5"},
		{"verbosity short separate", []string{"-v", "5"}, "v", "5"},
		{"verbosity long", []string{"--v=5"}, "v", "5"},
		{"equals value containing equals", []string{"--server=https://x/?a=b"}, "server", "https://x/?a=b"},
		{"as-user-extra repeated keeps last", []string{"--as-user-extra", "a=1", "--as-user-extra=b=2"}, "as-user-extra", "b=2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parseArgs(tt.args)
			if got, ok := p.value(tt.flag); !ok || got != tt.want {
				t.Errorf("value(%q) = %q, %v; want %q, true", tt.flag, got, ok, tt.want)
			}
			if len(p.positionals) != 0 {
				t.Errorf("positionals = %q, want none", p.positionals)
			}
		})
	}
}

func TestParseArgsRepeatedValuesInOrder(t *testing.T) {
	p := parseArgs([]string{"--context", "a", "--context=b", "-s", "X"})
	if got, want := p.flags["context"], []string{"a", "b"}; !reflect.DeepEqual(got, want) {
		t.Errorf("flags[context] = %q, want %q", got, want)
	}
}

func TestParseArgsMissingValue(t *testing.T) {
	for _, args := range [][]string{{"--namespace"}, {"get", "-n"}, {"get", "--server"}} {
		p := parseArgs(args)
		for _, name := range []string{"namespace", "server"} {
			if !p.has(name) {
				continue
			}
			if got, ok := p.value(name); ok || got != "" {
				t.Errorf("%q: value(%q) = %q, %v; want \"\", false", args, name, got, ok)
			}
		}
	}
	p := parseArgs([]string{"--namespace"})
	if !p.has("namespace") {
		t.Errorf("has(namespace) = false, want true for trailing value flag")
	}
	p = parseArgs([]string{"--server="})
	if !p.has("server") {
		t.Errorf("has(server) = false, want true for --server=")
	}
}

func TestBoolValue(t *testing.T) {
	tests := []struct {
		name string
		args []string
		flag string
		want bool
	}{
		{"client bare", []string{"version", "--client"}, "client", true},
		{"client true", []string{"version", "--client=true"}, "client", true},
		{"client then true positional", []string{"version", "--client", "true"}, "client", true},
		{"client false", []string{"version", "--client=false"}, "client", false},
		{"client last wins false", []string{"version", "--client=true", "--client=false"}, "client", false},
		{"client last wins true", []string{"version", "--client=false", "--client"}, "client", true},
		{"client absent", []string{"version"}, "client", false},
		{"short help", []string{"-h"}, "help", true},
		{"long help", []string{"--help"}, "help", true},
		{"short help false", []string{"-h=false"}, "help", false},
		{"long help false", []string{"--help=false"}, "help", false},
		{"help consumed as namespace value", []string{"--namespace", "--help"}, "help", false},
		{"help after terminator", []string{"exec", "pod", "--", "-h"}, "help", false},
		{"cluster help false", []string{"-hh=false"}, "help", false},
		{"cluster help then namespace", []string{"-hnns"}, "help", true},
		{"insecure bare", []string{"--insecure-skip-tls-verify", "get"}, "insecure-skip-tls-verify", true},
		{"insecure underscore", []string{"--insecure_skip_tls_verify=true"}, "insecure-skip-tls-verify", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseArgs(tt.args).boolValue(tt.flag); got != tt.want {
				t.Errorf("boolValue(%q) on %q = %v, want %v", tt.flag, tt.args, got, tt.want)
			}
		})
	}

	p := parseArgs([]string{"version", "--client", "true"})
	if !sameStrings(p.positionals, []string{"version", "true"}) {
		t.Errorf("positionals = %q, want [version true]", p.positionals)
	}
}

func TestShorthandClusters(t *testing.T) {
	p := parseArgs([]string{"-hnns", "get"})
	if !p.boolValue("help") {
		t.Errorf("-hnns: boolValue(help) = false, want true")
	}
	if got, ok := p.value("namespace"); !ok || got != "ns" {
		t.Errorf("-hnns: value(namespace) = %q, %v; want ns, true", got, ok)
	}
	if !sameStrings(p.positionals, []string{"get"}) {
		t.Errorf("-hnns: positionals = %q, want [get]", p.positionals)
	}

	p = parseArgs([]string{"-hn", "ns", "get"})
	if got, _ := p.value("namespace"); got != "ns" || !p.boolValue("help") {
		t.Errorf("-hn ns: namespace = %q help = %v", got, p.boolValue("help"))
	}

	p = parseArgs([]string{"-hn=ns"})
	if got, _ := p.value("namespace"); got != "ns" {
		t.Errorf("-hn=ns: namespace = %q, want ns", got)
	}

	p = parseArgs([]string{"-hh=false"})
	if p.boolValue("help") || !p.has("help") {
		t.Errorf("-hh=false: boolValue(help) = %v has = %v, want false, true", p.boolValue("help"), p.has("help"))
	}

	p = parseArgs([]string{"-hxs", "X"})
	if !p.has("x") || !p.boolValue("help") {
		t.Errorf("-hxs X: unknown letter in a cluster must be recorded: %v", p.flags)
	}
	if got, _ := p.value("server"); got != "X" {
		t.Errorf("-hxs X: server = %q, want X", got)
	}

	p = parseArgs([]string{"get", "-ojson", "pods"})
	if p.has("server") || p.has("namespace") || !p.has("o") {
		t.Errorf("-ojson must be one unknown flag, got %v", p.flags)
	}
	if !sameStrings(p.positionals, []string{"get", "pods"}) {
		t.Errorf("-ojson: positionals = %q", p.positionals)
	}
}

func TestUnknownShortFlag(t *testing.T) {
	p := parseArgs([]string{"get", "-o", "json", "-A"})
	if !p.has("o") || !p.has("A") {
		t.Errorf("unknown short flags must be recorded under their letter: %v", p.flags)
	}
	if !sameStrings(p.positionals, []string{"get", "json"}) {
		t.Errorf("positionals = %q, want [get json]", p.positionals)
	}
}

func TestRemoveFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
		flag string
		want []string
	}{
		{"separate value", []string{"--kubeconfig", "f", "get", "pods"}, "kubeconfig", []string{"get", "pods"}},
		{"equals value", []string{"--kubeconfig=f", "get", "pods"}, "kubeconfig", []string{"get", "pods"}},
		{"every occurrence", []string{"--kubeconfig", "a", "get", "--kubeconfig=b", "pods"}, "kubeconfig", []string{"get", "pods"}},
		{"underscore spelling", []string{"--client_certificate", "p", "get"}, "client-certificate", []string{"get"}},
		{"underscore in requested name", []string{"--client-certificate=p", "get"}, "client_certificate", []string{"get"}},
		{"after terminator untouched", []string{"exec", "pod", "--", "kubectl", "--kubeconfig", "f"}, "kubeconfig", []string{"exec", "pod", "--", "kubectl", "--kubeconfig", "f"}},
		{"before and after terminator", []string{"--kubeconfig=f", "exec", "--", "--kubeconfig=g"}, "kubeconfig", []string{"exec", "--", "--kubeconfig=g"}},
		{"no value", []string{"--kubeconfig"}, "kubeconfig", []string{}},
		{"short spelling", []string{"-n", "ns", "get", "-nother", "pods", "-n=x"}, "namespace", []string{"get", "pods"}},
		{"cluster remove namespace keeps help", []string{"-hnns", "get"}, "namespace", []string{"-h", "get"}},
		{"cluster remove help keeps namespace", []string{"-hnns", "get"}, "help", []string{"-nns", "get"}},
		{"cluster remove namespace with separate value", []string{"-hn", "ns", "get"}, "namespace", []string{"-h", "get"}},
		{"cluster remove help with trailing value flag", []string{"-hn", "ns", "get"}, "help", []string{"-n", "ns", "get"}},
		{"cluster remove only flag drops token", []string{"-hh=false", "get"}, "help", []string{"get"}},
		{"cluster remove namespace from -nhx", []string{"-nhx", "get"}, "namespace", []string{"get"}},
		{"consumed value is not a flag", []string{"--namespace", "--kubeconfig", "get"}, "kubeconfig", []string{"--namespace", "--kubeconfig", "get"}},
		{"other flags kept", []string{"--context", "c", "--kubeconfig", "f", "get"}, "kubeconfig", []string{"--context", "c", "get"}},
		{"absent", []string{"get", "pods"}, "kubeconfig", []string{"get", "pods"}},
		{"nil input", nil, "kubeconfig", []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var before []string
			if tt.args != nil {
				before = append([]string{}, tt.args...)
			}
			got := removeFlag(tt.args, tt.flag)
			if got == nil || !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeFlag(%q, %q) = %#v, want %#v", tt.args, tt.flag, got, tt.want)
			}
			if tt.args != nil && !reflect.DeepEqual(tt.args, before) {
				t.Errorf("input mutated: %q, was %q", tt.args, before)
			}
		})
	}
}

func TestRemoveFlagDoesNotAliasInput(t *testing.T) {
	in := []string{"get", "pods"}
	out := removeFlag(in, "kubeconfig")
	out[0] = "changed"
	if in[0] != "get" {
		t.Errorf("returned slice aliases the input")
	}
}

func TestFlagTables(t *testing.T) {
	wantValue := []string{
		"as", "as-group", "as-uid", "as-user-extra", "cache-dir", "certificate-authority",
		"client-certificate", "client-key", "cluster", "context", "kubeconfig", "kuberc",
		"log-flush-frequency", "namespace", "password", "profile", "profile-output",
		"proxy-url", "request-timeout", "server", "tls-server-name", "token", "user",
		"username", "v", "vmodule",
	}
	wantBool := []string{
		"insecure-skip-tls-verify", "match-server-version", "disable-compression",
		"warnings-as-errors", "help", "client",
	}
	if len(valueFlags) != len(wantValue) {
		t.Errorf("valueFlags has %d entries, want %d", len(valueFlags), len(wantValue))
	}
	for _, name := range wantValue {
		if !valueFlags[name] {
			t.Errorf("valueFlags missing %q", name)
		}
		if boolFlags[name] {
			t.Errorf("boolFlags must not contain value flag %q", name)
		}
	}
	if len(boolFlags) != len(wantBool) {
		t.Errorf("boolFlags has %d entries, want %d", len(boolFlags), len(wantBool))
	}
	for _, name := range wantBool {
		if !boolFlags[name] {
			t.Errorf("boolFlags missing %q", name)
		}
		if valueFlags[name] {
			t.Errorf("valueFlags must not contain boolean flag %q", name)
		}
	}
	wantShort := map[string]string{"n": "namespace", "s": "server", "v": "v", "h": "help"}
	if !reflect.DeepEqual(shortFlags, wantShort) {
		t.Errorf("shortFlags = %v, want %v", shortFlags, wantShort)
	}
}

func TestIsAmbiguous(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{"unknown cluster hides server shorthand", []string{"get", "pods", "-Ashttps://other"}, true},
		{"unknown cluster with host and port", []string{"get", "pods", "-Asother.example.com:6443"}, true},
		{"label selector with dotted key", []string{"get", "pods", "-lapp.kubernetes.io/name=x"}, false},
		{"label selector with two terms", []string{"get", "pods", "-lapp=x,tier=db"}, false},
		{"unknown cluster with namespace equals", []string{"get", "pods", "-An=x"}, false},
		{"label selector env", []string{"get", "pods", "-lenv=prod"}, false},
		{"label selector run", []string{"get", "pods", "-lrun=x"}, false},
		{"unknown cluster with server equals", []string{"get", "pods", "-As=https://x"}, true},
		{"unknown long flag then server", []string{"get", "pods", "--template", "--server=https://x"}, true},
		{"unknown long flag then server separate", []string{"get", "pods", "--template", "--server", "https://x"}, true},
		{"unknown long flag then short server", []string{"get", "pods", "--template", "-s", "https://x"}, true},
		{"unknown long flag then short server attached", []string{"get", "pods", "--template", "-shttps://x"}, true},
		{"unknown long flag then kubeconfig", []string{"get", "pods", "--template", "--kubeconfig", "f"}, true},
		{"unknown long flag then cluster", []string{"get", "pods", "--template", "--cluster=c"}, true},
		{"show-labels then namespace short", []string{"get", "pods", "--show-labels", "-n", "foo"}, false},
		{"show-labels then namespace attached", []string{"get", "pods", "--show-labels", "-nfoo"}, false},
		{"all-namespaces then namespace long", []string{"get", "pods", "--all-namespaces", "--namespace", "foo"}, false},
		{"show-labels then context", []string{"get", "pods", "--show-labels", "--context=prod"}, true},
		{"unknown long flag then auth flags", []string{"get", "pods", "--watch", "--token=t", "--user", "u"}, false},
		{"unknown long flag then client key", []string{"get", "pods", "--template", "--client_key=k"}, false},
		{"unknown long flag then context", []string{"get", "pods", "--template", "--context=c"}, true},
		{"output flag cluster", []string{"get", "pods", "-ojson"}, false},
		{"all namespaces", []string{"get", "pods", "-A"}, false},
		{"template with separate value", []string{"get", "pods", "--template", "{{.}}", "--server=x"}, false},
		{"known value flag first", []string{"--output", "yaml", "get", "pods"}, false},
		{"unknown long flag then non-targeting flag", []string{"get", "pods", "--watch", "--output=yaml"}, false},
		{"unknown long flag with equals then server", []string{"get", "pods", "--template={{.}}", "--server=x"}, false},
		{"unknown long flag then terminator", []string{"exec", "pod", "--stdin", "--", "--server=x"}, false},
		{"after the terminator", []string{"exec", "pod", "--", "-Ashttps://x", "--template", "--server=x"}, false},
		{"known cluster parses fully", []string{"-hAshttps://x"}, false},
		{"known value flag consumes targeting-looking token", []string{"--namespace", "--server=x", "version"}, false},
		{"unknown long flag at end", []string{"get", "--template"}, false},
		{"no args", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseArgs(tt.args).isAmbiguous(); got != tt.want {
				t.Errorf("isAmbiguous() on %q = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestAmbiguousStillRecordsFlags(t *testing.T) {
	p := parseArgs([]string{"get", "pods", "--template", "--server=https://x"})
	if got, _ := p.value("server"); got != "https://x" {
		t.Errorf("value(server) = %q, want recorded despite ambiguity", got)
	}
}
