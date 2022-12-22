package app

import (
	"github.com/arrow2nd/nekome/v2/cli"
	"github.com/spf13/pflag"
)

func (a *App) newHomeCmd() *cli.Command {
	return &cli.Command{
		Name:      "home",
		Shorthand: "h",
		Short:     "Add home timeline page",
		Validate:  cli.NoArgs(),
		Hidden:    shared.isCommandLineMode,
		SetFlag:   setTimelineFlags,
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			page, err := newTimelinePage(timelineTypeHome)
			if err != nil {
				return err
			}

			unfocus, _ := f.GetBool("unfocus")
			stream, _ := f.GetBool("stream")

			if err := a.view.AddPage(page, !unfocus); err != nil {
				return err
			}

			go func() {
				page.Load()
				if stream {
					page.startStreamMode()
				}
			}()

			return nil
		},
	}
}
