package app

import (
	"fmt"

	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/cli"
	"github.com/spf13/pflag"
)

// addAccount : アカウントを追加
func addAccount(setMain bool) error {
	// 認証
	authUser, err := shared.api.Auth(&shared.conf.Settings.Feature.Consumer)
	if err != nil {
		return err
	}
	shared.conf.Cred.Write(authUser)

	// メインユーザに設定
	if setMain {
		shared.conf.Settings.Feature.MainUser = authUser.UserName
	}

	return shared.conf.SaveAll()
}

// loginAccount : ログイン
func loginAccount(u string) error {
	// ログインするユーザを取得
	user, err := shared.conf.Cred.Get(u)
	if err != nil {
		return err
	}

	// 新しいユーザで初期化
	shared.api = api.New(&shared.conf.Settings.Feature.Consumer, user)
	return nil
}

func (a *App) newAccountCmd() *cli.Command {
	cmd := &cli.Command{
		Name:      "account",
		Shorthand: "a",
		Short:     "Manage your account",
		Validate:  cli.NoArgs(),
	}

	cmd.AddCommand(
		a.newAccountAddCmd(),
		a.newAccountDeleteCmd(),
		a.newAccountListCmd(),
		a.newAccountSwitchCmd(),
	)

	return cmd
}

func (a *App) newAccountAddCmd() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Shorthand: "a",
		Short:     "Add account",
		Hidden:    !shared.isCommandLineMode,
		Validate:  cli.NoArgs(),
		SetFlag: func(f *pflag.FlagSet) {
			f.BoolP("main", "m", false, "set as main user")
		},
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			main, _ := f.GetBool("main")
			return addAccount(main)
		},
	}
}

func (a *App) newAccountDeleteCmd() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Shorthand: "d",
		Short:     "Delete account",
		Hidden:    !shared.isCommandLineMode,
		Validate:  cli.RequireArgs(1),
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			return shared.conf.Cred.Delete(f.Arg(0))
		},
	}
}

func (a *App) newAccountListCmd() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Shorthand: "l",
		Short:     "Show accounts that have been added",
		Hidden:    !shared.isCommandLineMode,
		Validate:  cli.NoArgs(),
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			for _, u := range shared.conf.Cred.GetAllNames() {
				current := " "
				if u == shared.api.CurrentUser.UserName {
					current = "*"
				}

				fmt.Printf(" %s %s\n", current, u)
			}

			return nil
		},
	}
}

func (a *App) newAccountSwitchCmd() *cli.Command {
	return &cli.Command{
		Name:      "switch",
		Shorthand: "s",
		Short:     "Switch the account to be used",
		Validate:  cli.RequireArgs(1),
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			// ログイン
			if err := loginAccount(f.Arg(0)); err != nil {
				return err
			}

			// アプリを初期化
			a.view.Reset()
			a.statusBar.DrawAccountInfo()
			a.initAutocomplate()
			a.runStartupCommands()

			return nil
		},
	}
}
