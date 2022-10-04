package app

import (
	"fmt"
	"testing"
	"time"

	"github.com/arrow2nd/nekome/v2/config"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/stretchr/testify/assert"
)

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
				Layout: config.Layout{
					TweetPoll:       "{graph}\n{detail}",
					TweetPollGraph:  "{label}\n{graph} {per} {votes}",
					TweetPollDetail: "{status} | {all_votes} | {end_date}",
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
		want := fmt.Sprintf(
			`test_1
[style_poll_g]==[-:-:-] 0.2%% [style_poll_d](2)[-:-:-]
test_2
[style_poll_g]=====[-:-:-] 0.5%% [style_poll_d](5)[-:-:-]
test_3
[style_poll_g]===[-:-:-] 0.3%% [style_poll_d](3)[-:-:-]
[style_poll_d]closed | 10 | %s[-:-:-]`,
			d,
		)

		assert.Equal(t, want, s)
	})
}
