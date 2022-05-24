package oauth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const (
	clientId   = "cmVzRHRHa2haNUlhemJfSFdaM1I6MTpjaQ"
	listenAddr = "127.0.0.1:3000"
)

type Client struct {
	config   *oauth2.Config
	session  *Session
	response chan *Token
}

type Token struct {
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
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
		response: make(chan *Token),
	}
}

func (c *Client) Auth() (*Token, error) {
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
	c.session = newSession()

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

	tokenResponse := &Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}

	c.response <- tokenResponse

	w.Write([]byte("Authentication complete! You may close this page."))
	w.WriteHeader(http.StatusOK)
}

type TokenRefreshFunc func(*oauth2.Token) error

type TokenSource struct {
	src oauth2.TokenSource
	f   TokenRefreshFunc
}

func (s *TokenSource) Token() (*oauth2.Token, error) {
	t, err := s.src.Token()
	if err != nil {
		return nil, err
	}

	if err := s.f(t); err != nil {
		return t, err
	}

	return t, nil
}

func (c *Client) NewClient(token *Token, refreshFunc TokenRefreshFunc) *http.Client {
	t := &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    "bearer",
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}

	src := c.config.TokenSource(context.Background(), t)

	tokenSource := &TokenSource{
		src: src,
		f:   refreshFunc,
	}

	reuseSrc := oauth2.ReuseTokenSource(t, tokenSource)

	return oauth2.NewClient(context.Background(), reuseSrc)
}
