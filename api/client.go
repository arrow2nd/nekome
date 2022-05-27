package api

import (
	"github.com/arrow2nd/nekome/oauth"
	"github.com/g8rswimmer/go-twitter/v2"
)

// API TwitterAPI
type API struct {
	CurrentUser          *twitter.UserObj
	oauth                *oauth.OAuth
	token                *oauth.Token
	tokenRefreshCallback oauth.TokenRefreshCallback
}

// New 生成
func New() *API {
	return &API{
		CurrentUser: nil,
		oauth:       oauth.New(),
		token:       nil,
	}
}

// SetUser ユーザをセット
func (a *API) SetUser(token *oauth.Token) error {
	a.token = token

	user, err := a.AuthUserLookup()
	if err != nil {
		return err
	}

	a.CurrentUser = user

	return nil
}

// SetTokenRefreshCallback トークンリフレッシュ時のコールバックをセット
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
