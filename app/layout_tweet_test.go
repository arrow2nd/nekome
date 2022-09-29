package app

import (
	"testing"

	"github.com/arrow2nd/nekome/config"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreateTweetTag(t *testing.T) {
	t.Run("作成できるか", func(t *testing.T) {
		tag := createTweetTag(12345)
		assert.Equal(t, tag, "tweet_12345")
	})
}

func TestCreateAnotation(t *testing.T) {
	shared = Shared{
		conf: &config.Config{
			Style: &config.Style{
				Tweet: config.TweetStyle{
					Annotation: "-:-:-",
				},
			},
		},
	}

	a := &twitter.UserObj{
		Name:     "市川雛菜",
		UserName: "ickwhnn",
	}

	t.Run("作成できるか", func(t *testing.T) {
		ano := createAnnotation("RT by", a)
		assert.Equal(t, ano, "[-:-:-]RT by 市川雛菜 [::i]@ickwhnn[-:-:-]")
	})
}

// func TestCreateTweetLayout(t *testing.T) {
// }

func TestCreateUserInfoLayout(t *testing.T) {
	shared = Shared{
		conf: &config.Config{
			Pref: &config.Preferences{
				Icon: config.Icon{
					Verified: "v",
					Private:  "p",
				},
			},
			Style: &config.Style{
				User: config.UserStyle{
					Name:     "style_name",
					UserName: "style_user_name",
					Verified: "style_verified",
					Private:  "style_private",
				},
			},
		},
	}

	t.Run("通常のアカウント", func(t *testing.T) {
		u := &twitter.UserObj{
			Name:      "hoge",
			UserName:  "fuga",
			Verified:  false,
			Protected: false,
		}

		s := createUserInfoLayout(u, 0, 50)
		assert.Equal(t, s, `[style_name]["tweet_0"]hoge[""] [style_user_name]@fuga[-:-:-]`+"\n")
	})

	t.Run("認証済みアカウント", func(t *testing.T) {
		u := &twitter.UserObj{
			Name:      "櫻木真乃",
			UserName:  "sakuragi_mano_official",
			Verified:  true,
			Protected: false,
		}

		s := createUserInfoLayout(u, 0, 50)
		assert.Equal(t, s, `[style_name]["tweet_0"]櫻木真乃[""] [style_user_name]@sakuragi_mano_official[-:-:-][style_verified] v[-:-:-]`+"\n")
	})

	t.Run("非公開アカウント", func(t *testing.T) {
		u := &twitter.UserObj{
			Name:      "ルカ",
			UserName:  "ikrglc_0131",
			Verified:  false,
			Protected: true,
		}

		s := createUserInfoLayout(u, 0, 50)
		assert.Equal(t, s, `[style_name]["tweet_0"]ルカ[""] [style_user_name]@ikrglc_0131[-:-:-][style_private] p[-:-:-]`+"\n")
	})

	t.Run("認証済み&非公開アカウント", func(t *testing.T) {
		u := &twitter.UserObj{
			Name:      "N.Y.",
			UserName:  "ykmnm",
			Verified:  true,
			Protected: true,
		}

		s := createUserInfoLayout(u, 0, 50)
		assert.Equal(t, s, `[style_name]["tweet_0"]N.Y.[""] [style_user_name]@ykmnm[-:-:-][style_verified] v[-:-:-][style_private] p[-:-:-]`+"\n")
	})
}
