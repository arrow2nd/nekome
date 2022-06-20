package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

const (
	credFileName    = ".cred"
	setingsFileName = "settings.yml"
)

func getConfigDir() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("failed to get config directory")
	}

	return filepath.Join(path, ".nekome"), nil
}

// LoadAll : 一括読込み
func (c *Config) LoadAll() (bool, error) {
	if ok, err := c.hasAllFileExists(); !ok {
		return false, err
	}

	if err := c.load(credFileName, &c.Cred.users); err != nil {
		return false, err
	}

	if err := c.load(setingsFileName, c.Settings); err != nil {
		return false, err
	}

	return true, nil
}

// SaveAll : 一括保存
func (c *Config) SaveAll() error {
	if err := c.SaveCred(); err != nil {
		return err
	}

	if err := c.SaveSettings(); err != nil {
		return err
	}

	return nil
}

func (c *Config) hasAllFileExists() (bool, error) {
	// configディレクトリの存在チェック
	if _, err := os.Stat(c.dirPath); err != nil {
		if err := os.Mkdir(c.dirPath, 0777); err != nil {
			return false, fmt.Errorf("failed to create configuration directory: %v", err)
		}

		return false, nil
	}

	files := []string{
		credFileName,
		setingsFileName,
	}

	// ファイルの存在チェック
	for _, path := range files {
		if _, err := os.Stat(filepath.Join(c.dirPath, path)); err != nil {
			return false, nil
		}
	}

	return true, nil
}

func (c *Config) save(fileName string, in interface{}) error {
	buf, err := yaml.Marshal(in)
	if err != nil {
		return fmt.Errorf("failed to marshal (%s): %v", fileName, err)
	}

	path := filepath.Join(c.dirPath, fileName)

	if err := ioutil.WriteFile(path, buf, os.ModePerm); err != nil {
		return fmt.Errorf("failed to save (%s): %v", path, err)
	}

	return nil
}

func (c *Config) load(fileName string, out interface{}) error {
	path := filepath.Join(c.dirPath, fileName)

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to load (%s): %v", path, err)
	}

	if err := yaml.Unmarshal(buf, out); err != nil {
		return fmt.Errorf("failed to unmarshal (%s): %v", path, err)
	}

	return nil
}
