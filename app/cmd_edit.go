package app

import (
	"os"
	"path"

	"github.com/arrow2nd/nekome/v2/cli"
	"github.com/arrow2nd/nekome/v2/config"
	"github.com/manifoldco/promptui"
	"github.com/spf13/pflag"
)

func (a *App) newEditCmd() *cli.Command {
	return &cli.Command{
		Name:      "edit",
		Shorthand: "e",
		Short:     "Edit configuration file",
		Hidden:    !shared.isCommandLineMode,
		Validate:  cli.NoArgs(),
		SetFlag: func(f *pflag.FlagSet) {
			f.StringP("editor", "e", os.Getenv("EDITOR"), "specify which editor to use (default is $EDITOR)")
		},
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			// 選択肢に表示するファイルを取得
			items, err := config.GetConfigFileNames()
			if err != nil {
				return err
			}

			prompt := promptui.Select{
				Label: "File to edit",
				Items: items,
			}

			_, file, err := prompt.Run()
			if err != nil {
				return err
			}

			dir := shared.conf.DirPath
			editor, _ := f.GetString("editor")
			return a.openExternalEditor(editor, path.Join(dir, file))
		},
	}
}
