package app

import (
	"github.com/arrow2nd/nekome/cli"
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
	// TODO: カスタムしたキーバインドを反映するようにする
	text := `Global:
  ctrl-l   Redraw screen (window width changes are not reflected)
  ctrl-w   Close current page
  ctrl-q   Quit application

Navigation:
  j up      Focus the next tweet
  k down    Focus the previous tweet
  g home    Focus the top tweet
  G end     Focus the bottom tweet
  h left    Focus the previous tab
  l right   Focus the next tab
  :         Focus the command line

Scrolling:
  ctrl-j page up     Scroll up
  ctrl-k page down   Scroll down

Home Timeline Page:
  s   Start stream mode (similar to UserStream)
  S   Exit stream mode

Tweet Navigation:
  f   Like a tweet
  F   Unlike a tweet
  t   Retweet a tweet
  T   Unretweet a tweet
  q   Quote tweet
  r   Reply to
  D   Delete a tweet
  o   Open in browser
  i   Open author's user timeline page
  c   Copy link to clipboard

User Navigation:
  w   Follow a user
  W   Unfollow a user
  u   Mute a user
  U   Unmute a user
  x   Block a user
  X   Unblock a user
`
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
