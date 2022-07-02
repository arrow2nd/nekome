package app

import (
	"errors"
	"strings"

	"github.com/arrow2nd/nekome/cli"
	"github.com/spf13/pflag"
)

// setUnfocusFlagc : unfocusフラグを設定
func setUnfocusFlag(f *pflag.FlagSet) {
	f.BoolP("unfocus", "u", false, "no focus on page")
}

// newHomeCmd : homeコマンド生成
func (a *App) newHomeCmd() *cli.Command {
	return &cli.Command{
		Name:         "home",
		Alias:        "h",
		Short:        "Add home timeline page",
		ValidateFunc: cli.NoArgs(),
		Hidden:       shared.isCommandLineMode,
		SetFlagFunc:  setUnfocusFlag,
		RunFunc: func(c *cli.Command, f *pflag.FlagSet) error {
			unfocus, _ := f.GetBool("unfocus")
			return a.view.AddPage(newTimelinePage(homeTL), !unfocus)
		},
	}
}

// newMentionCmd : mentionコマンド生成
func (a *App) newMentionCmd() *cli.Command {
	return &cli.Command{
		Name:         "mention",
		Alias:        "m",
		Short:        "Add mention timeline page",
		ValidateFunc: cli.NoArgs(),
		Hidden:       shared.isCommandLineMode,
		SetFlagFunc:  setUnfocusFlag,
		RunFunc: func(c *cli.Command, f *pflag.FlagSet) error {
			unfocus, _ := f.GetBool("unfocus")
			return a.view.AddPage(newTimelinePage(mentionTL), !unfocus)
		},
	}
}

// newListCmd : listコマンド生成
func (a *App) newListCmd() *cli.Command {
	return &cli.Command{
		Name:         "list",
		Alias:        "l",
		Short:        "Add list timeline page",
		Example:      "list cathouse 1234567890",
		ValidateFunc: cli.RequireArgs(2),
		Hidden:       shared.isCommandLineMode,
		SetFlagFunc:  setUnfocusFlag,
		RunFunc: func(c *cli.Command, f *pflag.FlagSet) error {
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

// newUserCmd : userコマンド生成
func (a *App) newUserCmd() *cli.Command {
	return &cli.Command{
		Name:         "user",
		Alias:        "u",
		Short:        "Add user timeline page",
		Example:      "user github",
		ValidateFunc: cli.RangeArgs(0, 1),
		Hidden:       shared.isCommandLineMode,
		SetFlagFunc:  setUnfocusFlag,
		RunFunc: func(c *cli.Command, f *pflag.FlagSet) error {
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

// newSearchCmd : searchコマンド生成
func (a *App) newSearchCmd() *cli.Command {
	return &cli.Command{
		Name:         "search",
		Alias:        "s",
		Short:        "Add seaech result page",
		Example:      "  search golang",
		ValidateFunc: cli.RequireArgs(1),
		Hidden:       shared.isCommandLineMode,
		SetFlagFunc:  setUnfocusFlag,
		RunFunc: func(c *cli.Command, f *pflag.FlagSet) error {
			query := f.Arg(0)
			if query == "" {
				return errors.New("please specify search keywords")
			}

			unfocus, _ := f.GetBool("unfocus")
			return a.view.AddPage(newSearchPage(query), !unfocus)
		},
	}
}
