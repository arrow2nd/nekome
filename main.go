package main

import (
	"fmt"
	"os"

	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/app"
	"github.com/arrow2nd/nekome/config"
)

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeErrAuth
	exitCodeErrConfig
	exitCodeErrCred
	exitCodeErrFileIO
	exitCodeErrApp
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
		printErrorExit(err, exitCodeErrConfig)
	}

	// 認証情報を読込む
	ok, err := conf.LoadCred()
	if err != nil {
		printErrorExit(err, exitCodeErrCred)
	} else if !ok {
		createNewCred()
	}

	login()

	// アプリ初期化
	app := app.New()
	app.Init(client, conf)

	if err := app.Run(); err != nil {
		printErrorExit(err, exitCodeErrApp)
	}
}

func createNewCred() {
	authUser, err := client.Auth(&conf.Settings.Feature.Consumer)
	if err != nil {
		printErrorExit(err, exitCodeErrAuth)
	}

	conf.Cred.Write(authUser)
	conf.Settings.Feature.MainUser = authUser.UserName

	if err := conf.SaveAll(); err != nil {
		printErrorExit(err, exitCodeErrFileIO)
	}
}

func login() error {
	// ログインするユーザを取得
	userName := conf.Settings.Feature.MainUser
	user, err := conf.Cred.Get(userName)
	if err != nil {
		return err
	}

	client = api.New(&conf.Settings.Feature.Consumer, user)

	return nil
}

// printErrorExit : エラーを出力して終了
func printErrorExit(e error, c exitCode) {
	fmt.Fprintln(os.Stderr, e.Error())
	os.Exit(int(c))
}
