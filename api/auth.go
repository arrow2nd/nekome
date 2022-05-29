package api

import (
	"context"
	"fmt"

	"github.com/arrow2nd/nekome/oauth"
	"github.com/g8rswimmer/go-twitter/v2"
)

// Auth アプリケーション認証を行う
func (a *API) Auth() (*User, error) {
	token, err := a.oauth.Auth()
	if err != nil {
		return nil, err
	}

	// 認証したユーザの情報を取得
	authUser, err := a.authUserLookup(token)
	if err != nil {
		return nil, err
	}

	user := &User{
		UserName: authUser.UserName,
		ID:       authUser.ID,
		Token:    token,
	}

	return user, nil
}

// authUserLookup トークンに紐づいたユーザの情報を取得
func (a *API) authUserLookup(token *oauth.Token) (*twitter.UserObj, error) {
	client := a.newClient(token)

	opts := twitter.UserLookupOpts{}
	userResponse, err := client.AuthUserLookup(context.Background(), opts)
	if err != nil {
		return nil, fmt.Errorf("auth user lookup error: %v", err)
	}

	return userResponse.Raw.Users[0], nil
}
