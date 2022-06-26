package api

import (
	"net/http"

	"github.com/dghubble/oauth1"
	"github.com/g8rswimmer/go-twitter/v2"
)

// NOTE: 認証に必要なヘッダは oauth1 の Client() で設定してくれるので必要ないが
//       Add() は go-twitter で必ず呼ばれるのでダミーとして用意してる
type authorizer struct{}

func (a *authorizer) Add(req *http.Request) {}

// User : ユーザ情報
type User struct {
	UserName string
	ID       string
	Token    *oauth1.Token
}

// API : TwitterAPI
type API struct {
	CurrentUser *User
	client      *twitter.Client
}

// New : 作成
func New(client *oauth1.Token, user *User) *API {
	return &API{
		CurrentUser: user,
		client:      newClient(client, user.Token),
	}
}

func newClient(ct, ut *oauth1.Token) *twitter.Client {
	ct = getClientToken(ct)
	config := oauth1.NewConfig(ct.Token, ct.TokenSecret)
	httpClient := config.Client(oauth1.NoContext, ut)

	return &twitter.Client{
		Authorizer: &authorizer{},
		Client:     httpClient,
		Host:       "https://api.twitter.com",
	}
}
