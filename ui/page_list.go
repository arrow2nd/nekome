package ui

import "github.com/gdamore/tcell/v2"

type listPage struct {
	*basePage
	listID string
}

func newListPage(listID, listName string) *listPage {
	p := &listPage{
		basePage: newBasePage(listName),
		listID:   listID,
	}

	p.SetFrame(p.tweets.view)
	p.frame.SetInputCapture(p.handleKeyEvents)

	return p
}

// Load : リスト読み込み
func (l *listPage) Load() {
	l.mu.Lock()
	defer l.mu.Unlock()

	shared.SetStatus(l.name, "Loading...")

	tweets, rateLimit, err := shared.api.FetchListTweets(l.listID, 25)
	if err != nil {
		shared.SetErrorStatus(l.name, err.Error())
		return
	}

	sinceID := l.tweets.GetSinceID()

	// 新規ツイート数をカウント
	newTweetsCount := 0
	for ; newTweetsCount < len(tweets); newTweetsCount++ {
		if tweets[newTweetsCount].Tweet.ID == sinceID {
			break
		}
	}

	// 新規ツイートのみを登録
	if newTweetsCount > 0 {
		l.tweets.Register(tweets[0:newTweetsCount])
		l.tweets.Draw()
	}

	l.updateIndicator("", rateLimit)
	l.updateLoadedStatus(newTweetsCount)
}

// handleKeyEvents : リストページのキーハンドラ
func (l *listPage) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	return handleCommonPageKeyEvent(l, event)
}
