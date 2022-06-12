package ui

import (
	"fmt"
	"strings"
	"sync"

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
		// if len(content.ReferencedTweets) != 0 {
		// 	fmt.Fprintf(t.textView, "type = %s\n", content.ReferencedTweets[0].Reference.Type)
		// }

		text := t.createHeader(content.Author, i)
		text += content.Tweet.Text + "\n"
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
		header += "[cyan] [-:-:-]"
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

func createMetricsString(unit, color string, count int, reverse bool) string {
	if count <= 0 {
		return ""
	} else if count > 1 {
		unit += "s"
	}

	if reverse {
		return fmt.Sprintf("[%s:-:r] %d%s [-:-:-]", color, count, unit)
	}

	return fmt.Sprintf("[%s]%d%s[-:-:-] ", color, count, unit)
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
