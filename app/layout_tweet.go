package app

import (
	"fmt"
	"html"
	"math"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/g8rswimmer/go-twitter/v2"
)

// createTweetTag : ツイート管理用のタグ文字列を作成
func createTweetTag(id int) string {
	return fmt.Sprintf("tweet_%d", id)
}

// createAnnotation : アノテーションを作成
func createAnnotation(s string, author *twitter.UserObj) string {
	return fmt.Sprintf(
		"[%s]%s %s [::i]@%s[-:-:-]",
		shared.conf.Style.Tweet.Annotation,
		s,
		author.Name,
		author.UserName,
	)
}

// createTweetLayout : ツイートのレイアウトを作成
func createTweetLayout(c *twitter.TweetDictionary, i, w int) string {
	return createUserInfoLayout(c.Author, i, w) +
		createTweetTextLayout(&c.Tweet) +
		createPollLayout(c.AttachmentPolls) +
		createTweetDetailLayout(&c.Tweet)
}

// createUserInfoLayout : ユーザ情報のレイアウトを作成
func createUserInfoLayout(u *twitter.UserObj, i, w int) string {
	name := truncate(u.Name, w/2)
	userName := truncate("@"+u.UserName, w/2)

	// カーソル選択用のタグを埋め込む
	if i >= 0 {
		name = fmt.Sprintf(`["%s"]%s[""]`, createTweetTag(i), name)
	}

	style := shared.conf.Style
	icon := shared.conf.Pref.Icon

	// ニックネーム・ユーザ名
	header := fmt.Sprintf(
		`[%s]%s [%s]%s[-:-:-]`,
		style.User.Name,
		name,
		style.User.UserName,
		userName,
	)

	// 認証済みバッジを追加
	if u.Verified {
		header += fmt.Sprintf("[%s] %s[-:-:-]", style.User.Verified, icon.Verified)
	}

	// 非公開バッジを追加
	if u.Protected {
		header += fmt.Sprintf("[%s] %s[-:-:-]", style.User.Private, icon.Private)
	}

	return header + "\n"
}

// createTweetTextLayout : ツイート文のレイアウトを作成
func createTweetTextLayout(tweet *twitter.TweetObj) string {
	text := html.UnescapeString(tweet.Text) + "\n"

	// 全角記号を置換
	text = strings.ReplaceAll(text, "＃", "#")
	text = strings.ReplaceAll(text, "＠", "@")

	if tweet.Entities == nil {
		return text
	}

	// ハッシュタグをハイライト
	if len(tweet.Entities.HashTags) != 0 {
		text = highlightHashtags(text, tweet.Entities)
	}

	// メンションをハイライト
	if len(tweet.Entities.Mentions) != 0 {
		rep := regexp.MustCompile(`(^|[^\w@#$%&])@(\w+)`)
		highlight := fmt.Sprintf("$1[%s]@$2[-:-:-]", shared.conf.Style.Tweet.Mention)
		text = rep.ReplaceAllString(text, highlight)
	}

	return text
}

// highlightHashtags : ツイート文内のハッシュタグをハイライト
func highlightHashtags(text string, entities *twitter.EntitiesObj) string {
	result := ""
	runes := []rune(text)
	end := 0

	for _, hashtag := range entities.HashTags {
		hashtagText := fmt.Sprintf("#%s", hashtag.Tag)

		// NOTE: URLや絵文字を多く含むツイートなどで、APIで取得できるハッシュタグの開始位置が後方にズレていることがあるので
		//       +1 して意図的にズラした後、ハッシュタグ全文が見つかるまで開始位置を前方に移動することで正しい位置を見つける

		start := hashtag.Start + 1
		textLength := utf8.RuneCountInString(hashtag.Tag) + 1

		for ; start > end; start-- {
			e := start + textLength

			if e > len(runes) {
				continue
			}

			if string(runes[start:e]) == hashtagText {
				break
			}
		}

		// 前方の文とハイライトされたハッシュタグを結合
		result += fmt.Sprintf(
			"%s[%s]%s[-:-:-]",
			string(runes[end:start]),
			shared.conf.Style.Tweet.HashTag,
			hashtagText,
		)

		// ハッシュタグの終了位置
		end = start + utf8.RuneCountInString(hashtagText)
	}

	// 残りの文を結合
	if len(runes) > end {
		result += string(runes[end:])
	}

	return result
}

// createPollLayout : 投票のレイアウトを作成
func createPollLayout(p []*twitter.PollObj) string {
	if len(p) == 0 {
		return ""
	}

	style := shared.conf.Style.Tweet
	pref := shared.conf.Pref.Appearance

	// グラフの表示幅を計算
	windowWidth := float64(getWindowWidth())
	graphMaxWidth := float64(pref.GraphMaxWidth)

	if graphMaxWidth > windowWidth {
		graphMaxWidth = windowWidth
	}

	// 総投票数を計算
	allVotes := 0
	for _, o := range p[0].Options {
		allVotes += o.Votes
	}

	// グラフを作成
	text := "\n"
	for _, o := range p[0].Options {
		per := float64(0)
		if allVotes > 0 {
			per = float64(o.Votes) / float64(allVotes)
		}

		graph := strings.Repeat(pref.GraphChar, int(math.Floor(per*graphMaxWidth)))
		text += fmt.Sprintf(
			"%s[%s]%s[-:-:-] %.1f%% (%d)\n",
			o.Label,
			style.PollGraph,
			graph,
			per*100,
			o.Votes,
		)
	}

	// 投票の詳細情報
	endDate := convertDateString(p[0].EndDateTime)
	text += fmt.Sprintf(
		"[%s]%s | %d votes | ends on %s[-:-:-]\n\n",
		style.PollDetail,
		p[0].VotingStatus,
		allVotes,
		endDate,
	)

	return text
}

// createTweetDetailLayout : ツイート詳細のレイアウトを作成
func createTweetDetailLayout(tw *twitter.TweetObj) string {
	pref := shared.conf.Pref.Text
	style := shared.conf.Style.Tweet

	metrics := ""

	// いいね数
	likes := tw.PublicMetrics.Likes
	if likes != 0 {
		metrics += createMetricsString(pref.Like, style.Like, likes, false)
	}

	// リツイート数
	rts := tw.PublicMetrics.Retweets
	if rts != 0 {
		metrics += createMetricsString(pref.Retweet, style.Retweet, rts, false)
	}

	if metrics != "" {
		metrics = "\n" + metrics
	}

	// 投稿日時・投稿元クライアント
	date := convertDateString(tw.CreatedAt)
	return fmt.Sprintf("[%s]%s | via %s[-:-:-]%s", style.Detail, date, tw.Source, metrics)
}
