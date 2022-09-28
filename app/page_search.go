package app

import (
	"fmt"
	"strings"
)

type searchPage struct {
	*tweetsBasePage
	query string
}

func newSearchPage(query string) (*searchPage, error) {
	tabName := strings.Replace(shared.conf.Pref.Text.TabSearch, "{query}", query, 1)

	basePage, err := newTweetsBasePage(tabName)
	if err != nil {
		return nil, err
	}
	p := &searchPage{
		tweetsBasePage: basePage,
		query:          query,
	}

	p.SetFrame(p.tweets.view)

	handler, err := createCommonPageKeyHandler(p)
	if err != nil {
		return nil, err
	}

	p.frame.SetInputCapture(handler)

	return p, nil
}

// Load : 検索結果読み込み
func (s *searchPage) Load() {
	s.mu.Lock()
	defer s.mu.Unlock()

	shared.SetStatus(s.name, shared.conf.Pref.Text.Loading)

	// ツイートを検索（RTは除外）
	count := shared.conf.Pref.Feature.LoadTweetsLimit
	sinceId := s.tweets.GetSinceID()
	query := s.query + " -is:retweet"
	tweets, rateLimit, err := shared.api.SearchRecentTweets(query, sinceId, count)

	if err != nil {
		s.tweets.DrawMessage(err.Error())
		shared.SetErrorStatus(s.name, err.Error())
		return
	}

	s.tweets.Update(tweets)
	s.tweets.UpdateRateLimit(rateLimit)

	s.updateIndicator(fmt.Sprintf("Query: %s | ", s.query))
	s.updateLoadedStatus(len(tweets))
}
