package app

import (
	"fmt"
	"strings"

	"github.com/arrow2nd/nekome/api"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type userPage struct {
	*tweetsBasePage
	flex             *tview.Flex
	profile          *tview.TextView
	tweetsMetrics    *tview.TextView
	followingMetrics *tview.TextView
	followersMetrics *tview.TextView
	userName         string
	userDic          *api.UserDictionary
}

func newUserPage(userName string) *userPage {
	tabName := shared.conf.Settings.Text.TabUser
	tabName = strings.Replace(tabName, "{name}", userName, 1)

	tweetsColor := shared.conf.Style.Metrics.TweetsBackgroundColor.ToColor()
	followingColor := shared.conf.Style.Metrics.FollowingBackgroundColor.ToColor()
	followersColor := shared.conf.Style.Metrics.FollowersBackgroundColor.ToColor()

	p := &userPage{
		tweetsBasePage:   newTweetsBasePage(tabName),
		flex:             tview.NewFlex(),
		profile:          tview.NewTextView(),
		tweetsMetrics:    createMetricsView(tweetsColor),
		followingMetrics: createMetricsView(followingColor),
		followersMetrics: createMetricsView(followersColor),
		userName:         userName,
		userDic:          nil,
	}

	padding := shared.conf.Settings.Appearance.UserProfilePaddingX

	// プロフィール表示域
	p.profile.
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignCenter).
		SetBorderPadding(0, 1, padding, padding)

	// メトリクス表示域
	metrics := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(p.tweetsMetrics, 0, 1, false).
		AddItem(p.followingMetrics, 0, 1, false).
		AddItem(p.followersMetrics, 0, 1, false)

	p.tweets.view.SetBorderPadding(1, 0, 0, 0)

	p.flex.
		SetDirection(tview.FlexRow).
		AddItem(p.profile, 1, 1, false).
		AddItem(metrics, 1, 1, false).
		AddItem(p.tweets.view, 0, 1, true)

	p.SetFrame(p.flex)

	p.frame.SetInputCapture(p.handleKeyEvents)

	return p
}

// createMetricsView : 各メトリクス表示用のTextViewを作成
func createMetricsView(color tcell.Color) *tview.TextView {
	t := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	t.SetBackgroundColor(color)

	return t
}

// Load : ユーザタイムライン読み込み
func (u *userPage) Load() {
	u.mu.Lock()
	defer u.mu.Unlock()

	shared.SetStatus(u.name, shared.conf.Settings.Text.Loading)

	// ユーザの情報を取得
	if u.userDic == nil {
		if err := u.loadProfile(); err != nil {
			u.tweets.DrawMessage(err.Error())
			shared.SetErrorStatus(u.name, err.Error())
			return
		}
	}

	// ユーザのツイートを取得
	tweets, rateLimit, err := shared.api.FetchUserTimeline(
		u.userDic.User.ID,
		u.tweets.GetSinceID(),
		shared.conf.Settings.Feature.LoadTweetsLimit,
	)

	if err != nil {
		u.tweets.DrawMessage(err.Error())
		shared.SetErrorStatus(u.name, err.Error())
		return
	}

	u.drawProfile(u.userDic.User)

	if u.userDic.PinnedTweet != nil {
		u.tweets.RegisterPinned(u.userDic.PinnedTweet)
	}

	u.tweets.Update(tweets)
	u.tweets.UpdateRateLimit(rateLimit)

	u.updateIndicator("")
	u.updateLoadedStatus(len(tweets))
}

// loadProfile : プロフィール読み込み
func (u *userPage) loadProfile() error {
	users, err := shared.api.FetchUser([]string{u.userName})
	if err != nil {
		return err
	}

	if len(users) == 0 || users[0] == nil {
		return fmt.Errorf("no user profile data: %w", err)
	}

	u.userDic = users[0]

	return nil
}

// drawProfile : プロフィールを描画
func (u *userPage) drawProfile(ur *twitter.UserObj) {
	width := getWindowWidth()

	u.profile.Clear()

	// プロフィール
	profile, row := createProfileLayout(ur, width)
	fmt.Fprint(u.profile, profile)

	// プロフィールの行数に合わせて表示域をリサイズ（+1 は下辺の padding 分）
	u.flex.ResizeItem(u.profile, row+1, 1)

	style := shared.conf.Style.Metrics

	// ツイート数
	u.tweetsMetrics.SetText(fmt.Sprintf(
		"[%s]%d Tweets[-:-:-]",
		style.TweetsText,
		ur.PublicMetrics.Tweets,
	))

	// フォロイー数
	u.followingMetrics.SetText(fmt.Sprintf(
		"[%s]%d Following[-:-:-]",
		style.FollowingText,
		ur.PublicMetrics.Following,
	))

	// フォロワー数
	u.followersMetrics.SetText(fmt.Sprintf(
		"[%s]%d Followers[-:-:-]",
		style.FollowersText,
		ur.PublicMetrics.Followers,
	))
}
