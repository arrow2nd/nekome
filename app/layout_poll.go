package app

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/g8rswimmer/go-twitter/v2"
)

// createPollLayout : 投票のレイアウトを作成
func createPollLayout(p []*twitter.PollObj, w int) string {
	if len(p) == 0 {
		return ""
	}

	// 総投票数
	allVotes := 0
	for _, o := range p[0].Options {
		allVotes += o.Votes
	}

	layout := shared.conf.Pref.Layout.TweetPoll

	// グラフ
	graph := createPollGraphLayout(p[0].Options, allVotes, w)
	layout = replaceLayoutTag(layout, "{graph}", graph)

	// 詳細
	detail := createPollDetailLayout(p[0], allVotes)
	layout = replaceLayoutTag(layout, "{detail}", detail)

	return layout
}

// createPollGraphLayout : 投票率グラフのレイアウトを作成
func createPollGraphLayout(p []*twitter.PollOptionObj, allVotes, w int) string {
	appearance := shared.conf.Pref.Appearance
	style := shared.conf.Style.Tweet

	// グラフの表示幅を計算
	windowWidth := float64(w)
	graphMaxWidth := float64(appearance.GraphMaxWidth)

	if graphMaxWidth > windowWidth {
		graphMaxWidth = windowWidth
	}

	graphs := []string{}

	for _, o := range p {
		per := float64(0)
		if allVotes > 0 {
			per = float64(o.Votes) / float64(allVotes)
		}

		pl := shared.conf.Pref.Layout.TweetPollGraph

		// 項目ラベル
		pl = replaceLayoutTag(pl, "{label}", o.Label)

		// グラフ
		graphWidth := int(math.Floor(per * graphMaxWidth))
		graph := createStyledText(
			style.PollGraph,
			strings.Repeat(appearance.GraphChar, graphWidth),
		)
		pl = replaceLayoutTag(pl, "{graph}", graph)

		// 得票率
		pl = replaceLayoutTag(pl, "{per}", fmt.Sprintf("%.1f%%", per*100))

		// 投票数
		votes := createStyledText(style.PollDetail, fmt.Sprintf("(%d)", o.Votes))
		pl = replaceLayoutTag(pl, "{votes}", votes)

		graphs = append(graphs, pl)
	}

	return strings.Join(graphs, "\n")
}

// createPollDetailLayout : 投票詳細のレイアウトを作成
func createPollDetailLayout(p *twitter.PollObj, allVotes int) string {
	layout := shared.conf.Pref.Layout.TweetPollDetail

	// 投票ステータス・総投票数
	layout = replaceLayoutTag(layout, "{status}", p.VotingStatus)
	layout = replaceLayoutTag(layout, "{all_votes}", strconv.Itoa(allVotes))

	// 投票終了日時
	endDate := convertDateString(p.EndDateTime)
	layout = replaceLayoutTag(layout, "{end_date}", endDate)

	return createStyledText(
		shared.conf.Style.Tweet.PollDetail,
		layout,
	)
}
