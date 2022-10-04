package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	credFileName = ".cred.toml"
	prefFileName = "preferences.toml"
)

// GetConfigDir : 設定ディレクトリを取得
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("failed to get home directory")
	}

	homeDir = filepath.Join(homeDir, ".config", "nekome")

	// ディレクトリが無いなら作成
	if _, err := os.Stat(homeDir); err != nil {
		if err := os.MkdirAll(homeDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	return homeDir, nil
}

// GetConfigFileNames : 設定ディレクトリ以下のファイル名を取得
func GetConfigFileNames() ([]string, error) {
	path, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	fileNames := []string{}
	for _, e := range entries {
		fileNames = append(fileNames, e.Name())
	}

	return fileNames, nil
}

// CheckOldFile : 古い設定ファイルがあるかチェック
func (c *Config) CheckOldFile() {
	notice := `[ 🐈 Notice ]

Starting with nekome v2.0.0, the configuration file format has been changed from yaml to toml.
Please run 'nekome edit' and reconfigure the file.

For more information, please refer to the following pages.
https://github.com/arrow2nd/nekome/blob/v2/docs/en/migrate-v1-v2.md

(This notice will not appear if you delete the old configuration file)
`

	if c.hasFileExists("default.yml") || c.hasFileExists("settings.yml") {
		fmt.Println(notice)
	}
}

// LoadCred : 認証情報を読込む
func (c *Config) LoadCred() (bool, error) {
	if !c.hasFileExists(credFileName) {
		return false, c.SaveCred()
	}

	if err := c.load(credFileName, &c.Cred); err != nil {
		return false, err
	}

	return len(c.Cred.User.Accounts) > 0, nil
}

// LoadPreferences : 環境設定を読込む
func (c *Config) LoadPreferences() error {
	if !c.hasFileExists(prefFileName) {
		if err := c.SavePreferences(); err != nil {
			return err
		}
	}

	return c.load(prefFileName, c.Pref)
}

// LoadStyle : スタイル定義を読込む
func (c *Config) LoadStyle() error {
	fileName := c.Pref.Appearance.StyleFilePath

	if !c.hasFileExists(fileName) {
		if err := c.saveDefaultStyle(); err != nil {
			return err
		}
	}

	return c.load(fileName, c.Style)
}

// SaveCred : 認証情報を保存
func (c *Config) SaveCred() error {
	return c.save(credFileName, c.Cred)
}

// SavePreferences : 環境設定を保存
func (c *Config) SavePreferences() error {
	return c.save(prefFileName, c.Pref)
}

// saveDefaultStyle : デフォルトのスタイル定義を保存
func (c *Config) saveDefaultStyle() error {
	return c.save(c.Pref.Appearance.StyleFilePath, c.Style)
}

// SaveAll : 一括保存
func (c *Config) SaveAll() error {
	if err := c.SaveCred(); err != nil {
		return err
	}

	if err := c.SavePreferences(); err != nil {
		return err
	}

	return nil
}

// hasFileExists : ファイルが存在するか
func (c *Config) hasFileExists(file string) bool {
	if _, err := os.Stat(filepath.Join(c.dirPath, file)); err != nil {
		return false
	}

	return true
}

// save : 保存
func (c *Config) save(fileName string, in interface{}) error {
	buf := &bytes.Buffer{}

	if err := toml.NewEncoder(buf).Encode(in); err != nil {
		return fmt.Errorf("failed to marshal (%s): %w", fileName, err)
	}

	path := filepath.Join(c.dirPath, fileName)

	if err := ioutil.WriteFile(path, buf.Bytes(), os.ModePerm); err != nil {
		return fmt.Errorf("failed to save (%s): %w", path, err)
	}

	return nil
}

// load : 読み込み
func (c *Config) load(fileName string, out interface{}) error {
	path := filepath.Join(c.dirPath, fileName)

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to load (%s): %w", path, err)
	}

	if err := toml.Unmarshal(buf, out); err != nil {
		return fmt.Errorf("failed to unmarshal (%s): %w", path, err)
	}

	return nil
}
