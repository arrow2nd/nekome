package oauth

import (
	"fmt"

	"golang.org/x/oauth2"
)

const (
	clientId = "cmVzRHRHa2haNUlhemJfSFdaM1I6MTpjaQ"
)

type Client struct {
	config  *oauth2.Config
	session *Session
}

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
}

func New() *Client {
	return &Client{
		config: &oauth2.Config{
			ClientID:    clientId,
			RedirectURL: "http://localhost:3000/callback",
			Scopes:      []string{"tweet.read", "users.read", "offline.access"},
			Endpoint: oauth2.Endpoint{
				TokenURL: "https://api.twitter.com/2/oauth2/token",
				AuthURL:  "https://twitter.com/i/oauth2/authorize",
			},
		},
		session: newSession(),
	}
}

func (oa *Client) Auth() *TokenResponse {
	url := oa.buildAuthorizationURL()

	fmt.Printf("Please access the following URL to approve the application\n%s\n", url)

	return &TokenResponse{}
}

func (oa *Client) buildAuthorizationURL() string {
	url := oa.config.AuthCodeURL(
		oa.session.state,
		oauth2.SetAuthURLParam("code_challenge", oa.session.codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "s256"),
	)

	return url
}
