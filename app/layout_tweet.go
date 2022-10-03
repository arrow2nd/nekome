package app

import (
	"fmt"
	"html"
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
func createAnnotation(text string, author *twitter.UserObj) string {
	layout := shared.conf.Pref.Layout.TweetAnotation

	layout = replaceLayoutTag(layout, "{text}", text)
	layout = replaceLayoutTag(layout, "{author_name}", author.Name)
	layout = replaceLayoutTag(layout, "{author_username}", "@"+author.UserName)

	return createStyledText(
		shared.conf.Style.Tweet.Annotation,
		layout,
	)
}

// createTweetLayout : ツイートのレイアウトを作成
func createTweetLayout(a string, d *twitter.TweetDictionary, i, w int) string {
	layout := shared.conf.Pref.Layout.Tweet

	layout = replaceLayoutTag(layout, "{annotation}", a)
	layout = replaceLayoutTag(layout, "{user_info}", createUserInfoLayout(d.Author, i, w))
	layout = replaceLayoutTag(layout, "{text}", createTextLayout(&d.Tweet))
	layout = replaceLayoutTag(layout, "{poll}", createPollLayout(d.AttachmentPolls, w))
	layout = replaceLayoutTag(layout, "{detail}", createTweetDetailLayout(&d.Tweet))
	layout = replaceLayoutTag(layout, "{metrics}", createTweetMetricsLayout(&d.Tweet))

	return trimEndNewline(layout)
}

// createUserInfoLayout : ユーザ情報のレイアウトを作成
func createUserInfoLayout(u *twitter.UserObj, i, w int) string {
	layout := shared.conf.Pref.Layout.UserInfo
	style := shared.conf.Style
	icon := shared.conf.Pref.Icon

	// 名前
	name := truncate(u.Name, w/2)

	// カーソル選択用のタグを埋め込む
	if i >= 0 {
		name = fmt.Sprintf(`["%s"]%s[""]`, createTweetTag(i), name)
	}

	layout = replaceLayoutTag(
		layout,
		"{name}",
		createStyledText(style.User.Name, name),
	)

	// ユーザネーム
	userName := createStyledText(
		style.User.UserName,
		truncate("@"+u.UserName, w/2),
	)

	layout = replaceLayoutTag(layout, "{username}", userName)

	// バッジ
	badges := []string{}

	if u.Verified {
		badges = append(
			badges,
			createStyledText(style.User.Verified, icon.Verified),
		)
	}

	if u.Protected {
		badges = append(
			badges,
			createStyledText(style.User.Private, icon.Private),
		)
	}

	layout = replaceLayoutTag(layout, "{badge}", strings.Join(badges, " "))

	return strings.TrimSpace(layout)
}

// createTextLayout : ツイート文のレイアウトを作成
func createTextLayout(t *twitter.TweetObj) string {
	text := html.UnescapeString(t.Text)

	// 全角記号を置換
	text = strings.ReplaceAll(text, "＃", "#")
	text = strings.ReplaceAll(text, "＠", "@")

	if t.Entities == nil {
		return text
	}

	// ハッシュタグをハイライト
	if len(t.Entities.HashTags) != 0 {
		text = highlightHashtags(text, t.Entities)
	}

	// メンションをハイライト
	if len(t.Entities.Mentions) != 0 {
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

// createTweetDetailLayout : ツイート詳細のレイアウトを作成
func createTweetDetailLayout(t *twitter.TweetObj) string {
	layout := shared.conf.Pref.Layout.TweetDetail

	// 投稿日時
	date := convertDateString(t.CreatedAt)
	layout = replaceLayoutTag(layout, "{created_at}", date)

	// 投稿元クライアント
	layout = replaceLayoutTag(layout, "{via}", t.Source)

	// メトリクス
	metrics := createTweetMetricsLayout(t)
	layout = replaceLayoutTag(layout, "{metrics}", metrics)

	return createStyledText(
		shared.conf.Style.Tweet.Detail,
		trimEndNewline(layout),
	)
}

// createTweetMetricsLayout : ツイートメトリクスのレイアウトを作成
func createTweetMetricsLayout(t *twitter.TweetObj) string {
	pref := shared.conf.Pref.Text
	style := shared.conf.Style.Tweet

	metrics := []string{}

	// いいね数
	if likes := t.PublicMetrics.Likes; likes != 0 {
		metrics = append(
			metrics,
			createMetricsString(pref.Like, style.Like, likes),
		)
	}

	// リツイート数
	if rts := t.PublicMetrics.Retweets; rts != 0 {
		metrics = append(
			metrics,
			createMetricsString(pref.Retweet, style.Retweet, rts),
		)
	}

	return strings.Join(metrics, " ")
}
