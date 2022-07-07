package app

import (
	"errors"
	"strings"

	"github.com/arrow2nd/nekome/cli"
	"github.com/spf13/pflag"
)

func setUnfocusFlag(f *pflag.FlagSet) {
	f.BoolP("unfocus", "u", false, "no focus on page")
}

func (a *App) newHomeCmd() *cli.Command {
	return &cli.Command{
		Name:      "home",
		Shorthand: "h",
		Short:     "Add home timeline page",
		Validate:  cli.NoArgs(),
		Hidden:    shared.isCommandLineMode,
		SetFlag:   setUnfocusFlag,
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			unfocus, _ := f.GetBool("unfocus")
			return a.view.AddPage(newTimelinePage(homeTL), !unfocus)
		},
	}
}

func (a *App) newMentionCmd() *cli.Command {
	return &cli.Command{
		Name:      "mention",
		Shorthand: "m",
		Short:     "Add mention timeline page",
		Validate:  cli.NoArgs(),
		Hidden:    shared.isCommandLineMode,
		SetFlag:   setUnfocusFlag,
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			unfocus, _ := f.GetBool("unfocus")
			return a.view.AddPage(newTimelinePage(mentionTL), !unfocus)
		},
	}
}

func (a *App) newListCmd() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Shorthand: "l",
		Short:     "Add list timeline page",
		UsageArgs: "<list name> <list ID>",
		Example:   "list cathouse 1234567890",
		Validate:  cli.RequireArgs(2),
		Hidden:    shared.isCommandLineMode,
		SetFlag:   setUnfocusFlag,
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			name := f.Arg(0)
			if name == "" {
				return errors.New("please specify list name")
			}

			id := f.Arg(1)
			if id == "" {
				return errors.New("please specify list id")
			}

			unfocus, _ := f.GetBool("unfocus")
			return a.view.AddPage(newListPage(name, id), !unfocus)
		},
	}
}

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

			unfocus, _ := f.GetBool("unfocus")
			return a.view.AddPage(newUserPage(userName), !unfocus)
		},
	}
}

func (a *App) newSearchCmd() *cli.Command {
	return &cli.Command{
		Name:      "search",
		Shorthand: "s",
		Short:     "Add seaech result page",
		Long: `Add seaech result page.

If the query contains spaces, enclose it in double quotes.`,
		UsageArgs: "<query>",
		Example:   "search golang",
		Validate:  cli.RequireArgs(1),
		Hidden:    shared.isCommandLineMode,
		SetFlag:   setUnfocusFlag,
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			query := f.Arg(0)
			if query == "" {
				return errors.New("please specify search query")
			}

			unfocus, _ := f.GetBool("unfocus")
			return a.view.AddPage(newSearchPage(query), !unfocus)
		},
	}
}
