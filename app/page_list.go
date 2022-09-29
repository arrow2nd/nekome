package app

import "strings"

type listPage struct {
	*tweetsBasePage
	listId string
}

func newListPage(name, id string) (*listPage, error) {
	tabName := strings.Replace(shared.conf.Pref.Text.TabList, "{name}", name, 1)

	basePage, err := newTweetsBasePage(tabName)
	if err != nil {
		return nil, err
	}

	p := &listPage{
		tweetsBasePage: basePage,
		listId:         id,
	}

	p.SetFrame(p.tweets.view)

	handler, err := createCommonPageKeyHandler(p)
	if err != nil {
		return nil, err
	}

	p.frame.SetInputCapture(handler)

	return p, nil
}

// Load : リスト読み込み
func (l *listPage) Load() {
	pref := shared.conf.Pref

	l.mu.Lock()
	defer l.mu.Unlock()

	shared.SetStatus(l.name, pref.Text.Loading)

	// リスト内のツイートを取得
	count := pref.Feature.LoadTweetsLimit
	tweets, rateLimit, err := shared.api.FetchListTweets(l.listId, count)
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
