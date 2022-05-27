package api

import (
	"github.com/arrow2nd/nekome/oauth"
	"github.com/g8rswimmer/go-twitter/v2"
	"golang.org/x/oauth2"
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
	config := &oauth2.Config{
		ClientID:    "cmVzRHRHa2haNUlhemJfSFdaM1I6MTpjaQ",
		RedirectURL: "http://localhost:3000/callback",
		Scopes:      []string{"tweet.read", "users.read", "offline.access"},
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://api.twitter.com/2/oauth2/token",
			AuthURL:  "https://twitter.com/i/oauth2/authorize",
		},
	}

	return &API{
		CurrentUser: nil,
		oauth:       oauth.New(config),
		token:       nil,
	}
}

// SetUser ユーザをセット
func (a *API) SetUser(token *oauth.Token) error {
	user, err := a.authUserLookup(token)
	if err != nil {
		return err
	}

	a.token = token
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
