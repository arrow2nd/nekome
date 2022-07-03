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

// func createNewCred() {
// 	authUser, err := client.Auth(&conf.Settings.Feature.Consumer)
// 	if err != nil {
// 		log.ErrorExit(err.Error(), log.ExitCodeErrAuth)
// 	}

// 	conf.Cred.Write(authUser)
// 	conf.Settings.Feature.MainUser = authUser.UserName

// 	if err := conf.SaveAll(); err != nil {
// 		log.ErrorExit(err.Error(), log.ExitCodeErrFileIO)
// 	}
// }

// func login() error {
// 	// ログインするユーザを取得
// 	userName := conf.Settings.Feature.MainUser
// 	user, err := conf.Cred.Get(userName)
// 	if err != nil {
// 		return err
// 	}

// 	client = api.New(&conf.Settings.Feature.Consumer, user)

// 	return nil
// }
