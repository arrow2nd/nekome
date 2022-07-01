package config

import (
	"github.com/arrow2nd/nekome/log"
)

// Config : 設定
type Config struct {
	// Cred : 認証情報
	Cred *Cred
	// Settings : 設定情報
	Settings *Settings
	dirPath  string
}

// New : 生成
func New() *Config {
	path, err := GetConfigDir()
	if err != nil {
		log.ErrorExit(err.Error(), log.ExitCodeErrFileIO)
	}

	return &Config{
		Cred:     &Cred{},
		Settings: defaultSettings(),
		dirPath:  path,
	}
}
