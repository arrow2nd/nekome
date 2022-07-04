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
	tabName := shared.conf.Settings.Texts.TabUser
	tabName = strings.Replace(tabName, "{name}", userName, 1)

	tweetsColor := shared.conf.Style.User.TweetsMetricsBG.ToColor()
	followingColor := shared.conf.Style.User.FollowingMetricsBG.ToColor()
	followersColor := shared.conf.Style.User.FollowersMetricsBG.ToColor()

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

	p.profile.
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignCenter).
		SetBorderPadding(0, 1, padding, padding)

	p.tweets.view.SetBorderPadding(1, 0, 0, 0)

	metrics := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(p.tweetsMetrics, 0, 1, false).
		AddItem(p.followingMetrics, 0, 1, false).
		AddItem(p.followersMetrics, 0, 1, false)

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

	shared.SetStatus(u.name, shared.conf.Settings.Texts.Loading)

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
		shared.conf.Settings.Feature.LoadTweetsCount,
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

	u.tweets.Register(tweets)
	u.tweets.Draw()

	u.updateIndicator("", rateLimit)
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

	// ツイート数
	tweets := fmt.Sprintf(
		"[%s]%d Tweets[-:-:-]",
		shared.conf.Style.User.TweetsMetricsText,
		ur.PublicMetrics.Tweets,
	)
	u.tweetsMetrics.SetText(tweets)

	// フォロイー数
	following := fmt.Sprintf(
		"[%s]%d Following[-:-:-]",
		shared.conf.Style.User.FollowingMetricsText,
		ur.PublicMetrics.Following,
	)
	u.followingMetrics.SetText(following)

	// フォロワー数
	followers := fmt.Sprintf(
		"[%s]%d Followers[-:-:-]",
		shared.conf.Style.User.FollowersMetricsText,
		ur.PublicMetrics.Followers,
	)
	u.followersMetrics.SetText(followers)
}

// handleKeyEvents : ユーザページのキーハンドラ
func (u *userPage) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	return handleCommonPageKeyEvent(u, event)
}
