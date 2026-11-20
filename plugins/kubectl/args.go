package kubectl

import (
	"regexp"
	"strconv"
	"strings"
)

type parsedArgs struct {
	flags       map[string][]string
	positionals []string
	rest        []string
	ambiguous   bool
}

var valueFlags = map[string]bool{
	"as":                    true,
	"as-group":              true,
	"as-uid":                true,
	"as-user-extra":         true,
	"cache-dir":             true,
	"certificate-authority": true,
	"client-certificate":    true,
	"client-key":            true,
	"cluster":               true,
	"context":               true,
	"kubeconfig":            true,
	"kuberc":                true,
	"log-flush-frequency":   true,
	"namespace":             true,
	"password":              true,
	"profile":               true,
	"profile-output":        true,
	"proxy-url":             true,
	"request-timeout":       true,
	"server":                true,
	"tls-server-name":       true,
	"token":                 true,
	"user":                  true,
	"username":              true,
	"v":                     true,
	"vmodule":               true,
}

var boolFlags = map[string]bool{
	"insecure-skip-tls-verify": true,
	"match-server-version":     true,
	"disable-compression":      true,
	"warnings-as-errors":       true,
	"help":                     true,
	"client":                   true,
}

var shortFlags = map[string]string{"n": "namespace", "s": "server", "v": "v", "h": "help"}

type argToken struct {
	name  string
	value string
	start int
	end   int
	from  int
	to    int
	short bool
}

func normalizeFlagName(name string) string {
	return strings.ReplaceAll(name, "_", "-")
}

func scanArgs(args []string) (tokens []argToken, positionals []string, terminator int, ambiguous bool) {
	terminator = -1
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--":
			return tokens, positionals, i, ambiguous
		case strings.HasPrefix(arg, "--"):
			name, value, hasValue := strings.Cut(arg[2:], "=")
			name = normalizeFlagName(name)
			tok := argToken{name: name, start: i, end: i + 1}
			switch {
			case hasValue:
				tok.value = value
			case valueFlags[name]:
				if i+1 < len(args) {
					i++
					tok.value = args[i]
					tok.end = i + 1
				}
			default:
				tok.value = "true"
				if !boolFlags[name] && i+1 < len(args) && isTargetingFlag(args[i+1]) {
					ambiguous = true
				}
			}
			tokens = append(tokens, tok)
		case strings.HasPrefix(arg, "-") && len(arg) > 1:
			var consumed int
			tokens, consumed = scanShorthands(tokens, args, i)
			i += consumed
			if _, known := shortFlags[arg[1:2]]; !known && hidesTargetingShorthand(arg) {
				ambiguous = true
			}
		default:
			positionals = append(positionals, arg)
		}
	}
	return tokens, positionals, terminator, ambiguous
}

var targetingFlags = map[string]bool{
	"server": true, "context": true, "cluster": true, "kubeconfig": true,
}

func isTargetingFlag(arg string) bool {
	switch {
	case strings.HasPrefix(arg, "--") && len(arg) > 2:
		name, _, _ := strings.Cut(arg[2:], "=")
		return targetingFlags[normalizeFlagName(name)]
	case strings.HasPrefix(arg, "-") && len(arg) > 1:
		return targetingFlags[shortFlags[arg[1:2]]]
	}
	return false
}

var hostPortPattern = regexp.MustCompile(`:[0-9]+`)

func hidesTargetingShorthand(arg string) bool {
	for j := 2; j < len(arg); j++ {
		if arg[j] != 's' {
			continue
		}
		rest := arg[j+1:]
		if strings.HasPrefix(rest, "=") || strings.Contains(rest, "://") || hostPortPattern.MatchString(rest) {
			return true
		}
	}
	return false
}

func scanShorthands(tokens []argToken, args []string, i int) ([]argToken, int) {
	arg := args[i]
	if _, known := shortFlags[arg[1:2]]; !known {
		return append(tokens, argToken{name: arg[1:2], value: "true", start: i, end: i + 1, from: 1, to: len(arg), short: true}), 0
	}
	for j := 1; j < len(arg); j++ {
		letter := arg[j : j+1]
		name, known := shortFlags[letter]
		tok := argToken{name: letter, value: "true", start: i, end: i + 1, from: j, to: j + 1, short: true}
		if !known {
			if j+1 < len(arg) && arg[j+1] == '=' {
				tok.to = len(arg)
				return append(tokens, tok), 0
			}
			tokens = append(tokens, tok)
			continue
		}
		tok.name = name
		if valueFlags[name] {
			tok.to = len(arg)
			if j+1 < len(arg) {
				tok.value = strings.TrimPrefix(arg[j+1:], "=")
				return append(tokens, tok), 0
			}
			tok.value = ""
			consumed := 0
			if i+1 < len(args) {
				tok.value = args[i+1]
				tok.end = i + 2
				consumed = 1
			}
			return append(tokens, tok), consumed
		}
		if j+1 < len(arg) && arg[j+1] == '=' {
			tok.value = arg[j+2:]
			tok.to = len(arg)
			return append(tokens, tok), 0
		}
		tokens = append(tokens, tok)
	}
	return tokens, 0
}

func parseArgs(args []string) parsedArgs {
	tokens, positionals, terminator, ambiguous := scanArgs(args)
	p := parsedArgs{flags: map[string][]string{}, positionals: positionals, ambiguous: ambiguous}
	for _, tok := range tokens {
		p.flags[tok.name] = append(p.flags[tok.name], tok.value)
	}
	if terminator >= 0 {
		p.rest = append([]string(nil), args[terminator+1:]...)
	}
	return p
}

func (p parsedArgs) value(name string) (string, bool) {
	values := p.flags[normalizeFlagName(name)]
	if len(values) == 0 {
		return "", false
	}
	last := values[len(values)-1]
	return last, last != ""
}

func (p parsedArgs) has(name string) bool {
	_, ok := p.flags[normalizeFlagName(name)]
	return ok
}

func (p parsedArgs) boolValue(name string) bool {
	values := p.flags[normalizeFlagName(name)]
	if len(values) == 0 {
		return false
	}
	b, err := strconv.ParseBool(values[len(values)-1])
	return err == nil && b
}

func (p parsedArgs) isAmbiguous() bool { return p.ambiguous }

func (p parsedArgs) subcommand() string {
	if len(p.positionals) == 0 {
		return ""
	}
	return p.positionals[0]
}

func removeFlag(args []string, name string) []string {
	name = normalizeFlagName(name)
	tokens, _, _, _ := scanArgs(args)
	dropped := make([]bool, len(args))
	cuts := map[int][][2]int{}
	for _, tok := range tokens {
		if tok.name != name {
			continue
		}
		if !tok.short {
			for k := tok.start; k < tok.end; k++ {
				dropped[k] = true
			}
			continue
		}
		cuts[tok.start] = append(cuts[tok.start], [2]int{tok.from, tok.to})
		for k := tok.start + 1; k < tok.end; k++ {
			dropped[k] = true
		}
	}
	out := make([]string, 0, len(args))
	for k, arg := range args {
		if dropped[k] {
			continue
		}
		if ranges, ok := cuts[k]; ok {
			arg = cutRanges(arg, ranges)
			if arg == "-" {
				continue
			}
		}
		out = append(out, arg)
	}
	return out
}

func cutRanges(arg string, ranges [][2]int) string {
	var b strings.Builder
	pos := 0
	for _, r := range ranges {
		b.WriteString(arg[pos:r[0]])
		pos = r[1]
	}
	b.WriteString(arg[pos:])
	return b.String()
}
