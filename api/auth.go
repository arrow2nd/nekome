package api

import (
	"context"

	"github.com/arrow2nd/nekome/oauth"
	"github.com/g8rswimmer/go-twitter/v2"
)

// Auth アプリケーション認証を行う
func (a *API) Auth() (string, *oauth.Token, error) {
	token, err := a.oauth.Auth()
	if err != nil {
		return "", nil, err
	}

	// 認証したユーザの情報を取得
	user, err := a.AuthUserLookupByToken(token)
	if err != nil {
		return "", nil, err
	}

	return user.UserName, token, nil
}

// AuthUserLookup 現在のユーザの情報を取得
func (a *API) AuthUserLookup() (*twitter.UserObj, error) {
	return a.AuthUserLookupByToken(a.token)
}

// AuthUserLookupByToken トークンに紐づいたユーザの情報を取得
func (a *API) AuthUserLookupByToken(token *oauth.Token) (*twitter.UserObj, error) {
	client, err := a.newClient(token)
	if err != nil {
		return nil, err
	}

	opts := twitter.UserLookupOpts{}
	userResponse, err := client.AuthUserLookup(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	return userResponse.Raw.Users[0], nil
}
