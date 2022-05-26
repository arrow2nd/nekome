package api

import (
	"github.com/arrow2nd/nekome/oauth"
	"github.com/g8rswimmer/go-twitter/v2"
)

// API TwitterAPI
type API struct {
	UserName             string
	oauth                *oauth.OAuth
	token                *oauth.Token
	tokenRefreshCallback oauth.TokenRefreshCallback
}

// New 生成
func New() *API {
	return &API{
		UserName: "",
		oauth:    oauth.New(),
		token:    nil,
	}
}

// SetUser ユーザを設定
func (a *API) SetUser(userName string, token *oauth.Token) {
	a.UserName = userName
	a.token = token
}

// SetTokenRefreshCallback トークンリフレッシュ時のコールバックを設定
func (a *API) SetTokenRefreshCallback(callback oauth.TokenRefreshCallback) {
	a.tokenRefreshCallback = callback
}

// Auth アプリケーション認証を行う
func (a *API) Auth() (*oauth.Token, error) {
	return a.oauth.Auth()
}

func (a *API) newClient() (*twitter.Client, error) {
	// expiry, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", "2022-05-25 11:53:32.258818 +0900 JST m=+7210.477936373")

	httpClient := a.oauth.NewClient(a.token, a.tokenRefreshCallback)

	client := &twitter.Client{
		Authorizer: a.token,
		Client:     httpClient,
		Host:       "https://api.twitter.com",
	}

	return client, nil
}
