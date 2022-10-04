package config

import (
	"github.com/arrow2nd/nekome/v2/log"
)

// Config : 設定
type Config struct {
	// Cred : 認証情報
	Cred *Cred
	// Pref : 環境設定
	Pref *Preferences
	// Style : スタイル定義
	Style *Style

	// dirPath : 設定ディレクトリのパス
	dirPath string
}

// New : 新規作成
func New() *Config {
	path, err := GetConfigDir()
	if err != nil {
		log.ErrorExit(err.Error(), log.ExitCodeErrFileIO)
	}

	return &Config{
		Cred:    &Cred{},
		Pref:    defaultPreferences(),
		Style:   defaultStyle(),
		dirPath: path,
	}
}
