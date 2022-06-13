package main

import (
	"fmt"
	"log"

	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
	"github.com/arrow2nd/nekome/oauth"
	"github.com/arrow2nd/nekome/ui"
	"golang.org/x/oauth2"
)

var (
	client *api.API
	conf   *config.Config
)

func init() {
	client = api.New()
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
	// 使用するトークンを取得
	userName := conf.Settings.MainUser
	user, err := conf.Cred.Get(userName)
	if err != nil {
		return err
	}

	// クライアントを初期化
	client.SetUser(user)
	client.SetTokenRefreshCallback(handleTokenRefresh)

	return nil
}

func handleTokenRefresh(rawToken *oauth2.Token) error {
	token := &oauth.Token{
		AccessToken:  rawToken.AccessToken,
		RefreshToken: rawToken.RefreshToken,
		Expiry:       rawToken.Expiry,
	}

	// クライアントと設定ファイルのトークンを更新
	client.SetToken(token)
	conf.Cred.Write(client.CurrentUser)

	return conf.SaveCred()
}
