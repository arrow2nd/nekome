package ui

import (
	"fmt"
	"html"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tweets struct {
	textView *tview.TextView
	contents []*twitter.TweetDictionary
	count    int
	mu       sync.Mutex
}

func newTweets() *tweets {
	t := &tweets{
		textView: tview.NewTextView(),
		contents: []*twitter.TweetDictionary{},
		count:    0,
	}

	t.textView.SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true).
		SetRegions(true).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			t.textView.ScrollToHighlight()
		})

	t.textView.SetBackgroundColor(tcell.ColorDefault)
	t.textView.SetInputCapture(t.handleKeyEvents)

	return t
}

func (t *tweets) createTweetId(id int) string {
	return fmt.Sprintf("tweet_%d", id)
}

func (t *tweets) getSinceID() string {
	if len(t.contents) == 0 {
		return ""
	}

	return t.contents[0].Tweet.ID
}

func (t *tweets) register(tweets []*twitter.TweetDictionary) {
	t.mu.Lock()
	defer t.mu.Unlock()

	length := len(tweets)
	t.contents = append(tweets, t.contents...)
	t.count = len(t.contents)

	shared.setStatus(fmt.Sprintf("%d tweets loaded", length))
}

func (t *tweets) draw() {
	width := getWindowWidth()

	t.textView.Clear()

	for i, content := range t.contents {
		// 参照元のツイートに置換
		if len(content.ReferencedTweets) != 0 {
			ref := content.ReferencedTweets[0]

			switch ref.Reference.Type {
			case "retweeted":
				fmt.Fprintf(t.textView, "[green:-:-]RT by %s [:-:i]@%s[-:-:-]\n", content.Author.Name, content.Author.UserName)
				content = ref.TweetDictionary
			case "quoted":
				fmt.Fprintf(t.textView, "type = %s\n", content.ReferencedTweets[0].Reference.Type)
			}
		}

		text := t.createHeader(content.Author, i)
		text += t.createTweetText(&content.Tweet)
		text += t.createFooter(&content.Tweet)

		fmt.Fprintf(t.textView, "%s\n", text)

		if i < t.count-1 {
			fmt.Fprintf(t.textView, "[gray]%s[-:-:-]\n", strings.Repeat("─", width-1))
		}
	}

	t.scrollToTweet(0)
}

func (t *tweets) createHeader(u *twitter.UserObj, i int) string {
	id := t.createTweetId(i)
	header := fmt.Sprintf(`[lightgray::b]["%s"]%s[""] [gray::i]@%s[-:-:-]`, id, u.Name, u.UserName)

	if u.Verified {
		header += "[blue] [-:-:-]"
	}

	if u.Protected {
		header += "[gray] [-:-:-]"
	}

	return header + "\n"
}

func (t *tweets) createFooter(tw *twitter.TweetObj) string {
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

func (t *tweets) createTweetText(tweet *twitter.TweetObj) string {
	text := html.UnescapeString(tweet.Text) + "\n"

	if tweet.Entities == nil {
		return text
	}

	// ハッシュタグをハイライト
	if len(tweet.Entities.HashTags) != 0 {
		text = t.highlightHashtags(text, tweet.Entities)
	}

	// メンションをハイライト
	if len(tweet.Entities.Mentions) != 0 {
		rep := regexp.MustCompile(`(^|[^\w@#$%&])[@＠](\w+)`)
		text = rep.ReplaceAllString(text, "$1[green]@$2[-:-:-]")
	}

	return text
}

func (t *tweets) highlightHashtags(text string, entities *twitter.EntitiesObj) string {
	result := ""
	runes := []rune(text)
	endPos := 0

	for _, hashtag := range entities.HashTags {
		// ハッシュタグの開始位置 ("#"を含まない)
		beginPos := hashtag.Start + 1
		textLength := utf8.RuneCountInString(hashtag.Tag) + 1

		// NOTE: 絵文字の表示幅の関係で、開始位置が実際の値より大きい場合があるので
		//       ハッシュタグが見つかるまで開始位置を前方にズラし、切り出した文字列がタグ名を含むかチェックする
		//       終了条件が i > 0 なので、beginPos は "#" を含むハッシュタグの開始位置になる
		for ; beginPos > endPos; beginPos-- {
			if i := strings.Index(string(runes[beginPos:beginPos+textLength]), hashtag.Tag); i > 0 {
				break
			}
		}

		// 前方の文とハイライトされたハッシュタグを結合
		hashtagText := fmt.Sprintf("#%s", hashtag.Tag)
		result += fmt.Sprintf("%s[blue]%s[-:-:-]", string(runes[endPos:beginPos]), hashtagText)

		// ハッシュタグの終了位置
		endPos = beginPos + utf8.RuneCountInString(hashtagText)
	}

	// 残りの文を結合
	if len(runes) > endPos {
		result += string(runes[endPos:])
	}

	return result
}

func (t *tweets) scrollToTweet(i int) {
	t.textView.Highlight(t.createTweetId(i))
}

func (t *tweets) cursorUp() {
	idx := getHighlightId(t.textView.GetHighlights())
	if idx == -1 {
		return
	}

	if idx--; idx < 0 {
		idx = t.count - 1
	}

	t.scrollToTweet(idx)
}

func (t *tweets) cursorDown() {
	idx := getHighlightId(t.textView.GetHighlights())
	if idx == -1 {
		return
	}

	idx = (idx + 1) % t.count

	t.scrollToTweet(idx)
}

func (t *tweets) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	keyRune := event.Rune()

	if key == tcell.KeyUp || keyRune == 'k' {
		t.cursorUp()
		return nil
	}

	if key == tcell.KeyDown || keyRune == 'j' {
		t.cursorDown()
		return nil
	}

	if keyRune == 'g' {
		t.scrollToTweet(0)
		return nil
	}

	if keyRune == 'G' {
		t.scrollToTweet(t.count - 1)
		return nil
	}

	return event
}
