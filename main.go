package main

import (
	"fmt"
	"log"

	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
	"github.com/arrow2nd/nekome/oauth"
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
		initConfig()
	}

	// 使用するトークンを取得
	userName := conf.Settings.MainUser
	token, err := conf.Cred.Get(userName)
	if err != nil {
		log.Fatal(err)
	}

	// クライアントにユーザをセット
	client.SetTokenRefreshCallback(handleTokenRefresh)
	if err := client.SetUser(token); err != nil {
		log.Fatal(err)
	}

	// NOTE: テスト用
	fmt.Printf("Name: %s / UserName: %s / UserID: %s\n", client.CurrentUser.Name, client.CurrentUser.UserName, client.CurrentUser.ID)

	tweets, err := client.UserTimeline(client.CurrentUser.ID)
	if err != nil {
		log.Fatal(err)
	}

	for i, tweet := range tweets {
		fmt.Printf("[%d] %s : %s\n", i, tweet.AuthorID, tweet.Text)
	}
}

func initConfig() {
	userName, token, err := client.Auth()
	if err != nil {
		log.Fatal(err)
	}

	conf.Cred.Write(userName, token)
	conf.Settings.MainUser = userName

	if err := conf.SaveAll(); err != nil {
		log.Fatal(err)
	}
}

func handleTokenRefresh(rawToken *oauth2.Token) error {
	prevUserName := client.CurrentUser.UserName
	token := &oauth.Token{
		AccessToken:  rawToken.AccessToken,
		RefreshToken: rawToken.RefreshToken,
		Expiry:       rawToken.Expiry,
	}

	// トークンを更新
	if err := client.SetUser(token); err != nil {
		return err
	}

	conf.Cred.Write(client.CurrentUser.UserName, token)

	// UserNameが変更されていたら、前のユーザデータを削除
	if client.CurrentUser.UserName != prevUserName {
		conf.Cred.Delete(prevUserName)
	}

	return conf.SaveCred()
}
