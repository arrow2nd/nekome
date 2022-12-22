package app

import (
	"strings"

	"github.com/arrow2nd/nekome/v2/cli"
	"github.com/spf13/pflag"
)

func (a *App) newUserCmd() *cli.Command {
	return &cli.Command{
		Name:      "user",
		Shorthand: "u",
		Short:     "Add user timeline page",
		Long: `Add user timeline page.

The @ in the user name can be omitted.
If no user name is specified, the currently logged-in user is specified.`,
		UsageArgs: "[user name]",
		Example:   "user github",
		Validate:  cli.RangeArgs(0, 1),
		Hidden:    shared.isCommandLineMode,
		SetFlag:   setUnfocusFlag,
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			userName := shared.api.CurrentUser.UserName

			// ユーザの指定があるなら置き換え
			if f.NArg() > 0 {
				userName = f.Arg(0)
			}

			// @を除去
			userName = strings.ReplaceAll(userName, "@", "")

			page, err := newUserPage(userName)
			if err != nil {
				return err
			}

			unfocus, _ := f.GetBool("unfocus")
			if err := a.view.AddPage(page, !unfocus); err != nil {
				return err
			}

			go page.Load()

			return nil
		},
	}
}
