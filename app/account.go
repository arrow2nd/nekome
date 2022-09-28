package app

import "github.com/arrow2nd/nekome/api"

// addAccount : アカウントを追加
func addAccount(setMain bool) error {
	// 認証
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

// loginAccount : ログイン処理
func loginAccount(u string) error {
	// ログインするユーザを取得
	user, err := shared.conf.Cred.Get(u)
	if err != nil {
		return err
	}

	// 新しいユーザでクライアントを生成
	api, err := api.New(&shared.conf.Cred.Consumer, user)
	if err != nil {
		return err
	}

	shared.api = api

	return nil
}
