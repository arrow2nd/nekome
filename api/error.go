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

	// calloutエラーではないならそのまま返す
	tErr := &twitter.ErrorResponse{}
	if !errors.As(err, &tErr) {
		return err
	}

	// レート制限
	if tErr.StatusCode == 429 {
		t := tErr.RateLimit.Reset.Time().Local().Format("15:04:05")
		return fmt.Errorf("Rate limit exceeded (Reset time: %s)", t)
	}

	return fmt.Errorf("server error: %d %s | %s", tErr.StatusCode, tErr.Title, tErr.Detail)
}
