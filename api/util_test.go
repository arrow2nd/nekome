package api

import (
	"testing"

	"github.com/dghubble/oauth1"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetConsumerToken(t *testing.T) {
	consumerToken = "token"
	consumerSecret = "secret"

	t.Run("内臓のコンシューマキーが取得できるか", func(t *testing.T) {
		tk := &oauth1.Token{}
		c, err := getConsumerToken(tk)

		want := oauth1.Token{
			Token:       "token",
			TokenSecret: "secret",
		}

		assert.NoError(t, err)
		assert.Equal(t, want, *c)
	})

	t.Run("引数に渡したコンシューマキーが返るか", func(t *testing.T) {
		want := oauth1.Token{
			Token:       "token",
			TokenSecret: "secret",
		}

		c, err := getConsumerToken(&want)

		assert.NoError(t, err)
		assert.Equal(t, want, *c)
	})

	consumerToken = ""
	consumerSecret = ""

	t.Run("取得可能なコンシューマキーがない場合にエラーが返るか", func(t *testing.T) {
		tk := &oauth1.Token{}
		_, err := getConsumerToken(tk)

		assert.EqualError(
			t,
			err,
			"no consumer key, please run 'nekome edit' and set it to .cred.toml",
		)
	})
}

func TestCreateTweetSlice(t *testing.T) {
	t.Run("作成できるか", func(t *testing.T) {
		r := &twitter.TweetRaw{
			Tweets: []*twitter.TweetObj{
				{
					ID:                 "1234567890",
					Text:               "tweet text",
					Attachments:        &twitter.TweetAttachmentsObj{},
					AuthorID:           "",
					ContextAnnotations: []*twitter.TweetContextAnnotationObj{},
					ConversationID:     "",
					CreatedAt:          "",
					Entities:           &twitter.EntitiesObj{},
					Geo:                &twitter.TweetGeoObj{},
					InReplyToUserID:    "",
					Language:           "",
					NonPublicMetrics:   &twitter.TweetMetricsObj{},
					OrganicMetrics:     &twitter.TweetMetricsObj{},
					PossiblySensitive:  false,
					PromotedMetrics:    &twitter.TweetMetricsObj{},
					PublicMetrics:      &twitter.TweetMetricsObj{},
					ReferencedTweets:   []*twitter.TweetReferencedTweetObj{},
					Source:             "",
					WithHeld:           &twitter.WithHeldObj{},
				},
			},
		}

		ok, d := createTweetSlice(r)

		assert.True(t, ok)
		assert.Len(t, d, 1)

		wantId := "1234567890"
		assert.Equal(t, wantId, d[0].Tweet.ID)
	})

	t.Run("データがない場合にfalseが返るか", func(t *testing.T) {
		r := &twitter.TweetRaw{
			Tweets: []*twitter.TweetObj{},
		}

		ok, d := createTweetSlice(r)

		assert.False(t, ok)
		assert.Nil(t, d)
	})
}

func TestCreateUserSlice(t *testing.T) {
	t.Run("作成できるか", func(t *testing.T) {
		r := &twitter.UserRaw{
			Users: []*twitter.UserObj{
				{
					ID:            "1234567890",
					Name:          "name",
					UserName:      "user_name",
					Description:   "bio",
					Entities:      &twitter.EntitiesObj{},
					Protected:     false,
					PublicMetrics: &twitter.UserMetricsObj{},
					Verified:      false,
					WithHeld:      &twitter.WithHeldObj{},
				},
			},
		}

		ok, d := createUserSlice(r)

		assert.True(t, ok)
		assert.Len(t, d, 1)

		wantId := "1234567890"
		assert.Equal(t, wantId, d[0].User.ID)
	})

	t.Run("データがない場合にfalseが返るか", func(t *testing.T) {
		r := &twitter.UserRaw{
			Users: []*twitter.UserObj{},
		}

		ok, d := createUserSlice(r)

		assert.False(t, ok)
		assert.Nil(t, d)
	})
}
