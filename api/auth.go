package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/dghubble/oauth1"
	twauth "github.com/dghubble/oauth1/twitter"
	"github.com/g8rswimmer/go-twitter/v2"
)

// Auth : アプリケーション認証を行なって、ユーザ情報とトークンを取得
func (a *API) Auth(client *oauth1.Token) (*User, error) {
	ct, err := getConsumerToken(client)
	if err != nil {
		return nil, err
	}

	config := oauth1.Config{
		ConsumerKey:    ct.Token,
		ConsumerSecret: ct.TokenSecret,
		CallbackURL:    "oob",
		Endpoint:       twauth.AuthorizeEndpoint,
	}

	// リクエストトークンを取得
	requestToken, _, err := config.RequestToken()
	if err != nil {
		return nil, fmt.Errorf("failed to request token: %w", err)
	}

	// 認証URLを取得
	authURL, err := config.AuthorizationURL(requestToken)
	if err != nil {
		return nil, fmt.Errorf("failed to issue authentication URL: %w", err)
	}

	fmt.Println("🐈 Go to the following URL to authenticate the application and enter the PIN that is displayed")
	fmt.Println()
	fmt.Println(authURL.String())
	fmt.Println()

	// PINの入力受付
	verifier, err := inputPIN()
	if err != nil {
		return nil, fmt.Errorf("failed to read PIN: %w", err)
	}

	// アクセストークン取得
	accessToken, accessSecret, err := config.AccessToken(requestToken, "", verifier)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain token: %w", err)
	}

	newToken := oauth1.NewToken(accessToken, accessSecret)

	// 認証ユーザの詳細を取得
	user, err := a.authUserLookup(client, newToken)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain authenticated user: %w", err)
	}

	return &User{
		UserName: user.UserName,
		ID:       user.ID,
		Token:    newToken,
	}, nil
}

// authUserLookup : トークンに紐づいたユーザの情報を取得
func (a *API) authUserLookup(ct, ut *oauth1.Token) (*twitter.UserObj, error) {
	client, err := newClient(ct, ut)
	if err != nil {
		return nil, err
	}

	opts := twitter.UserLookupOpts{}
	res, err := client.AuthUserLookup(context.Background(), opts)

	if e := checkError(err); e != nil {
		return nil, e
	}

	if res.Raw == nil {
		return nil, errors.New("empty response data")
	}

	if e := checkPartialError(res.Raw.Errors); res.Raw.Users[0] == nil && e != nil {
		return nil, e
	}

	return res.Raw.Users[0], nil
}
