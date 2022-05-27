package api

import (
	"github.com/arrow2nd/nekome/oauth"
	"github.com/g8rswimmer/go-twitter/v2"
)

// API TwitterAPI
type API struct {
	CurrentUserName      string
	oauth                *oauth.OAuth
	token                *oauth.Token
	tokenRefreshCallback oauth.TokenRefreshCallback
}

// New 生成
func New() *API {
	return &API{
		CurrentUserName: "",
		oauth:           oauth.New(),
		token:           nil,
	}
}

// SetUser ユーザを設定
func (a *API) SetUser(userName string, token *oauth.Token) {
	a.CurrentUserName = userName
	a.token = token
}

// SetTokenRefreshCallback トークンリフレッシュ時のコールバックを設定
func (a *API) SetTokenRefreshCallback(callback oauth.TokenRefreshCallback) {
	a.tokenRefreshCallback = callback
}

func (a *API) newClient(token *oauth.Token) (*twitter.Client, error) {
	httpClient := a.oauth.NewClient(token, a.tokenRefreshCallback)

	client := &twitter.Client{
		Authorizer: token,
		Client:     httpClient,
		Host:       "https://api.twitter.com",
	}

	return client, nil
}
