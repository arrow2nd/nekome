package main

import (
	"fmt"
	"log"

	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
	"github.com/arrow2nd/nekome/ui"
)

var (
	client *api.API
	conf   *config.Config
)

func init() {
	conf = config.New()
}

func main() {
	// 設定ファイル読み込み
	ok, err := conf.LoadAll()
	if err != nil {
		log.Fatal(err)
	}

	// 設定ファイルが無い場合,新規作成
	if !ok {
		createNewCred()
	}

	login()

	// NOTE: テスト用
	fmt.Printf("Name: %s / UserName: %s / UserID: %s\n", client.CurrentUser.UserName, client.CurrentUser.UserName, client.CurrentUser.ID)

	// UI初期化
	tui := ui.New()
	tui.Init(client, conf)

	if err := tui.Run(); err != nil {
		log.Fatal(err)
	}
}

func createNewCred() {
	authUser, err := client.Auth()
	if err != nil {
		log.Fatal(err)
	}

	conf.Cred.Write(authUser)
	conf.Settings.MainUser = authUser.UserName

	if err := conf.SaveAll(); err != nil {
		log.Fatal(err)
	}
}

func login() error {
	// ログインするユーザを取得
	userName := conf.Settings.MainUser
	user, err := conf.Cred.Get(userName)
	if err != nil {
		return err
	}

	client = api.New(user)

	return nil
}
