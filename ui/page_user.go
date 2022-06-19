package ui

import (
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type userPage struct {
	*basePage
	detail *tview.TextView
	// tweetCount     *tview.TextView
	// followingCount *tview.TextView
	// followersCount *tview.TextView
	targetUserID string
	user         *twitter.UserObj
}

func newUserPage(name, userID string) *userPage {
	page := &userPage{
		basePage: newBasePage(name),
		detail:   tview.NewTextView(),
		// tweetCount:     newCountTextView(0xe06c75),
		// followingCount: newCountTextView(0xc678dd),
		// followersCount: newCountTextView(0x56b6c2),
		targetUserID: userID,
		user:         nil,
	}

	page.detail.SetDynamicColors(true).
		SetWrap(true)

	page.tweets.view.SetBorderPadding(1, 0, 0, 0)

	// countView := tview.NewFlex().
	// 	AddItem(page.tweetCount, 0, 1, false).
	// 	AddItem(page.followingCount, 0, 1, false).
	// 	AddItem(page.followersCount, 0, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(page.detail, 5, 1, false).
		AddItem(page.tweets.view, 0, 1, true)

	page.SetFrame(layout)

	page.frame.SetInputCapture(page.handleKeyEvents)

	return page
}

func newCountTextView(color int32) *tview.TextView {
	t := tview.NewTextView().
		SetDynamicColors(true).
		SetTextColor(tcell.ColorBlack).
		SetTextAlign(tview.AlignCenter)

	t.SetBackgroundColor(tcell.NewHexColor(color))

	return t
}

func (u *userPage) Load() {
	shared.setStatus(u.name, "Loading...")

	sinceID := u.tweets.getSinceID()
	tweets, rateLimit, err := shared.api.FetchUserTimeline(u.targetUserID, sinceID, 25)
	if err != nil {
		shared.setErrorStatus(u.name, err.Error())
	}

	u.tweets.register(tweets)
	u.tweets.draw()

	if u.tweets.count != 0 {
		u.setUserDetail(tweets[0].Author)
	}

	u.showLoadedStatus(rateLimit)
}

func (u *userPage) setUserDetail(user *twitter.UserObj) {
	u.detail.Clear()

	fmt.Fprint(u.detail, createHeader(user, -1))
	fmt.Fprintf(u.detail, "[white:-:-]%s\n", user.Description)
	fmt.Fprintf(u.detail, " : %s\n", user.Location)
	fmt.Fprintf(u.detail, " : %s\n", user.URL)
	fmt.Fprintf(u.detail, "%d Tweets / %d Following / %d Followers", user.PublicMetrics.Tweets, user.PublicMetrics.Following, user.PublicMetrics.Followers)

	// u.tweetCount.SetText(fmt.Sprintf("%d Tweets", user.PublicMetrics.Tweets))
	// u.followingCount.SetText(fmt.Sprintf("%d Following", user.PublicMetrics.Following))
	// u.followersCount.SetText(fmt.Sprintf("%d Followers", user.PublicMetrics.Followers))
}

func (u *userPage) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	return handlePageKeyEvents(u, event)
}
