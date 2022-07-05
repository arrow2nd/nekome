package app

import (
	"github.com/arrow2nd/nekome/cli"
	"github.com/spf13/pflag"
)

func (a *App) newDocsCmd() *cli.Command {
	cmd := &cli.Command{
		Name:      "docs",
		Shorthand: "d",
		Short:     "Show the document",
		Hidden:    shared.isCommandLineMode,
		Validate:  cli.NoArgs(),
	}

	cmd.AddCommand(a.newDocShortcutsCmd())

	return cmd
}

func (a *App) newDocShortcutsCmd() *cli.Command {
	text := `[-:-:b]System
[black:gray:-] ctrl+l [-:-:-] Redraw screen (window width changes are not reflected)
[black:gray:-] ctrl+w [-:-:-] Close current page
[black:gray:-] ctrl+q [-:-:-] Exit Application

[-:-:b]Navigation
[black:gray:-] j [-:-:-] [black:gray:-] up [-:-:-]    Focus the next tweet
[black:gray:-] k [-:-:-] [black:gray:-] down [-:-:-]  Focus the previous tweet
[black:gray:-] g [-:-:-] [black:gray:-] home [-:-:-]  Focus the top tweet
[black:gray:-] G [-:-:-] [black:gray:-] end [-:-:-]   Focus the bottom tweet
[black:gray:-] h [-:-:-] [black:gray:-] left [-:-:-]  Focus the previous tab
[black:gray:-] l [-:-:-] [black:gray:-] right [-:-:-] Focus the next tab
[black:gray:-] : [-:-:-]         Focus the command line

[-:-:b]Scrolling
[black:gray:-] ctrl+j [-:-:-] [black:gray:-] page up [-:-:-]   Scroll up
[black:gray:-] ctrl+k [-:-:-] [black:gray:-] page down [-:-:-] Scroll down

[-:-:b]Tweet Navigation
[black:gray:-] f [-:-:-] Like a tweet
[black:gray:-] F [-:-:-] Unlike a tweet
[black:gray:-] t [-:-:-] Retweet a tweet
[black:gray:-] T [-:-:-] Unretweet a tweet
[black:gray:-] q [-:-:-] Quote tweet
[black:gray:-] r [-:-:-] Reply to
[black:gray:-] D [-:-:-] Delete a tweet
[black:gray:-] o [-:-:-] Open in browser
[black:gray:-] i [-:-:-] Open author's user timeline page
[black:gray:-] c [-:-:-] Copy link to clipboard

[-:-:b]User Navigation
[black:gray:-] w [-:-:-] Follow a user
[black:gray:-] W [-:-:-] Unfollow a user
[black:gray:-] u [-:-:-] Mute a user
[black:gray:-] U [-:-:-] Unmute a user
[black:gray:-] x [-:-:-] Block a user
[black:gray:-] X [-:-:-] Unblock a user
`
	return &cli.Command{
		Name:      "shortcuts",
		Shorthand: "s",
		Short:     "Documentation for shortcut keys",
		Validate:  cli.NoArgs(),
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			return a.view.AddPage(newDocsPage("Shortcuts", text), true)
		},
	}
}
