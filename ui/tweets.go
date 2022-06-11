package ui

import (
	"fmt"
	"sync"

	"github.com/g8rswimmer/go-twitter/v2"
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
	t.textView.Clear()

	for _, content := range t.contents {
		if len(content.ReferencedTweets) != 0 {
			fmt.Fprintf(t.textView, "type = %s\n", content.ReferencedTweets[0].Reference.Type)
		}

		fmt.Fprintf(t.textView, "[white::b]%s [gray::i]@%s\n", content.Author.Name, content.Author.UserName)
		fmt.Fprintf(t.textView, "[-:-:-]%s\n", content.Tweet.Text)

		likes := content.Tweet.PublicMetrics.Likes
		rts := content.Tweet.PublicMetrics.Retweets
		via := content.Tweet.Source
		fmt.Fprintf(t.textView, "[pink]%dLikes [green]%dRTs [-]- %s - via.%s\n", likes, rts, convertDateString(content.Tweet.CreatedAt), via)

		fmt.Fprintln(t.textView, "")
	}
}
