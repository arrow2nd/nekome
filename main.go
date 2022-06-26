package main

import (
	"log"

	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/app"
	"github.com/arrow2nd/nekome/config"
)

var (
	client *api.API
	conf   *config.Config
)

func init() {
	conf = config.New()
}

func main() {
	// 設定を読込む
	if err := conf.LoadSettings(); err != nil {
		log.Fatal(err)
	}

	// 認証情報を読込む
	ok, err := conf.LoadCred()
	if err != nil {
		log.Fatal(err)
	} else if !ok {
		createNewCred()
	}

	login()

	// アプリ初期化
	app := app.New()
	app.Init(client, conf)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func createNewCred() {
	authUser, err := client.Auth(&conf.Settings.Feature.Client)
	if err != nil {
		log.Fatal(err)
	}

	conf.Cred.Write(authUser)
	conf.Settings.Feature.MainUser = authUser.UserName

	if err := conf.SaveAll(); err != nil {
		log.Fatal(err)
	}
}

func login() error {
	// ログインするユーザを取得
	userName := conf.Settings.Feature.MainUser
	user, err := conf.Cred.Get(userName)
	if err != nil {
		return err
	}

	client = api.New(&conf.Settings.Feature.Client, user)

	return nil
}
