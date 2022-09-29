package config

import (
	"fmt"

	"github.com/arrow2nd/nekome/api"
	"github.com/dghubble/oauth1"
)

// User : ユーザの認証情報
type User struct {
	Accounts []api.User `toml:"accounts"`
}

// Cred : 認証情報
type Cred struct {
	Consumer oauth1.Token `toml:"consumer"`
	User     User         `toml:"user"`
}

// Get : 取得
func (c *Cred) Get(userName string) (*api.User, error) {
	for _, user := range c.User.Accounts {
		if user.UserName == userName {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user not found: %s", userName)
}

// GetAllNames : 全てのユーザ名を取得
func (c *Cred) GetAllNames() []string {
	ls := []string{}

	for _, u := range c.User.Accounts {
		ls = append(ls, u.UserName)
	}

	return ls
}

// Write : 書込む
func (c *Cred) Write(newUser *api.User) {
	// 同じIDをもつユーザが居れば上書き
	for i, user := range c.User.Accounts {
		if user.ID == newUser.ID {
			c.User.Accounts[i] = *newUser
			return
		}
	}

	// 新規追加
	c.User.Accounts = append(c.User.Accounts, *newUser)
}

// Delete : 削除
func (c *Cred) Delete(userName string) error {
	err := fmt.Errorf("user not found: %s", userName)
	tmp := []api.User{}

	for _, user := range c.User.Accounts {
		if user.UserName == userName {
			err = nil
			continue
		}

		tmp = append(tmp, user)
	}

	if err == nil {
		c.User.Accounts = tmp
	}

	return err
}
