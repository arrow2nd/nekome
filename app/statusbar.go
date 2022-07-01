package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type statusBar struct {
	flex          *tview.Flex
	accountInfo   *tview.TextView
	pageIndicator *tview.TextView
}

func newStatusBar() *statusBar {
	return &statusBar{
		flex:          tview.NewFlex(),
		accountInfo:   tview.NewTextView(),
		pageIndicator: tview.NewTextView(),
	}
}

// Init : 初期化
func (s *statusBar) Init() {
	bgColor := tcell.NewHexColor(shared.conf.Theme.Statusbar.BG)

	s.accountInfo.
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetBackgroundColor(bgColor)

	s.pageIndicator.
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight).
		SetBackgroundColor(bgColor)

	s.flex.
		SetDirection(tview.FlexColumn).
		AddItem(s.accountInfo, 0, 1, false).
		AddItem(s.pageIndicator, 0, 1, false)
}

// DrawAccountInfo : ログイン中のアカウント情報を描画
func (s *statusBar) DrawAccountInfo() {
	s.accountInfo.Clear()

	fmt.Fprintf(
		s.accountInfo,
		" [%s]@%s[-:-:-]",
		shared.conf.Theme.Statusbar.Text,
		shared.api.CurrentUser.UserName,
	)
}

// DrawPageIndicator : 現在のページのインジケータを描画
func (s *statusBar) DrawPageIndicator(d string) {
	s.pageIndicator.Clear()

	fmt.Fprintf(
		s.pageIndicator,
		"[%s]%s[-:-:-] ",
		shared.conf.Theme.Statusbar.Text,
		d,
	)
}
