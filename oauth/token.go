package oauth

import "golang.org/x/oauth2"

type TokenRefreshCallback func(*oauth2.Token) error

type TokenSource struct {
	src      oauth2.TokenSource
	callback TokenRefreshCallback
}

func (ts *TokenSource) Token() (*oauth2.Token, error) {
	t, err := ts.src.Token()
	if err != nil {
		return nil, err
	}

	if err := ts.callback(t); err != nil {
		return t, err
	}

	return t, nil
}
