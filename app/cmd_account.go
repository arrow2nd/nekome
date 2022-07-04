package app

import (
	"errors"
	"fmt"

	"github.com/arrow2nd/nekome/cli"
	"github.com/manifoldco/promptui"
	"github.com/spf13/pflag"
)

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
			if err := addAccount(main); err != nil {
				return err
			}

			fmt.Println("successfully added")
			return nil
		},
	}
}

func (a *App) newAccountDeleteCmd() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Shorthand: "d",
		Short:     "Delete account",
		Long: `Delete account.
If you do not specify an account name, you can select it interactively.`,
		UsageArgs: "[user name]",
		Example:   "delete arrow2nd",
		Hidden:    !shared.isCommandLineMode,
		Validate:  cli.RangeArgs(0, 1),
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			target := f.Arg(0)

			// 指定が無い場合、ユーザを選択
			if target == "" {
				prompt := promptui.Select{
					Label: "Account to delete",
					Items: shared.conf.Cred.GetAllNames(),
				}

				_, seletecd, err := prompt.Run()
				if err != nil {
					return err
				}

				target = seletecd
			}

			// 削除実行
			if err := shared.conf.Cred.Delete(target); err != nil {
				return err
			}

			// 書き込み
			if err := shared.conf.SaveCred(); err != nil {
				return err
			}

			fmt.Printf("successfully deleted: %s\n", f.Arg(0))
			return nil
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
		UsageArgs: "[user name]",
		Example:   "switch arrow2nd",
		Hidden:    shared.isCommandLineMode,
		Validate:  cli.RequireArgs(1),
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			// 既にログイン中なら切り替えない
			if f.Arg(0) == shared.api.CurrentUser.UserName {
				return errors.New("account currently logged in")
			}

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
