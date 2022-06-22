package api

import (
	"errors"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

func checkError(err error) error {
	// HTTPエラー
	httpErr := &twitter.HTTPError{}
	if errors.As(err, &httpErr) {
		return fmt.Errorf("http error: %d %s", httpErr.StatusCode, httpErr.Status)
	}

	// calloutエラー
	tErr := &twitter.ErrorResponse{}
	if errors.As(err, &tErr) {
		return fmt.Errorf("server error: %d %s %s", tErr.StatusCode, tErr.Title, tErr.Detail)
	}

	// レート制限
	if rateLimit, has := twitter.RateLimitFromError(err); has {
		t := rateLimit.Reset.Time().Local().Format("15:04:05")
		return fmt.Errorf("Rate limit exceeded (Reset time: %s)", t)
	}

	return err
}
