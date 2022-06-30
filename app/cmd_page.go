package app

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
)

// newHomeCmd : homeコマンド生成
func (a *App) newHomeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "home",
		Aliases: []string{"h"},
		Short:   "add home timeline page",
		Args:    cobra.NoArgs,
		Hidden:  shared.isCommandLineMode,
	}

	cmd.Flags().BoolP("unfocus", "u", false, "no focus on page")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		unfocus, _ := cmd.Flags().GetBool("unfocus")
		return a.view.AddPage(newTimelinePage(homeTL), !unfocus)
	}

	return cmd
}

// newMentionCmd : mentionコマンド生成
func (a *App) newMentionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mention",
		Aliases: []string{"m"},
		Short:   "add mention timeline page",
		Args:    cobra.NoArgs,
		Hidden:  shared.isCommandLineMode,
	}

	cmd.Flags().BoolP("unfocus", "u", false, "no focus on page")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		unfocus, _ := cmd.Flags().GetBool("unfocus")
		return a.view.AddPage(newTimelinePage(mentionTL), !unfocus)
	}

	return cmd
}

// newListCmd : listコマンド生成
func (a *App) newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "add list timeline page",
		Example: "list <list name> <list id>",
		Args:    cobra.ExactValidArgs(2),
		Hidden:  shared.isCommandLineMode,
	}

	cmd.Flags().BoolP("unfocus", "u", false, "no focus on page")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if name == "" {
			return errors.New("please specify list name")
		}

		id := args[1]
		if id == "" {
			return errors.New("please specify list id")
		}

		unfocus, _ := cmd.Flags().GetBool("unfocus")
		return a.view.AddPage(newListPage(name, id), !unfocus)
	}

	return cmd
}

// newUserCmd : userコマンド生成
func (a *App) newUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user",
		Aliases: []string{"u"},
		Short:   "add user timeline page",
		Example: "user <user name>",
		Args:    cobra.RangeArgs(0, 1),
		Hidden:  shared.isCommandLineMode,
	}

	cmd.Flags().BoolP("unfocus", "u", false, "no focus on page")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		userName := shared.api.CurrentUser.UserName

		// ユーザの指定があるなら置き換え
		if len(args) > 0 {
			userName = args[0]
		}

		// @を除去
		userName = strings.ReplaceAll(userName, "@", "")

		unfocus, _ := cmd.Flags().GetBool("unfocus")
		return a.view.AddPage(newUserPage(userName), !unfocus)
	}

	return cmd
}

// newSearchCmd : searchコマンド生成
func (a *App) newSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "search",
		Aliases: []string{"s"},
		Short:   "add seaech result page",
		Example: "search <query>",
		Args:    cobra.ExactValidArgs(1),
		Hidden:  shared.isCommandLineMode,
	}

	cmd.Flags().BoolP("unfocus", "u", false, "no focus on page")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		query := args[0]
		if query == "" {
			return errors.New("please specify search keywords")
		}

		unfocus, _ := cmd.Flags().GetBool("unfocus")
		return a.view.AddPage(newSearchPage(query), !unfocus)
	}

	return cmd
}
