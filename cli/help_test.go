package cli_test

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestHelp(t *testing.T) {
	a := newCmd("testhelp")
	a.Short = "short desc"
	a.SetFlagFunc = func(f *pflag.FlagSet) {
		f.Bool("bool", false, "bool usage")
		f.StringP("string", "s", "", "string usage")
	}

	r := newCmd("root")
	r.Alias = "r"
	r.Short = "root command"
	r.Example = "example"
	r.HelpFunc = func(h string) {
		assert.Equal(t, h, `root command

Usage:
  root [command] [flags]

Alias:
  r

Example:
  example

Commands:
  testhelp	short desc

Flags:
  -h, --help   help for root
`,
			"正しいヘルプが生成されているか",
		)
	}
	r.AddCommand(a)

	r.Execute([]string{"--help"})
}
