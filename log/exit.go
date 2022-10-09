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
	// ExitCodeErrInit : 初期化エラー
	ExitCodeErrInit
	// ExitCodeErrApp : アプリケーションエラー
	ExitCodeErrApp
	// ExitCodeErrFileIO : ファイルIOエラー
	ExitCodeErrFileIO
	// ExitCodeErrTerm : 端末関連のエラー
	ExitCodeErrTerm
)

// Exit : ログを出力して終了
func Exit(s string) {
	fmt.Println(s)
	os.Exit(ExitCodeOK.GetInt())
}

// ErrorExit : エラーを出力して終了
func ErrorExit(e string, c ExitCode) {
	fmt.Fprintf(os.Stderr, "[Error] %s\n", e)
	os.Exit(c.GetInt())
}
