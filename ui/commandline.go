package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

// initCommandLine : コマンドラインを初期化
func (u *UI) initCommandLine() {
	u.commandLine.
		SetPlaceholderStyle(tcell.StyleDefault).
		SetBackgroundColor(tcell.ColorDefault)

	u.commandLine.
		SetChangedFunc(func(text string) {
			if text == "" {
				u.app.SetFocus(u.view.pagesView)
			}
		}).
		SetFocusFunc(func() {
			u.commandLine.SetText(":")
		})

	u.commandLine.SetInputCapture(u.handleCommandLineKeyEvents)
}

// updateStatusMessage : ステータスメッセージを更新
func (u *UI) updateStatusMessage(s string) {
	color := tcell.ColorDefault

	// エラーステータスなら文字色を赤に
	if strings.HasPrefix(s, "[ERR") {
		color = tcell.ColorRed
	}

	u.commandLine.
		SetPlaceholderTextColor(color).
		SetPlaceholder(s)

	u.app.Draw()
}

// handleCommandLineKeyEvents : コマンドラインのキーハンドラ
func (u *UI) handleCommandLineKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()

	// フォーカスをページへ移す
	if key == tcell.KeyEsc {
		u.commandLine.SetText("")
		u.app.SetFocus(u.view.pagesView)
		return nil

	}

	return event
}
