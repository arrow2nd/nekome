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
	ok, err := conf.LoadAll()
	if err != nil {
		log.Fatal(err)
	}

	// 初期設定
	if !ok {
		createNewConfig()
	}

	login()

	// NOTE: テスト用
	fmt.Printf("Name: %s / UserName: %s / UserID: %s\n", client.CurrentUser.UserName, client.CurrentUser.UserName, client.CurrentUser.ID)

	tui := ui.New()
	tui.Init(client, conf)

	if err := tui.App.Run(); err != nil {
		log.Fatal(err)
	}
}

func createNewConfig() {
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

	// トークンを更新
	client.SetToken(token)
	conf.Cred.Write(client.CurrentUser)

	return conf.SaveCred()
}
