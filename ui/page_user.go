package ui

import (
	"fmt"

	"github.com/arrow2nd/nekome/api"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const userProfilePaddingX = 4

type userPage struct {
	*basePage
	flex             *tview.Flex
	profile          *tview.TextView
	tweetMetrics     *tview.TextView
	followingMetrics *tview.TextView
	followersMetrics *tview.TextView
	userName         string
	userDic          *api.UserDictionary
}

func newUserPage(userName string) *userPage {
	p := &userPage{
		basePage:         newBasePage("@" + userName),
		flex:             tview.NewFlex(),
		profile:          tview.NewTextView(),
		tweetMetrics:     createMetricsView(0xa094c7),
		followingMetrics: createMetricsView(0x84a0c6),
		followersMetrics: createMetricsView(0x89b8c2),
		userName:         userName,
		userDic:          nil,
	}

	p.profile.SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignCenter).
		SetBorderPadding(0, 1, userProfilePaddingX, userProfilePaddingX)

	p.tweets.view.SetBorderPadding(1, 0, 0, 0)

	metrics := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(p.tweetMetrics, 0, 1, false).
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
func createMetricsView(color int32) *tview.TextView {
	t := tview.NewTextView().
		SetDynamicColors(true).
		SetTextColor(tcell.ColorBlack).
		SetTextAlign(tview.AlignCenter)

	t.SetBackgroundColor(tcell.NewHexColor(color))

	return t
}

// Load : ユーザタイムライン読み込み
func (u *userPage) Load() {
	shared.SetStatus(u.name, "Loading...")

	// ユーザの情報を取得
	if u.userDic == nil {
		if err := u.loadProfile(); err != nil {
			shared.SetErrorStatus(u.name, err.Error())
			return
		}
	}

	// ユーザのツイートを取得
	sinceID := u.tweets.GetSinceID()
	tweets, rateLimit, err := shared.api.FetchUserTimeline(u.userDic.User.ID, sinceID, 25)
	if err != nil {
		shared.SetErrorStatus(u.name, err.Error())
		return
	}

	u.drawProfile(u.userDic.User)

	if u.userDic.PinnedTweet != nil {
		u.tweets.RegisterPinned(u.userDic.PinnedTweet)
	}

	u.tweets.Register(tweets)
	u.tweets.Draw()

	u.showLoadedStatus(len(tweets), rateLimit)
}

// loadProfile : プロフィール情報の読み込み
func (u *userPage) loadProfile() error {
	users, err := shared.api.FetchUser([]string{u.userName})
	if err != nil {
		return err
	}

	if len(users) == 0 || users[0] == nil {
		return err
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

	// プロフィールの行数に合わせてリサイズ（+1 は下辺の padding 分）
	u.flex.ResizeItem(u.profile, row+1, 1)

	// ツイート・フォロイー・フォロワー数
	u.tweetMetrics.SetText(fmt.Sprintf("%d Tweets", ur.PublicMetrics.Tweets))
	u.followingMetrics.SetText(fmt.Sprintf("%d Following", ur.PublicMetrics.Following))
	u.followersMetrics.SetText(fmt.Sprintf("%d Followers", ur.PublicMetrics.Followers))
}

// handleKeyEvents : ユーザページのキーハンドラ
func (u *userPage) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	return handlePageKeyEvents(u, event)
}
