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

	return nil, fmt.Errorf("user \"%s\" does not exist", userName)
}

// Write : 書込む
func (c *Cred) Write(newUser *api.User) {
	// 同じIDをもつユーザが居れば上書きする
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
func (c *Cred) Delete(userName string) {
	var tmp []api.User

	for _, user := range c.users {
		if user.UserName != userName {
			tmp = append(tmp, user)
		}
	}

	c.users = tmp
}
