package app

import (
	"fmt"
	"strings"
)

type likesPage struct {
	*tweetsBasePage
	userId   string
	userName string
}

func newLikesPage(userName string) (*likesPage, error) {
	pref := shared.conf.Pref

	tabName := strings.Replace(pref.Text.TabLikes, "{name}", userName, 1)
	basePage, err := newTweetsBasePage(tabName)
	if err != nil {
		return nil, err
	}

	p := &likesPage{
		tweetsBasePage: basePage,
		userId:         "",
		userName:       userName,
	}

	p.SetFrame(p.tweets.view)

	handler, err := createCommonPageKeyHandler(p)
	if err != nil {
		return nil, err
	}

	p.frame.SetInputCapture(handler)

	return p, nil
}

// Load : いいね済みツイート読み込み
func (l *likesPage) Load() {
	l.mu.Lock()
	defer l.mu.Unlock()

	shared.SetStatus(l.name, shared.conf.Pref.Text.Loading)

	// ユーザIDが空なら取得
	if l.userId == "" {
		if err := l.fetchUserId(); err != nil {
			l.tweets.DrawMessage(err.Error())
			shared.SetErrorStatus(l.name, err.Error())
			return
		}
	}

	count := shared.conf.Pref.Feature.LoadTweetsLimit
	tweets, rateLimit, err := shared.api.FetchLikedTweets(l.userId, count)
	if err != nil {
		l.tweets.DrawMessage(err.Error())
		shared.SetErrorStatus(l.name, err.Error())
		return
	}

	// 新規いいね済みツイート数をカウント
	newLikesCount := 0
	sinceId := l.tweets.GetSinceId()
	for ; newLikesCount < len(tweets); newLikesCount++ {
		if tweets[newLikesCount].Tweet.ID == sinceId {
			break
		}
	}

	// 新規のいいね済みツイートがある場合のみ画面を更新
	if newLikesCount > 0 {
		l.tweets.Update(tweets[0:newLikesCount])
	}

	l.tweets.UpdateRateLimit(rateLimit)

	l.updateIndicator("")
	l.updateLoadedStatus(newLikesCount)
}

// fetchUserId : ユーザIDを取得
func (l *likesPage) fetchUserId() error {
	users, err := shared.api.FetchUser([]string{l.userName})
	if err != nil {
		return err
	}

	if len(users) == 0 || users[0] == nil {
		return fmt.Errorf("failed to get user id: %w", err)
	}

	l.userId = users[0].User.ID

	return nil
}
