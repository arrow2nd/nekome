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

	errText := fmt.Sprintf("%d %s | %s", tErr.StatusCode, tErr.Title, tErr.Detail)

	// 認証
	if tErr.StatusCode == 401 || tErr.StatusCode == 402 || tErr.StatusCode == 403 {
		return fmt.Errorf("client error: %s (API key may no longer be available)", errText)
	}

	// レート制限
	if tErr.StatusCode == 429 {
		t := tErr.RateLimit.Reset.Time().Local().Format("15:04:05")
		return fmt.Errorf("Rate limit exceeded (Reset time: %s)", t)
	}

	return fmt.Errorf("server error: %s", errText)
}

// checkPartialError : 部分エラーが無いかチェック
func checkPartialError(errs []*twitter.ErrorObj) error {
	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("partial error: %s %s", errs[0].Title, errs[0].Detail)
}
