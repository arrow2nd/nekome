package app

import (
	"fmt"
	"testing"
	"time"

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
		assert.Equal(t, s, "test @ #\n")
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
		assert.Equal(t, s, "test [style_hashtag]#hashtag[-:-:-]\n")
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
		assert.Equal(t, s, "test [style_mention]@mention[-:-:-]\n")
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
		assert.Equal(t, s, "test tweet [style_hashtag]#hashtag[-:-:-] [style_hashtag]#ハッシュタグ[-:-:-]")
	})

	t.Run("開始位置がズレている場合", func(t *testing.T) {
		e2 := e

		for i := 0; i < 2; i++ {
			e2.HashTags[i].Start += 2
			e2.HashTags[i].End += 2
		}

		s := highlightHashtags(text, &e2)
		assert.Equal(t, s, "test tweet [style_hashtag]#hashtag[-:-:-] [style_hashtag]#ハッシュタグ[-:-:-]")
	})
}

func TestCreatePollLayout(t *testing.T) {
	shared = Shared{
		conf: &config.Config{
			Pref: &config.Preferences{
				Appearance: config.Appearancene{
					DateFormat:    "2006/01/02",
					TimeFormat:    "15:04:05",
					GraphMaxWidth: 10,
					GraphChar:     "=",
				},
			},
			Style: &config.Style{
				Tweet: config.TweetStyle{
					PollGraph:  "style_poll_g",
					PollDetail: "style_poll_d",
				},
			},
		},
	}

	p := []*twitter.PollObj{
		{
			ID: "1234567890",
			Options: []*twitter.PollOptionObj{
				{
					Position: 1,
					Label:    "test_1",
					Votes:    2,
				},
				{
					Position: 2,
					Label:    "test_2",
					Votes:    5,
				},
				{
					Position: 3,
					Label:    "test_3",
					Votes:    3,
				},
			},
			DurationMinutes: 60,
			EndDateTime:     "2022-04-18T15:00:00.000Z",
			VotingStatus:    "closed",
		},
	}

	t.Run("生成できるか", func(t *testing.T) {
		s := createPollLayout(p, 120)

		p, _ := time.Parse(time.RFC3339, p[0].EndDateTime)
		d := p.Local().Format("2006/01/02 15:04:05")
		want := fmt.Sprintf(`
test_1[style_poll_g]==[-:-:-] 20.0%% (2)
test_2[style_poll_g]=====[-:-:-] 50.0%% (5)
test_3[style_poll_g]===[-:-:-] 30.0%% (3)
[style_poll_d]closed | 10 votes | ends on %s[-:-:-]

`, d)

		assert.Equal(t, s, want)
	})
}

func TestCreateDetailLayout(t *testing.T) {
	shared = Shared{
		conf: &config.Config{
			Pref: &config.Preferences{
				Appearance: config.Appearancene{
					DateFormat: "2006/01/02",
					TimeFormat: "15:04:05",
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

	t.Run("作成できるか", func(t *testing.T) {
		s := createTweetDetailLayout(o)

		p, _ := time.Parse(time.RFC3339, o.CreatedAt)
		d := p.Local().Format("2006/01/02 15:04:05")
		want := fmt.Sprintf(`[style_detail]%s | via nekome for term[-:-:-]
[style_like]10likes[-:-:-] [style_rt]5rts[-:-:-]`, d)

		assert.Equal(t, s, want)
	})
}
