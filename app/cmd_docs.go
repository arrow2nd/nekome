package app

import (
	"fmt"

	"github.com/arrow2nd/nekome/v2/cli"
	"github.com/arrow2nd/nekome/v2/config"
	"github.com/spf13/pflag"
)

func (a *App) newDocsCmd() *cli.Command {
	cmd := &cli.Command{
		Name:      "docs",
		Shorthand: "d",
		Short:     "Show documentation",
		Hidden:    shared.isCommandLineMode,
		Validate:  cli.NoArgs(),
	}

	cmd.AddCommand(a.newDocsKeybindingsCmd())

	return cmd
}

func (a *App) newDocsKeybindingsCmd() *cli.Command {
	k := shared.conf.Pref.Keybindings

	global := fmt.Sprintf(
		`[Global]
  %-20s Quit application

`,
		k.Global.GetString(config.ActionQuit),
	)

	view := fmt.Sprintf(
		`[View]
  %-20s Select previous tab
  %-20s Select next tab
  %-20s Close current page
  %-20s Redraw screen (window width changes are not reflected)
  %-20s Focus the command line
  %-20s Show documentation for keybindings

`,
		k.View.GetString(config.ActionSelectPrevTab),
		k.View.GetString(config.ActionSelectNextTab),
		k.View.GetString(config.ActionClosePage),
		k.View.GetString(config.ActionRedraw),
		k.View.GetString(config.ActionFocusCmdLine),
		k.View.GetString(config.ActionShowHelp),
	)

	page := fmt.Sprintf(
		`[Common Page]
  %-20s Reload page

`,
		k.Page.GetString(config.ActionReloadPage),
	)

	home := fmt.Sprintf(
		`[Home Timeline Page]
  %-20s Start stream mode (similar to UserStream)
  %-20s Stop stream mode

`,
		k.HomeTimeline.GetString(config.ActionStreamModeStart),
		k.HomeTimeline.GetString(config.ActionStreamModeStop),
	)

	tweet := fmt.Sprintf(
		`[Tweet View]
  %-20s Scroll up
  %-20s Scroll down
  %-20s Move cursor up
  %-20s Move cursor down
  %-20s Move cursor top
  %-20s Move cursor bottom

  %-20s Like
  %-20s Unlike
  %-20s Retweet
  %-20s Unretweet
  %-20s New tweet
  %-20s Quote tweet
  %-20s Reply to tweet
  %-20s Delete a tweet
  %-20s Open in browser
  %-20s Copy link to clipboard
  
  %-20s Follow
  %-20s Unfollow
  %-20s Mute
  %-20s Unmute
  %-20s Block
  %-20s Unblock
  %-20s Open user timeline page
  %-20s Open user likes page
`,
		k.TweetView.GetString(config.ActionScrollUp),
		k.TweetView.GetString(config.ActionScrollDown),
		k.TweetView.GetString(config.ActionCursorUp),
		k.TweetView.GetString(config.ActionCursorDown),
		k.TweetView.GetString(config.ActionCursorTop),
		k.TweetView.GetString(config.ActionCursorBottom),
		k.TweetView.GetString(config.ActionTweetLike),
		k.TweetView.GetString(config.ActionTweetUnlike),
		k.TweetView.GetString(config.ActionTweetRetweet),
		k.TweetView.GetString(config.ActionTweetUnretweet),
		k.TweetView.GetString(config.ActionTweet),
		k.TweetView.GetString(config.ActionQuote),
		k.TweetView.GetString(config.ActionReply),
		k.TweetView.GetString(config.ActionTweetDelete),
		k.TweetView.GetString(config.ActionOpenBrowser),
		k.TweetView.GetString(config.ActionCopyUrl),
		k.TweetView.GetString(config.ActionUserFollow),
		k.TweetView.GetString(config.ActionUserUnfollow),
		k.TweetView.GetString(config.ActionUserMute),
		k.TweetView.GetString(config.ActionUserUnmute),
		k.TweetView.GetString(config.ActionUserBlock),
		k.TweetView.GetString(config.ActionUserUnblock),
		k.TweetView.GetString(config.ActionOpenUserPage),
		k.TweetView.GetString(config.ActionOpenUserLikes),
	)

	text := global + view + page + home + tweet

	return &cli.Command{
		Name:      "keybindings",
		Shorthand: "k",
		Short:     "Documentation for keybindings",
		Validate:  cli.NoArgs(),
		SetFlag:   setUnfocusFlag,
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			unfocus, _ := f.GetBool("unfocus")
			return a.view.AddPage(newDocsPage("Keybindings", text), !unfocus)
		},
	}
}
