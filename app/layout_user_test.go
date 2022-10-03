package app

import (
	"testing"

	"github.com/arrow2nd/nekome/config"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/assert"
)

func TestCreateProfileLayout(t *testing.T) {
	shared = Shared{
		conf: &config.Config{
			Pref: &config.Preferences{
				Appearance: config.Appearancene{
					UserBIOMaxRow:       3,
					UserProfilePaddingX: 4,
				},
				Layout: config.Layout{
					User:     "{user_info}\n{bio}\n{user_detail}",
					UserInfo: "{name} {username} {badge}",
				},
				Icon: config.Icon{
					Geo:  "g",
					Link: "l",
				},
			},
			Style: &config.Style{
				User: config.UserStyle{
					Name:     "style_name",
					UserName: "style_user_name",
					Detail:   "style_detail",
				},
			},
		},
	}

	t.Run("詳細情報あり", func(t *testing.T) {
		u := &twitter.UserObj{
			Name:        "name",
			UserName:    "user_name",
			Description: "bio",
			Location:    "location",
			URL:         "url",
		}

		s, r := createProfileLayout(u, 100)
		want := `[style_name]name[-:-:-] [style_user_name]@user_name[-:-:-]
bio
[style_detail]g locationl url[-:-:-]`

		assert.Equal(t, want, s)
		assert.Equal(t, 3, r)
	})

	t.Run("詳細情報なし", func(t *testing.T) {
		u := &twitter.UserObj{
			Name:        "name",
			UserName:    "user_name",
			Description: "bio",
		}

		s, r := createProfileLayout(u, 100)
		want := `[style_name]name[-:-:-] [style_user_name]@user_name[-:-:-]
bio`

		assert.Equal(t, want, s)
		assert.Equal(t, 2, r)
	})
}

func TestCreateUserBioLayout(t *testing.T) {
	shared = Shared{
		conf: &config.Config{
			Pref: &config.Preferences{
				Appearance: config.Appearancene{
					UserBIOMaxRow: 3,
				},
			},
		},
	}

	t.Run("改行が空白に置換されるか", func(t *testing.T) {
		bio := `テスト
テスト
テスト`
		want := "テスト テスト テスト"

		s := createUserBioLayout(bio, 100)
		assert.Equal(t, want, s)
	})

	t.Run("最大表示幅を超えた場合に省略されるか", func(t *testing.T) {
		bio := "1234567890"
		want := "12345678…"

		runewidth.DefaultCondition.EastAsianWidth = false

		s := createUserBioLayout(bio, 3)
		assert.Equal(t, want, s)
	})
}

func TestCreateUserDetailLayout(t *testing.T) {
	shared = Shared{
		conf: &config.Config{
			Pref: &config.Preferences{
				Appearance: config.Appearancene{
					UserDetailSeparator: " | ",
				},
				Icon: config.Icon{
					Geo:  "g",
					Link: "l",
				},
			},
			Style: &config.Style{
				User: config.UserStyle{
					Detail: "style_detail",
				},
			},
		},
	}

	u := &twitter.UserObj{
		Location: "location",
		URL:      "url",
	}

	t.Run("作成できるか", func(t *testing.T) {
		s := createUserDetailLayout(u)
		want := "[style_detail]g location | l url[-:-:-]"

		assert.Equal(t, want, s)
	})
}
