package app

import "strings"

type listPage struct {
	*tweetsBasePage
	listID string
}

func newListPage(name, id string) *listPage {
	tabName := shared.conf.Pref.Text.TabList
	tabName = strings.Replace(tabName, "{name}", name, 1)

	p := &listPage{
		tweetsBasePage: newTweetsBasePage(tabName),
		listID:         id,
	}

	p.SetFrame(p.tweets.view)
	p.frame.SetInputCapture(p.handleKeyEvents)

	return p
}

// Load : リスト読み込み
func (l *listPage) Load() {
	l.mu.Lock()
	defer l.mu.Unlock()

	shared.SetStatus(l.name, shared.conf.Pref.Text.Loading)

	// リスト内のツイートを取得
	count := shared.conf.Pref.Feature.LoadTweetsLimit
	tweets, rateLimit, err := shared.api.FetchListTweets(l.listID, count)
	if err != nil {
		l.tweets.DrawMessage(err.Error())
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

	// 新規ツイートがある場合のみ画面を更新
	if newTweetsCount > 0 {
		l.tweets.Update(tweets[0:newTweetsCount])
	}

	l.tweets.UpdateRateLimit(rateLimit)

	l.updateIndicator("")
	l.updateLoadedStatus(newTweetsCount)
}
