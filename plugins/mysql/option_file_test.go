package mysql

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Asserts what MySQL reads back, rather than the generated text.
func TestConfigFileEntryRoundTripsThroughMySQL(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"# starts the value", "#4b"},
		{"# inside the value", "a#b"},
		{"# ends the value", "pass#"},
		{"backslash", `back\slash`},
		{"backslash escape sequence", `tab\there`},
		{"trailing backslash", `secret\`},
		{"leading space", " leading"},
		{"trailing space", "trailing "},
		{"vertical tab and form feed at the ends", "\vsecret\f"},
		{"newline", "line\nbreak"},
		{"carriage return", "line\rbreak"},
		{"double quote", `double"quote`},
		{"single quote", "single'quote"},
		{"already wrapped in double quotes", `"secret"`},
		{"already wrapped in single quotes", "'secret'"},
		{"# with a double quote", `a"b#c`},
		{"# with a single quote", "a'b#c"},
		{"# with both kinds of quotation mark", `a'b"c#d`},
		{"semicolon", "semi;colon"},
		{"empty", ""},
		{"no special characters", "123456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line := configFileEntry("password", tt.value)

			assert.Equal(t, tt.value, readOptionFileValue(line),
				"mysql reads %q out of %q, but the 1Password item holds %q",
				readOptionFileValue(line), line, tt.value)
		})
	}
}

// Worked examples from the MySQL and MariaDB manuals and mysys/my_default.cc.
// This must pass whatever the encoder does, so a failure above cannot be blamed
// on the reader.
func TestReadOptionFileValueMatchesMySQL(t *testing.T) {
	tests := []struct {
		line string
		want string
	}{
		{`password=#4b`, ``},
		{`password=abc#def`, `abc`},
		// ";" only begins a comment at the start of a line.
		{`password=a;b`, `a;b`},
		{`password =  secret  `, `secret`},
		{`password="#4b"`, `#4b`},
		{`password='#4b'`, `#4b`},
		{`basedir=C:\\Program\sFiles\\MySQL`, `C:\Program Files\MySQL`},
		// An unrecognised sequence keeps its backslash.
		{`password=a\Sb`, `a\Sb`},
		// \" is decoded by mysys/my_default.cc, and does not close the quoting.
		{`password="a\"b"`, `a"b`},
		{`password="a\"b#c"`, `a"b#c`},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, readOptionFileValue(tt.line+"\n"),
			"reading option-file line %s", tt.line)
	}
}

// readOptionFileValue follows mysys/my_default.cc, in this order: drop the
// comment, delete surrounding whitespace, unwrap quotation marks, decode escapes.
func readOptionFileValue(line string) string {
	// MySQL reads one line at a time, so a raw line break ends the value.
	line, _, _ = strings.Cut(line, "\n")

	_, value, _ := strings.Cut(line, "=")

	// my_isspace() for latin1 accepts exactly this set.
	value = strings.Trim(removeOptionFileComment(value), " \t\n\v\f\r")

	if len(value) >= 2 && (value[0] == '"' || value[0] == '\'') && value[len(value)-1] == value[0] {
		value = value[1 : len(value)-1]
	}

	return decodeOptionFileEscapes(value)
}

// A '#' is only a comment outside quotation marks, and a backslash inside them
// hides a closing quotation mark.
func removeOptionFileComment(value string) string {
	var quote byte
	escaped := false

	for i := 0; i < len(value); i++ {
		c := value[i]

		if (c == '\'' || c == '"') && !escaped {
			switch quote {
			case 0:
				quote = c
			case c:
				quote = 0
			}
		}

		if quote == 0 && c == '#' {
			return value[:i]
		}

		escaped = quote != 0 && c == '\\' && !escaped
	}

	return value
}

// An unrecognised sequence keeps its backslash.
func decodeOptionFileEscapes(value string) string {
	var decoded strings.Builder

	for i := 0; i < len(value); i++ {
		// A backslash in the final position is not an escape character.
		if value[i] != '\\' || i == len(value)-1 {
			decoded.WriteByte(value[i])
			continue
		}

		i++
		switch value[i] {
		case 'b':
			decoded.WriteByte('\b')
		case 't':
			decoded.WriteByte('\t')
		case 'n':
			decoded.WriteByte('\n')
		case 'r':
			decoded.WriteByte('\r')
		case 's':
			decoded.WriteByte(' ')
		case '"':
			decoded.WriteByte('"')
		case '\'':
			decoded.WriteByte('\'')
		case '\\':
			decoded.WriteByte('\\')
		default:
			decoded.WriteByte('\\')
			decoded.WriteByte(value[i])
		}
	}

	return decoded.String()
}
