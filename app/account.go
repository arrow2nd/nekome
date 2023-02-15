package app

import (
	"fmt"
	"os"

	"github.com/arrow2nd/nekome/v2/api"
)

// printConsumerKeyWarn : コンシューマーキーに関する警告
func printConsumerKeyWarn() {
	fmt.Fprintln(
		os.Stderr,
		`WARN: The built-in API key may have expired due to the Twitter API no longer being provided free of charge.
      See https://github.com/arrow2nd/nekome for details.`,
	)
}

// addAccount : アカウントを追加
func addAccount(setMain bool) error {
	authUser, err := shared.api.Auth(&shared.conf.Cred.Consumer)
	if err != nil {
		printConsumerKeyWarn()
		return err
	}

	shared.conf.Cred.Write(authUser)

	// メインユーザに設定
	if setMain {
		shared.conf.Pref.Feature.MainUser = authUser.UserName
	}

	return shared.conf.SaveAll()
}

// loginAccount : ログイン
func loginAccount(u string) error {
	// ユーザ名が空なら新規追加
	if u == "" {
		if err := addAccount(true); err != nil {
			return err
		}

		u = shared.conf.Pref.Feature.MainUser
	}

	// ログインするユーザを取得
	user, err := shared.conf.Cred.Get(u)
	if err != nil {
		return err
	}

	api, err := api.New(&shared.conf.Cred.Consumer, user)
	if err != nil {
		return err
	}

	shared.api = api
	return nil
}
