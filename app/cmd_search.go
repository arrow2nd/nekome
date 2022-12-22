package app

import (
	"errors"
	"strings"

	"github.com/arrow2nd/nekome/v2/cli"
	"github.com/spf13/pflag"
)

func (a *App) newSearchCmd() *cli.Command {
	return &cli.Command{
		Name:      "search",
		Shorthand: "s",
		Short:     "Add seaech result page",
		UsageArgs: "<query>",
		Example:   "search golang",
		Hidden:    shared.isCommandLineMode,
		SetFlag:   setUnfocusFlag,
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			query := strings.Join(f.Args(), " ")
			if query == "" {
				return errors.New("please specify search query")
			}

			page, err := newSearchPage(query)
			if err != nil {
				return err
			}

			unfocus, _ := f.GetBool("unfocus")
			if err := a.view.AddPage(page, !unfocus); err != nil {
				return err
			}

			go page.Load()

			return nil
		},
	}
}
