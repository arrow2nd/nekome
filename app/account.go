package app

import "github.com/arrow2nd/nekome/v2/api"

// addAccount : アカウントを追加
func addAccount(setMain bool) error {
	authUser, err := shared.api.Auth(&shared.conf.Cred.Consumer)
	if err != nil {
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
