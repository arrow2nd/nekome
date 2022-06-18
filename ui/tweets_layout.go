package ui

import (
	"fmt"
	"html"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/g8rswimmer/go-twitter/v2"
)

func createSeparator(s string, width int) string {
	return fmt.Sprintf("[gray:-:-]%s[-:-:-]", strings.Repeat(s, width))
}

func createTweetId(id int) string {
	return fmt.Sprintf("tweet_%d", id)
}

func createAnnotation(s string, author *twitter.UserObj) string {
	return fmt.Sprintf("[blue:-:-]%s %s [:-:i]@%s[-:-:-]", s, author.Name, author.UserName)
}

func createTweetLayout(content *twitter.TweetDictionary, index int) string {
	text := createHeader(content.Author, index)
	text += createTweetText(&content.Tweet)
	text += createFooter(&content.Tweet)

	return text
}

func createHeader(u *twitter.UserObj, i int) string {
	name := u.Name

	if i >= 0 {
		name = fmt.Sprintf(`["%s"]%s[""]`, createTweetId(i), u.Name)
	}

	header := fmt.Sprintf(`[lightgray::b]%s [gray::i]@%s[-:-:-]`, name, u.UserName)

	if u.Verified {
		header += "[blue] [-:-:-]"
	}

	if u.Protected {
		header += "[gray] [-:-:-]"
	}

	return header + "\n"
}

func createFooter(tw *twitter.TweetObj) string {
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
	info := fmt.Sprintf("[gray]%s - via %s[-:-:-]", createAt, tw.Source)

	return info + metrics
}

func createTweetText(tweet *twitter.TweetObj) string {
	text := html.UnescapeString(tweet.Text) + "\n"

	if tweet.Entities == nil {
		return text
	}

	// ハッシュタグをハイライト
	if len(tweet.Entities.HashTags) != 0 {
		text = highlightHashtags(text, tweet.Entities)
	}

	// メンションをハイライト
	if len(tweet.Entities.Mentions) != 0 {
		rep := regexp.MustCompile(`(^|[^\w@#$%&])[@＠](\w+)`)
		text = rep.ReplaceAllString(text, "$1[blue]@$2[-:-:-]")
	}

	return text
}

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
