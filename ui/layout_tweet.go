package ui

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

// createAnnotation : 「RT by」等のアノテーション文字列を作成
func createAnnotation(s string, author *twitter.UserObj) string {
	return fmt.Sprintf("[blue:-:-]%s %s [:-:i]@%s[-:-:-]", s, author.Name, author.UserName)
}

// createUserInfoLayout : レイアウト済みのユーザ情報文字列を作成
func createUserInfoLayout(u *twitter.UserObj, i, w int) string {
	name := truncate(u.Name, w/2)
	userName := truncate("@"+u.UserName, w/2)

	// カーソル選択用のタグを埋め込む
	if i >= 0 {
		name = fmt.Sprintf(`["%s"]%s[""]`, createTweetTag(i), name)
	}

	header := fmt.Sprintf(`[lightgray:-:b]%s [gray::i]%s[-:-:-]`, name, userName)

	// 認証済み
	if u.Verified {
		header += "[blue] \uf4a1[-:-:-]"
	}

	// 鍵付き
	if u.Protected {
		header += "[gray] \uf83d[-:-:-]"
	}

	return header + "\n"
}

// createPollLayout : レイアウト済みの投票文字列を作成
func createPollLayout(p []*twitter.PollObj) string {
	if len(p) == 0 {
		return ""
	}

	// グラフの表示幅を計算
	windowWidth := float64(getWindowWidth())
	graphMaxWidth := float64(30)

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
		per := float64(o.Votes) / float64(allVotes)
		graph := strings.Repeat("▇", int(math.Floor(per*graphMaxWidth)))

		text += fmt.Sprintf("\uf14a %s\n", o.Label)
		text += fmt.Sprintf("[blue]%s[-:-:-] %.1f%% (%d)\n", graph, per*100, o.Votes)
	}

	// 投票の詳細情報
	endDate := convertDateString(p[0].EndDateTime)
	text += fmt.Sprintf("[gray]%s | %d votes | ends on %s[-:-:-]\n\n", p[0].VotingStatus, allVotes, endDate)

	return text
}

// createTweetDetailLayout : レイアウト済みのツイート詳細文字列を作成
func createTweetDetailLayout(tw *twitter.TweetObj) string {
	metrics := ""

	likes := tw.PublicMetrics.Likes
	if likes != 0 {
		metrics += createMetricsString("Like", "pink", likes, false)
	}

	rts := tw.PublicMetrics.Retweets
	if rts != 0 {
		metrics += createMetricsString("RT", "green", rts, false)
	}

	if metrics != "" {
		metrics = "\n" + metrics
	}

	createAt := convertDateString(tw.CreatedAt)

	return fmt.Sprintf("[gray]%s | via %s[-:-:-]%s", createAt, tw.Source, metrics)
}

// createTweetTextLayout : レイアウト済みのツイート文字列を作成
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
		text = rep.ReplaceAllString(text, "$1[blue]@$2[-:-:-]")
	}

	return text
}

// highlightHashtags : ハッシュタグをハイライトしたツイート文字列を作成
func highlightHashtags(text string, entities *twitter.EntitiesObj) string {
	result := ""
	runes := []rune(text)
	end := 0

	for _, hashtag := range entities.HashTags {
		hashtagText := fmt.Sprintf("#%s", hashtag.Tag)

		// NOTE: URLや絵文字を多く含むツイートなどで、ハッシュタグの開始位置が後方にズレていることがあるので
		//       +1 して意図的にズラした後、ハッシュタグ全文が見つかるまで開始位置を前方に移動することで正しい位置を見つける

		start := hashtag.Start + 1
		textLength := utf8.RuneCountInString(hashtag.Tag) + 1

		for ; start > end; start-- {
			e := start + textLength
			if l := len(runes); e > l {
				e = l
			}

			if string(runes[start:e]) == hashtagText {
				break
			}
		}

		// 前方の文とハイライトされたハッシュタグを結合
		result += fmt.Sprintf("%s[blue]%s[-:-:-]", string(runes[end:start]), hashtagText)

		// ハッシュタグの終了位置
		end = start + utf8.RuneCountInString(hashtagText)
	}

	// 残りの文を結合
	if len(runes) > end {
		result += string(runes[end:])
	}

	return result
}

// createTweetLayout : レイアウト済みのツイート文字列を作成
func createTweetLayout(c *twitter.TweetDictionary, i, w int) string {
	return createUserInfoLayout(c.Author, i, w) +
		createTweetTextLayout(&c.Tweet) +
		createPollLayout(c.AttachmentPolls) +
		createTweetDetailLayout(&c.Tweet)
}
