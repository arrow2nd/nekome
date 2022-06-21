package ui

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
	s := &statusBar{
		flex:          tview.NewFlex(),
		accountInfo:   tview.NewTextView(),
		pageIndicator: tview.NewTextView(),
	}

	s.accountInfo.SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorBlack).
		SetBackgroundColor(tcell.ColorDarkGray)

	s.pageIndicator.SetDynamicColors(true).
		SetTextAlign(tview.AlignRight).
		SetTextColor(tcell.ColorBlack).
		SetBackgroundColor(tcell.ColorDarkGray)

	s.flex.SetDirection(tview.FlexColumn).
		AddItem(s.accountInfo, 0, 1, false).
		AddItem(s.pageIndicator, 0, 1, false)

	return s
}

// DrawAccountInfo : ログイン中のアカウント情報を描画
func (s *statusBar) DrawAccountInfo() {
	s.accountInfo.Clear()
	fmt.Fprintf(s.accountInfo, " @%s", shared.api.CurrentUser.UserName)
}

// DrawPageIndicator : 現在のページのインジケータを描画
func (s *statusBar) DrawPageIndicator(d string) {
	s.pageIndicator.Clear()
	fmt.Fprintf(s.pageIndicator, "%s ", d)
}
