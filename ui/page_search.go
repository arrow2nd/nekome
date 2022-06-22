package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type searchPage struct {
	*basePage
	query string
}

func newSearchPage(query string) *searchPage {
	p := &searchPage{
		basePage: newBasePage("Search: " + query),
		query:    query,
	}

	p.SetFrame(p.tweets.view)
	p.frame.SetInputCapture(p.handleKeyEvents)

	return p
}

// Load : 検索結果読み込み
func (s *searchPage) Load() {
	s.mu.Lock()
	defer s.mu.Unlock()

	shared.SetStatus(s.name, "Loading...")

	sinceID := s.tweets.GetSinceID()
	query := s.query + " -is:retweet"
	tweets, rateLimit, err := shared.api.SearchRecentTweets(query, sinceID, 25)

	if err != nil {
		shared.SetErrorStatus(s.name, err.Error())
		return
	}

	s.tweets.Register(tweets)
	s.tweets.Draw()

	s.updateIndicator(fmt.Sprintf("Query: %s | ", s.query), rateLimit)
	s.updateLoadedStatus(len(tweets))
}

// handleKeyEvents : 検索ページのキーハンドラ
func (s *searchPage) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	return handleCommonPageKeyEvent(s, event)
}
