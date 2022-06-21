package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type statusBar struct {
	flex      *tview.Flex
	leftView  *tview.TextView
	rightView *tview.TextView
}

func newStatusBar() *statusBar {
	s := &statusBar{
		flex:      tview.NewFlex(),
		leftView:  tview.NewTextView(),
		rightView: tview.NewTextView(),
	}

	s.leftView.SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorBlack).
		SetBackgroundColor(tcell.ColorDarkGray)

	s.rightView.SetDynamicColors(true).
		SetTextAlign(tview.AlignRight).
		SetTextColor(tcell.ColorBlack).
		SetBackgroundColor(tcell.ColorDarkGray)

	s.flex.SetDirection(tview.FlexColumn).
		AddItem(s.leftView, 0, 1, false).
		AddItem(s.rightView, 0, 1, false)

	return s
}

// DrawLeft : 左側のビューを描画（認証前に呼ぶとエラー）
func (s *statusBar) DrawLeft() {
	s.leftView.Clear()
	fmt.Fprintf(s.leftView, " @%s", shared.api.CurrentUser.UserName)
}

// DrawRight : 右側のビューを描画
func (s *statusBar) DrawRight(d string) {
	s.rightView.Clear()
	fmt.Fprintf(s.rightView, "%s ", d)
}
