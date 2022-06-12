package ui

import (
	"fmt"
	"strings"
	"sync"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

type tweets struct {
	textView *tview.TextView
	contents []*twitter.TweetDictionary
	index    int
	mu       sync.Mutex
}

func newTweets() *tweets {
	t := &tweets{
		textView: tview.NewTextView(),
		contents: []*twitter.TweetDictionary{},
		index:    0,
	}

	t.textView.SetDynamicColors(true).
		SetScrollable(true).
		SetRegions(true).
		SetWrap(false).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			t.textView.ScrollToHighlight()
		})

	t.textView.SetBackgroundColor(tcell.ColorDefault)
	t.textView.SetInputCapture(t.handleKeyEvents)

	return t
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

	t.index = 0
	length := len(tweets)
	t.contents = append(tweets, t.contents...)

	shared.setStatus(fmt.Sprintf("%d tweets loaded", length))
}

func (t *tweets) draw() {
	width := getWindowWidth()
	t.textView.Clear()

	for i, content := range t.contents {
		// if len(content.ReferencedTweets) != 0 {
		// 	fmt.Fprintf(t.textView, "type = %s\n", content.ReferencedTweets[0].Reference.Type)
		// }

		text := t.createHeader(content.Author)
		text += runewidth.Wrap(content.Tweet.Text, width) + "\n"
		text += t.createFooter(&content.Tweet)

		t.printTweet(i, text)
	}

	t.textView.Highlight(t.createTweetId(0))
}

func (t *tweets) createHeader(u *twitter.UserObj) string {
	header := fmt.Sprintf("[white::b]%s [gray::i]@%s[-:-:-]", u.Name, u.UserName)

	if u.Verified {
		header += "[cyan] [default]"
	}

	if u.Protected {
		header += "[gray] [default]"
	}

	return header + "\n"
}

func (t *tweets) createFooter(tw *twitter.TweetObj) string {
	metrics := ""

	likes := tw.PublicMetrics.Likes
	if likes != 0 {
		metrics += createMetricsString("Like", "pink", likes, false) + " "
	}

	rts := tw.PublicMetrics.Retweets
	if rts != 0 {
		metrics += createMetricsString("RT", "green", rts, false) + " "
	}

	createAt := convertDateString(tw.CreatedAt)
	via := fmt.Sprintf("via %s", tw.Source)

	return fmt.Sprintf("%s%s - %s", metrics, createAt, via)
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

	return fmt.Sprintf("[%s]%d%s[default]", color, count, unit)
}

func (t *tweets) printTweet(i int, text string) {
	cursor := fmt.Sprintf(`[blue]["tweet_%d"] [""][default] `, i)
	fmt.Fprintf(t.textView, "%s%s", cursor, strings.Replace(text, "\n", "\n"+cursor, -1))
	fmt.Fprint(t.textView, "\n\n")
}

func (t *tweets) cursorUp() {
	idx := getHighlightId(t.textView.GetHighlights())
	if idx == -1 {
		return
	}

	if idx--; idx < 0 {
		idx = len(t.contents) - 1
	}

	t.textView.Highlight(t.createTweetId(idx))
}

func (t *tweets) cursorDown() {
	idx := getHighlightId(t.textView.GetHighlights())
	if idx == -1 {
		return
	}

	idx = (idx + 1) % len(t.contents)

	t.textView.Highlight(t.createTweetId(idx))
}

func (t *tweets) createTweetId(id int) string {
	return fmt.Sprintf("tweet_%d", id)
}

func (t *tweets) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyUp:
		t.cursorUp()
		return nil
	case tcell.KeyDown:
		t.cursorDown()
		return nil
	case tcell.KeyRune:
		switch event.Rune() {
		case 'k':
			t.cursorUp()
			return nil
		case 'j':
			t.cursorDown()
			return nil
		case 'g':
			t.textView.Highlight(t.createTweetId(0))
			return nil
		case 'G':
			t.textView.Highlight(t.createTweetId(len(t.contents) - 1))
			return nil
		}
	}

	return event
}
