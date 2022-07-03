package config

import (
	"fmt"

	"github.com/arrow2nd/nekome/api"
)

// Cred : 認証情報
type Cred struct {
	users []api.User
}

// Get : 取得
func (c *Cred) Get(userName string) (*api.User, error) {
	for _, user := range c.users {
		if user.UserName == userName {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user %s does not exist", userName)
}

// GetAllNames : 全てのユーザ名を取得
func (c *Cred) GetAllNames() []string {
	ls := []string{}

	for _, u := range c.users {
		ls = append(ls, u.UserName)
	}

	return ls
}

// Write : 書込む
func (c *Cred) Write(newUser *api.User) {
	// 同じIDをもつユーザが居れば上書き
	for i, user := range c.users {
		if user.ID == newUser.ID {
			c.users[i] = *newUser
			return
		}
	}

	// 新規追加
	c.users = append(c.users, *newUser)
}

// Delete : 削除
func (c *Cred) Delete(userName string) error {
	var err error = nil
	tmp := []api.User{}

	for _, user := range c.users {
		if user.UserName == userName {
			err = fmt.Errorf("user %s does not exist", userName)
			continue
		}

		tmp = append(tmp, user)
	}

	if err == nil {
		c.users = tmp
	}

	return err
}
