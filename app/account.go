package app

import "github.com/arrow2nd/nekome/api"

// addAccount : アカウントを追加
func addAccount(setMain bool) error {
	// 認証
	authUser, err := shared.api.Auth(&shared.conf.Settings.Feature.Consumer)
	if err != nil {
		return err
	}
	shared.conf.Cred.Write(authUser)

	// メインユーザに設定
	if setMain {
		shared.conf.Settings.Feature.MainUser = authUser.UserName
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

	// 新しいユーザで初期化
	shared.api = api.New(&shared.conf.Settings.Feature.Consumer, user)
	return nil
}
