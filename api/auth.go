package api

import (
	"context"
	"fmt"

	"github.com/dghubble/oauth1"
	twauth "github.com/dghubble/oauth1/twitter"
	"github.com/g8rswimmer/go-twitter/v2"
)

var (
	consumerKey    = "mYt6BHZC7gFIgHWLAcFKLKAca"
	consumerSecret = "uUkUPybUlc88IkJWUsd2PCNuW4I8HtSqbRfWNEabX8hqUtUrJg"
)

// Auth : アプリケーション認証を行う
func (a *API) Auth() (*User, error) {
	config := oauth1.Config{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		CallbackURL:    "oob",
		Endpoint:       twauth.AuthorizeEndpoint,
	}

	requestToken, _, err := config.RequestToken()
	if err != nil {
		return nil, err
	}

	authURL, err := config.AuthorizationURL(requestToken)
	if err != nil {
		return nil, err
	}

	fmt.Println("🐈 Go to the following URL to authenticate the application and enter the PIN that is displayed")
	fmt.Println("-----")
	fmt.Println(authURL.String())
	fmt.Print("PIN: ")

	var verifier string

	_, err = fmt.Scanf("%s", &verifier)
	if err != nil {
		return nil, err
	}

	accessToken, accessSecret, err := config.AccessToken(requestToken, "", verifier)
	if err != nil {
		return nil, err
	}

	newToken := oauth1.NewToken(accessToken, accessSecret)

	user, err := a.authUserLookup(newToken)
	if err != nil {
		return nil, err
	}

	return &User{
		UserName: user.UserName,
		ID:       user.ID,
		Token:    newToken,
	}, nil
}

// authUserLookup : トークンに紐づいたユーザの情報を取得
func (a *API) authUserLookup(token *oauth1.Token) (*twitter.UserObj, error) {
	client := newClient(token)

	opts := twitter.UserLookupOpts{}

	res, err := client.AuthUserLookup(context.Background(), opts)
	if e := checkError(err); e != nil {
		return nil, e
	}

	return res.Raw.Users[0], nil
}
