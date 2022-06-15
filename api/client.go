package api

import (
	"sync"

	"github.com/arrow2nd/nekome/oauth"
	"github.com/g8rswimmer/go-twitter/v2"
	"golang.org/x/oauth2"
)

// User ユーザ情報
type User struct {
	UserName string
	ID       string
	Token    *oauth.Token
}

// API TwitterAPI
type API struct {
	CurrentUser          *User
	oauth                *oauth.OAuth
	tokenRefreshCallback oauth.TokenRefreshCallback
	mu                   sync.Mutex
}

// New 生成
func New() *API {
	config := &oauth2.Config{
		ClientID:    "cmVzRHRHa2haNUlhemJfSFdaM1I6MTpjaQ",
		RedirectURL: "http://localhost:3000/callback",
		Scopes: []string{
			"tweet.read",
			"tweet.write",
			"users.read",
			"like.read",
			"like.write",
			"bookmark.read",
			"bookmark.write",
			"follows.write",
			"block.write",
			"mute.write",
			"list.read",
			"offline.access",
		},
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://api.twitter.com/2/oauth2/token",
			AuthURL:  "https://twitter.com/i/oauth2/authorize",
		},
	}

	return &API{
		CurrentUser: nil,
		oauth:       oauth.New(config),
	}
}

// SetUser ユーザをセット
func (a *API) SetUser(user *User) {
	a.CurrentUser = user
}

// SetToken トークンをセット
func (a *API) SetToken(token *oauth.Token) {
	a.CurrentUser.Token = token
}

// SetTokenRefreshCallback トークンリフレッシュ時のコールバックをセット
func (a *API) SetTokenRefreshCallback(callback oauth.TokenRefreshCallback) {
	a.tokenRefreshCallback = callback
}

func (a *API) newClient(token *oauth.Token) *twitter.Client {
	httpClient := a.oauth.NewClient(token, a.tokenRefreshCallback)

	return &twitter.Client{
		Authorizer: token,
		Client:     httpClient,
		Host:       "https://api.twitter.com",
	}
}
