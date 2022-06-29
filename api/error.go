package api

import (
	"errors"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

// checkError : エラーを解析して適切なエラーメッセージを返す
func checkError(err error) error {
	// HTTPエラー
	httpErr := &twitter.HTTPError{}
	if errors.As(err, &httpErr) {
		return fmt.Errorf("http error: %s", httpErr.Status)
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

// checkPartialError : 部分エラーが無いかチェック
func checkPartialError(errs []*twitter.ErrorObj) error {
	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("partial error: %s %s", errs[0].Title, errs[0].Detail)
}
