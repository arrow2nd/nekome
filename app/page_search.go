package app

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type searchPage struct {
	*tweetsBasePage
	query string
}

// newSearchPage : 検索ページを作成
func newSearchPage(query string) *searchPage {
	tabName := shared.conf.Settings.Texts.TabSearch
	tabName = strings.Replace(tabName, "{query}", query, 1)

	p := &searchPage{
		tweetsBasePage: newTweetsBasePage(tabName),
		query:          query,
	}

	p.SetFrame(p.tweets.view)
	p.frame.SetInputCapture(p.handleKeyEvents)

	return p
}

// Load : 検索結果読み込み
func (s *searchPage) Load() {
	s.mu.Lock()
	defer s.mu.Unlock()

	shared.SetStatus(s.name, shared.conf.Settings.Texts.Loading)

	// ツイートを検索（RTは除外）
	count := shared.conf.Settings.Feature.LoadTweetsCount
	sinceID := s.tweets.GetSinceID()
	query := s.query + " -is:retweet"
	tweets, rateLimit, err := shared.api.SearchRecentTweets(query, sinceID, count)

	if err != nil {
		s.tweets.DrawMessage(err.Error())
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
