package app

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
)

func (a *App) newHomeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "home",
		Long: "Add Home Timeline Page",
		Args: cobra.NoArgs,
	}

	cmd.Flags().BoolP("unfocus", "u", false, "No focus on page")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		unfocus, _ := cmd.Flags().GetBool("unfocus")
		return a.view.AddPage(newTimelinePage(homeTL), !unfocus)
	}

	return cmd
}

func (a *App) newMentionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "mention",
		Long: "Add Mention Timeline Page",
		Args: cobra.NoArgs,
	}

	cmd.Flags().BoolP("unfocus", "u", false, "No focus on page")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		unfocus, _ := cmd.Flags().GetBool("unfocus")
		return a.view.AddPage(newTimelinePage(mentionTL), !unfocus)
	}

	return cmd
}

func (a *App) newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Long:    "Add List Timeline Page",
		Example: "list <list name> <list id>",
		Args:    cobra.ExactValidArgs(2),
	}

	cmd.Flags().BoolP("unfocus", "u", false, "No focus on page")

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

func (a *App) newUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user",
		Long:    "Add User Timeline Page",
		Example: "user <user name>",
		Args:    cobra.ExactValidArgs(1),
	}

	cmd.Flags().BoolP("unfocus", "u", false, "No focus on page")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		userName := args[0]

		// ユーザの指定がないなら自分を指定
		if userName == "" {
			userName = shared.api.CurrentUser.UserName
		}

		// @を除去
		userName = strings.ReplaceAll(userName, "@", "")

		unfocus, _ := cmd.Flags().GetBool("unfocus")
		return a.view.AddPage(newUserPage(userName), !unfocus)
	}

	return cmd
}

func (a *App) newSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "search",
		Long:    "Add Seaech Result Page",
		Example: "search <query>",
		Args:    cobra.ExactValidArgs(1),
	}

	cmd.Flags().BoolP("unfocus", "u", false, "No focus on page")

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
