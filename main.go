package main

import (
	"github.com/arrow2nd/nekome/app"
	"github.com/arrow2nd/nekome/log"
)

func main() {
	app := app.New()

	if err := app.Init(); err != nil {
		log.ErrorExit(err.Error(), log.ExitCodeErrInit)
	}

	if err := app.Run(); err != nil {
		log.ErrorExit(err.Error(), log.ExitCodeErrApp)
	}
}
