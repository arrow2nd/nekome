package config_test

import (
	"testing"

	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
	"github.com/dghubble/oauth1"
	"github.com/stretchr/testify/assert"
)

func newTestCred() *config.Cred {
	return &config.Cred{
		Consumer: oauth1.Token{
			Token:       "consumer_token",
			TokenSecret: "consumer_secret",
		},
		User: config.User{
			Accounts: []api.User{
				{
					UserName: "user_name_a",
					ID:       "user_id_a",
					Token: &oauth1.Token{
						Token:       "user_token_a",
						TokenSecret: "user_secret_a",
					},
				},
				{
					UserName: "user_name_b",
					ID:       "user_id_b",
					Token: &oauth1.Token{
						Token:       "user_token_b",
						TokenSecret: "user_secret_b",
					},
				},
			},
		},
	}
}

func TestCredGet(t *testing.T) {
	c := newTestCred()

	t.Run("取得できるか", func(t *testing.T) {
		u, err := c.Get("user_name_a")
		assert.NoError(t, err)
		assert.Equal(t, *u, c.User.Accounts[0])
	})

	t.Run("見つからなかった際にエラーが返るか", func(t *testing.T) {
		_, err := c.Get("hoge")
		assert.ErrorContains(t, err, "user not found: hoge")
	})
}

func TestGetAllNames(t *testing.T) {
	c := newTestCred()

	t.Run("取得できるか", func(t *testing.T) {
		l := c.GetAllNames()
		want := []string{
			"user_name_a",
			"user_name_b",
		}

		assert.Equal(t, l, want)
	})
}

func TestWrite(t *testing.T) {
	t.Run("追加できるか", func(t *testing.T) {
		c := newTestCred()

		want := &api.User{
			UserName: "hiori",
			ID:       "hio_hio",
			Token: &oauth1.Token{
				Token:       "hio_token",
				TokenSecret: "hio_secret",
			},
		}

		c.Write(want)

		u, _ := c.Get("hiori")
		assert.Equal(t, u, want)
	})

	t.Run("同じIDを持つユーザを上書きできるか", func(t *testing.T) {
		c := newTestCred()

		want := &api.User{
			UserName: "meguru",
			ID:       "user_id_b",
			Token: &oauth1.Token{
				Token:       "meguru_token",
				TokenSecret: "meguru_secret",
			},
		}

		c.Write(want)

		u, _ := c.Get("meguru")
		assert.Equal(t, u, want)
	})
}

func TestDelete(t *testing.T) {
	t.Run("削除できるか", func(t *testing.T) {
		c := newTestCred()

		err := c.Delete("user_name_a")
		assert.NoError(t, err)

		_, err = c.Get("user_name_a")
		assert.ErrorContains(t, err, "user not found: user_name_a")
	})

	t.Run("見つからない場合にエラーが返るか", func(t *testing.T) {
		c := newTestCred()

		err := c.Delete("hoge")
		assert.ErrorContains(t, err, "user not found: hoge")
	})
}
