package log

import (
	"fmt"
	"os"
)

// ExitCode : 終了コード
type ExitCode int

// GetInt : 数値を取得
func (e ExitCode) GetInt() int {
	return int(e)
}

const (
	// ExitCodeOK : 正常
	ExitCodeOK ExitCode = iota
	// ExitCodeErrAuth : 認証エラー
	ExitCodeErrAuth
	// ExitCodeErrFileIO : ファイルIOエラー
	ExitCodeErrFileIO
	// ExitCodeErrApp : アプリケーションエラー
	ExitCodeErrApp
)

// LogExit : ログを出力して終了
func LogExit(s string) {
	fmt.Println(s)
	os.Exit(ExitCodeOK.GetInt())
}

// ErrorExit : エラーを出力して終了
func ErrorExit(e string, c ExitCode) {
	fmt.Fprintln(os.Stderr, e)
	os.Exit(c.GetInt())
}
