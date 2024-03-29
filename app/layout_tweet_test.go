package app

import (
	"fmt"
	"testing"
	"time"

	"github.com/arrow2nd/nekome/v2/config"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreateTweetTag(t *testing.T) {
	t.Run("作成できるか", func(t *testing.T) {
		s := createTweetTag(12345)
		want := "tweet_12345"

		assert.Equal(t, want, s)
	})
}

func TestCreateAnotation(t *testing.T) {
	shared = Shared{
		conf: &config.Config{
			Pref: &config.Preferences{
				Layout: config.Layout{
					TweetAnotation: "{text} {author_name} {author_username}",
				},
			},
			Style: &config.Style{
				Tweet: config.TweetStyle{
					Annotation: "style_anno",
				},
			},
		},
	}

	a := &twitter.UserObj{
		Name:     "市川雛菜",
		UserName: "ickwhnn",
	}

	t.Run("作成できるか", func(t *testing.T) {
		s := createAnnotation("RT by", a)
		want := "[style_anno]RT by 市川雛菜 @ickwhnn[-:-:-]"

		assert.Equal(t, want, s)
	})
}

func TestCreateTweetLayout(t *testing.T) {
	shared = Shared{
		conf: &config.Config{
			Pref: &config.Preferences{
				Appearance: config.Appearancene{
					DateFormat: "2006/01/02",
					TimeFormat: "15:04:05",
				},
				Layout: config.Layout{
					Tweet:       "{annotation}\n{user_info}\n{text}\n{poll}\n{detail}",
					TweetDetail: "{created_at} | via {via}\n{metrics}",
					User:        "{user_info}\n{bio}\n{user_detail}",
					UserInfo:    "{name} {username} {badge}",
				},
				Icon: config.Icon{
					Verified: "v",
					Private:  "p",
				},
				Text: config.Text{
					Like:    "like",
					Retweet: "rt",
				},
			},
			Style: &config.Style{
				Tweet: config.TweetStyle{
					Like:    "style_like",
					Retweet: "style_rt",
					Detail:  "style_detail",
					HashTag: "style_hashtag",
					Mention: "style_mention",
				},
				User: config.UserStyle{
					Name:     "style_name",
					UserName: "style_user_name",
					Verified: "style_verified",
					Private:  "style_private",
				},
			},
		},
	}

	td := &twitter.TweetDictionary{
		Tweet: twitter.TweetObj{
			ID:        "1234",
			Text:      "text",
			CreatedAt: "2022-04-18T15:00:00.000Z",
			PublicMetrics: &twitter.TweetMetricsObj{
				Likes:    2,
				Retweets: 2,
				Quotes:   1,
			},
			Source: "nekome for term",
		},
		Author: &twitter.UserObj{
			ID:       "5678",
			Name:     "user",
			UserName: "user_name",
		},
	}

	t.Run("作成できるか", func(t *testing.T) {
		p, _ := time.Parse(time.RFC3339, td.Tweet.CreatedAt)
		d := p.Local().Format("2006/01/02 15:04:05")

		s := createTweetLayout("annotation", td, 0, 250)
		want := fmt.Sprintf(`annotation
[style_name]["tweet_0"]user[""][-:-:-] [style_user_name]@user_name[-:-:-]
text
[style_detail]%s | via nekome for term
[style_like]2likes[-:-:-] [style_rt]2rts[-:-:-][-:-:-]`,
			d,
		)

		assert.Equal(t, want, s)
	})
}

func TestCreateUserInfoLayout(t *testing.T) {
	shared = Shared{
		conf: &config.Config{
			Pref: &config.Preferences{
				Layout: config.Layout{
					UserInfo: "{name} {username} {badge}",
				},
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
		want := `[style_name]["tweet_0"]hoge[""][-:-:-] [style_user_name]@fuga[-:-:-]`

		assert.Equal(t, want, s)
	})

	t.Run("認証済みアカウント", func(t *testing.T) {
		u := &twitter.UserObj{
			Name:      "櫻木真乃",
			UserName:  "sakuragi_mano_official",
			Verified:  true,
			Protected: false,
		}

		s := createUserInfoLayout(u, 0, 50)
		want := `[style_name]["tweet_0"]櫻木真乃[""][-:-:-] [style_user_name]@sakuragi_mano_official[-:-:-] [style_verified]v[-:-:-]`

		assert.Equal(t, want, s)
	})

	t.Run("非公開アカウント", func(t *testing.T) {
		u := &twitter.UserObj{
			Name:      "ルカ",
			UserName:  "ikrglc_0131",
			Verified:  false,
			Protected: true,
		}

		s := createUserInfoLayout(u, 0, 50)
		want := `[style_name]["tweet_0"]ルカ[""][-:-:-] [style_user_name]@ikrglc_0131[-:-:-] [style_private]p[-:-:-]`

		assert.Equal(t, want, s)
	})

	t.Run("認証済み&非公開アカウント", func(t *testing.T) {
		u := &twitter.UserObj{
			Name:      "N.Y.",
			UserName:  "ykmnm",
			Verified:  true,
			Protected: true,
		}

		s := createUserInfoLayout(u, 0, 50)
		want := `[style_name]["tweet_0"]N.Y.[""][-:-:-] [style_user_name]@ykmnm[-:-:-] [style_verified]v[-:-:-] [style_private]p[-:-:-]`

		assert.Equal(t, want, s)
	})
}

func TestCreateTextLayout(t *testing.T) {
	shared = Shared{
		conf: &config.Config{
			Style: &config.Style{
				Tweet: config.TweetStyle{
					HashTag: "style_hashtag",
					Mention: "style_mention",
				},
			},
		},
	}

	t.Run("全角記号が半角記号に置換されるか", func(t *testing.T) {
		o := &twitter.TweetObj{
			Text: "test ＠ ＃",
		}

		s := createTextLayout(o)
		want := "test @ #"

		assert.Equal(t, want, s)
	})

	t.Run("ハッシュタグがハイライトされるか", func(t *testing.T) {
		o := &twitter.TweetObj{
			Text: "test #hashtag",
			Entities: &twitter.EntitiesObj{
				HashTags: []twitter.EntityTagObj{
					{
						EntityObj: twitter.EntityObj{
							Start: 5,
							End:   12,
						},
						Tag: "hashtag",
					},
				},
			},
		}

		s := createTextLayout(o)
		want := "test [style_hashtag]#hashtag[-:-:-]"

		assert.Equal(t, want, s)
	})

	t.Run("メンションがハイライトされるか", func(t *testing.T) {
		o := &twitter.TweetObj{
			Text: "test @mention",
			Entities: &twitter.EntitiesObj{
				Mentions: []twitter.EntityMentionObj{
					{
						EntityObj: twitter.EntityObj{
							Start: 5,
							End:   12,
						},
					},
				},
			},
		}

		s := createTextLayout(o)
		want := "test [style_mention]@mention[-:-:-]"

		assert.Equal(t, want, s)
	})
}

func TestHighlightHashtags(t *testing.T) {
	shared = Shared{
		conf: &config.Config{
			Style: &config.Style{
				Tweet: config.TweetStyle{
					HashTag: "style_hashtag",
				},
			},
		},
	}

	text := "test tweet #hashtag #ハッシュタグ"
	e := twitter.EntitiesObj{
		HashTags: []twitter.EntityTagObj{
			{
				EntityObj: twitter.EntityObj{
					Start: 11,
					End:   19,
				},
				Tag: "hashtag",
			},
			{
				EntityObj: twitter.EntityObj{
					Start: 20,
					End:   26,
				},
				Tag: "ハッシュタグ",
			},
		},
	}

	t.Run("開始位置の値が正しい場合", func(t *testing.T) {
		s := highlightHashtags(text, &e)
		want := "test tweet [style_hashtag]#hashtag[-:-:-] [style_hashtag]#ハッシュタグ[-:-:-]"

		assert.Equal(t, want, s)
	})

	t.Run("開始位置がズレている場合", func(t *testing.T) {
		e2 := e

		for i := 0; i < 2; i++ {
			e2.HashTags[i].Start += 2
			e2.HashTags[i].End += 2
		}

		s := highlightHashtags(text, &e2)
		want := "test tweet [style_hashtag]#hashtag[-:-:-] [style_hashtag]#ハッシュタグ[-:-:-]"

		assert.Equal(t, want, s)
	})
}

func TestCreateTweetDetailLayout(t *testing.T) {
	shared = Shared{
		conf: &config.Config{
			Pref: &config.Preferences{
				Appearance: config.Appearancene{
					DateFormat: "2006/01/02",
					TimeFormat: "15:04:05",
				},
				Layout: config.Layout{
					TweetDetail: "{created_at} | via {via}\n{metrics}",
				},
				Text: config.Text{
					Like:    "like",
					Retweet: "rt",
				},
			},
			Style: &config.Style{
				Tweet: config.TweetStyle{
					Like:    "style_like",
					Retweet: "style_rt",
					Detail:  "style_detail",
				},
			},
		},
	}

	o := &twitter.TweetObj{
		CreatedAt: "2022-04-18T15:00:00.000Z",
		Source:    "nekome for term",
		PublicMetrics: &twitter.TweetMetricsObj{
			Likes:    10,
			Retweets: 5,
		},
	}

	p, _ := time.Parse(time.RFC3339, o.CreatedAt)
	d := p.Local().Format("2006/01/02 15:04:05")

	t.Run("作成できるか", func(t *testing.T) {
		s := createTweetDetailLayout(o)

		want := fmt.Sprintf(
			`[style_detail]%s | via nekome for term
[style_like]10likes[-:-:-] [style_rt]5rts[-:-:-][-:-:-]`,
			d,
		)

		assert.Equal(t, want, s)
	})

	t.Run("viaが空文字の場合にunknownが入るか", func(t *testing.T) {
		o := *o
		o.Source = ""
		s := createTweetDetailLayout(&o)

		want := fmt.Sprintf(
			`[style_detail]%s | via unknown
[style_like]10likes[-:-:-] [style_rt]5rts[-:-:-][-:-:-]`,
			d,
		)

		assert.Equal(t, want, s)
	})
}

func TestCreateTweetMetricsLayout(t *testing.T) {
	shared = Shared{
		conf: &config.Config{
			Pref: &config.Preferences{
				Text: config.Text{
					Like:    "like",
					Retweet: "rt",
				},
			},
			Style: &config.Style{
				Tweet: config.TweetStyle{
					Like:    "style_like",
					Retweet: "style_rt",
					Detail:  "style_detail",
				},
			},
		},
	}

	o := &twitter.TweetObj{
		PublicMetrics: &twitter.TweetMetricsObj{
			Likes:    10,
			Retweets: 5,
		},
	}

	t.Run("作成できるか", func(t *testing.T) {
		s := createTweetMetricsLayout(o)
		want := "[style_like]10likes[-:-:-] [style_rt]5rts[-:-:-]"

		assert.Equal(t, want, s)
	})
}
