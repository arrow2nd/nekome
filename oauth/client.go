package oauth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

const (
	clientId   = "cmVzRHRHa2haNUlhemJfSFdaM1I6MTpjaQ"
	listenAddr = "127.0.0.1:3000"
)

type Client struct {
	config   *oauth2.Config
	session  *Session
	response chan *TokenResponse
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
		session:  newSession(),
		response: make(chan *TokenResponse),
	}
}

func (c *Client) Auth() (*TokenResponse, error) {
	// 認可URLを作成
	url := c.buildAuthorizationURL()

	fmt.Printf("Please access the following URL to approve the application\n%s\n", url)

	// サーバを立ててリダイレクトを待機
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", c.handleCallback)

	serve := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	go func() {
		if err := serve.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// 認可フローの終了を待つ
	tokenResponse := <-c.response

	// サーバを閉じる
	if err := serve.Shutdown(context.Background()); err != nil {
		return nil, err
	}

	return tokenResponse, nil
}

func (c *Client) buildAuthorizationURL() string {
	url := c.config.AuthCodeURL(
		c.session.state,
		oauth2.SetAuthURLParam("code_challenge", c.session.codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "s256"),
	)

	return url
}

func (c *Client) handleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != c.session.state {
		log.Println("invalid state")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		log.Println("code not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := c.config.Exchange(
		context.Background(),
		code,
		oauth2.SetAuthURLParam("code_verifier", c.session.codeVerifier),
	)
	if err != nil {
		log.Printf("Failed to obtain token %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenResponse := &TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	c.response <- tokenResponse

	w.Write([]byte("Authentication complete! You may close this page."))
	w.WriteHeader(http.StatusOK)
}
