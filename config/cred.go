package config

import (
	"fmt"

	"github.com/arrow2nd/nekome/oauth"
)

// Cred 認証情報
type Cred struct {
	tokens map[string]*oauth.Token
}

// Get ユーザ名から認証情報を取得
func (c *Cred) Get(userName string) (*oauth.Token, error) {
	if token, ok := c.tokens[userName]; ok {
		return token, nil
	}

	return nil, fmt.Errorf("user \"%s\" does not exist", userName)
}

// Write 認証情報を書込む
func (c *Cred) Write(userName string, token *oauth.Token) {
	c.tokens[userName] = token
}

// Delete 認証情報を削除
func (c *Cred) Delete(userName string) {
	delete(c.tokens, userName)
}

// SaveCred 認証情報を保存
func (c *Config) SaveCred() error {
	return c.save(credFileName, c.Cred.tokens)
}
