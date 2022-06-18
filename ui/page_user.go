package ui

import (
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type userPage struct {
	*basePage
	detail         *tview.TextView
	tweetCount     *tview.TextView
	followingCount *tview.TextView
	followersCount *tview.TextView
	user           *twitter.UserObj
}

func newUserPage() *userPage {
	page := &userPage{
		basePage:       newBasePage(),
		detail:         tview.NewTextView(),
		tweetCount:     newCountTextView(0xe06c75),
		followingCount: newCountTextView(0xc678dd),
		followersCount: newCountTextView(0x56b6c2),
		user:           nil,
	}

	page.detail.SetDynamicColors(true).
		SetWrap(true)

	countView := tview.NewFlex().
		AddItem(page.tweetCount, 0, 1, false).
		AddItem(page.followingCount, 0, 1, false).
		AddItem(page.followersCount, 0, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(page.detail, 4, 1, false).
		AddItem(countView, 1, 1, false).
		AddItem(page.tweets.view, 0, 1, false)

	page.SetFrame(layout)

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
	u.setUserDetail("ユーザー名", "user", "[blue]following", "BIO", "place")
	u.setCounts("10000", "10000", "10000")
}

func (u *userPage) setUserDetail(name, userName, relation, bio, place string) {
	u.detail.Clear()

	fmt.Fprintf(u.detail, "[white:-:b]%s\n", name)
	fmt.Fprintf(u.detail, "[gray:-:i]@%s[-:-:-] %s\n", userName, relation)
	fmt.Fprintf(u.detail, "[white:-:-]📄 : %s\n", bio)
	fmt.Fprintf(u.detail, "📍 : %s", place)
}

func (u *userPage) setCounts(tweets, following, followers string) {
	u.tweetCount.SetText(fmt.Sprintf("%s Tweets", tweets))
	u.followingCount.SetText(fmt.Sprintf("%s Following", following))
	u.followersCount.SetText(fmt.Sprintf("%s Followers", followers))
}
