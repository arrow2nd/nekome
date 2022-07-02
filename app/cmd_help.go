package app

import (
	"github.com/arrow2nd/nekome/cli"
	"github.com/spf13/pflag"
)

func (a *App) newHelpShortCuts() *cli.Command {
	help := `[-:-:b]System
[-:gray:-] ctrl+l [-:-:-] Redraw screen (window width changes are not reflected)
[-:gray:-] ctrl+w [-:-:-] Close current page
[-:gray:-] ctrl+q [-:-:-] Exit Application

[-:-:b]Navigation
[-:gray:-] j, up [-:-:-]    Focus the next tweet
[-:gray:-] k, down [-:-:-]  Focus the previous tweet
[-:gray:-] g, home [-:-:-]  Focus the top tweet
[-:gray:-] G, end [-:-:-]   Focus the bottom tweet
[-:gray:-] h, left [-:-:-]  Focus the previous tab
[-:gray:-] l, right [-:-:-] Focus the next tab
[-:gray:-] : [-:-:-]        Focus the command line

[-:-:b]Scrolling
[-:gray:-] ctrl+j, page up [-:-:-]   Scroll up
[-:gray:-] ctrl+k, page down [-:-:-] Scroll down

[-:-:b]Tweet Navigation
[-:gray:-] f [-:-:-] Like a tweet
[-:gray:-] F [-:-:-] Unlike a tweet
[-:gray:-] t [-:-:-] Retweet a tweet
[-:gray:-] T [-:-:-] Unretweet a tweet
[-:gray:-] q [-:-:-] Quote tweet
[-:gray:-] r [-:-:-] Reply to
[-:gray:-] D [-:-:-] Delete a tweet
[-:gray:-] o [-:-:-] Open in browser
[-:gray:-] i [-:-:-] Open author's user timeline page
[-:gray:-] c [-:-:-] Copy link to clipboard

[-:-:b]User Navigation
[-:gray:-] w [-:-:-] Follow a user
[-:gray:-] W [-:-:-] Unfollow a user
[-:gray:-] u [-:-:-] Mute a user
[-:gray:-] U [-:-:-] Unmute a user
[-:gray:-] x [-:-:-] Block a user
[-:gray:-] X [-:-:-] Unblock a user
`

	return &cli.Command{
		Name:         "helpshortcuts",
		Short:        "Show help for shortcut keys",
		ValidateFunc: cli.NoArgs(),
		Hidden:       shared.isCommandLineMode,
		RunFunc: func(c *cli.Command, f *pflag.FlagSet) error {
			return a.view.AddPage(newHelpPage("Shortcuts", help), true)
		},
	}
}
