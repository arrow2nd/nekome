package main

import (
	"github.com/arrow2nd/nekome/v2/app"
	"github.com/arrow2nd/nekome/v2/log"
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
