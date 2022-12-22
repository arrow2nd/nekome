package app

import (
	"fmt"

	"github.com/arrow2nd/nekome/v2/cli"
	"github.com/arrow2nd/nekome/v2/log"
	"github.com/spf13/pflag"
)

func newCmd() *cli.Command {
	return &cli.Command{
		Name:  "nekome",
		Short: "TUI Twitter client 🐈",
		Long:  "nekome is a TUI Twitter client that runs on the terminal 🐈",
		SetFlag: func(f *pflag.FlagSet) {
			f.StringP("user", "u", shared.conf.Pref.Feature.MainUser, "specify user to use")
			f.BoolP("version", "v", false, "show version")
		},
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			// コマンドラインからの実行ならバージョンを表示
			if ver, _ := f.GetBool("version"); shared.isCommandLineMode && ver {
				log.Exit(fmt.Sprintf("🐈 nekome for v.%s", version))
			}

			arg := f.Arg(0)
			if arg != "" {
				arg = ": " + arg
			}

			return fmt.Errorf("unavailable or not found command%s", arg)
		},
	}
}

func setUnfocusFlag(f *pflag.FlagSet) {
	f.BoolP("unfocus", "u", false, "no focus on page")
}

func setTimelineFlags(f *pflag.FlagSet) {
	f.BoolP("stream", "s", false, "start stream mode")
	setUnfocusFlag(f)
}
