package app

import (
	"github.com/arrow2nd/nekome/cli"
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
		Name:     "add",
		Short:    "Add account",
		Hidden:   !shared.isCommandLineMode,
		Validate: cli.NoArgs(),
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			return nil
		},
	}
}

func (a *App) newAccountDeleteCmd() *cli.Command {
	return &cli.Command{
		Name:     "delete",
		Short:    "Delete account",
		Hidden:   !shared.isCommandLineMode,
		Validate: cli.RequireArgs(1),
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			return nil
		},
	}
}

func (a *App) newAccountListCmd() *cli.Command {
	return &cli.Command{
		Name:     "list",
		Short:    "Show accounts that have been added",
		Hidden:   !shared.isCommandLineMode,
		Validate: cli.NoArgs(),
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			return nil
		},
	}
}

func (a *App) newAccountSwitchCmd() *cli.Command {
	return &cli.Command{
		Name:     "switch",
		Short:    "Switch the account to be used",
		Validate: cli.RequireArgs(1),
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			return nil
		},
	}
}
