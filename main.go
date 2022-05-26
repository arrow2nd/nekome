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

	if !ok {
		initConfig()
	}

	userName := conf.Settings.MainUser
	token, err := conf.Cred.Get(userName)
	if err != nil {
		log.Fatal(err)
	}

	client.SetUser(userName, token)
	client.SetTokenRefreshCallback(handleTokenRefresh)

	// NOTE: テスト用
	user, err := client.AuthUserLookup()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Name: %s / UserName: %s\n", user.Name, user.UserName)
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
	userName := client.UserName
	token := &oauth.Token{
		AccessToken:  rawToken.AccessToken,
		RefreshToken: rawToken.RefreshToken,
		Expiry:       rawToken.Expiry,
	}

	client.SetUser(userName, token)

	conf.Cred.Write(userName, token)
	conf.SaveCred()

	return nil
}
