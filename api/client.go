package api

import (
	"net/http"

	"github.com/dghubble/oauth1"
	"github.com/g8rswimmer/go-twitter/v2"
)

// NOTE: 認証に必要なヘッダは oauth1 の Client() 内で設定してくれるので必要ないが
//       go-twitter で Add() が呼ばれるのでダミーとして用意してる

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

// New : 新規作成
func New(ct *oauth1.Token, u *User) (*API, error) {
	client, err := newClient(ct, u.Token)
	if err != nil {
		return nil, err
	}

	return &API{
		CurrentUser: u,
		client:      client,
	}, nil
}

func newClient(ct, ut *oauth1.Token) (*twitter.Client, error) {
	ct, err := getConsumerToken(ct)
	if err != nil {
		return nil, err
	}

	config := oauth1.NewConfig(ct.Token, ct.TokenSecret)
	httpClient := config.Client(oauth1.NoContext, ut)

	return &twitter.Client{
		Authorizer: &authorizer{},
		Client:     httpClient,
		Host:       "https://api.twitter.com",
	}, nil
}
