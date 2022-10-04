package cli_test

import (
	"testing"

	"github.com/arrow2nd/nekome/v2/cli"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestHelp(t *testing.T) {
	a := newCmd("testhelp")
	a.Short = "short desc"
	a.SetFlag = func(f *pflag.FlagSet) {
		f.Bool("bool", false, "bool usage")
		f.StringP("string", "s", "", "string usage")
	}

	r := newCmd("root")
	r.Shorthand = "r"
	r.Short = "root command"
	r.Example = "example"
	r.Help = func(c *cli.Command, h string) {
		want := `root command

Usage:
  root [flags] [command]

Shorthand:
  r

Example:
  example

Commands:
  testhelp   short desc

Flags:
  -h, --help   help for root
`

		assert.Equal(t, want, h, "正しいヘルプが生成されているか")
	}

	r.AddCommand(a)

	t.Run("実行できるか", func(t *testing.T) {
		assert.NoError(t, r.Execute([]string{"--help"}))
	})
}
