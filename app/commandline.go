package app

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

// initCommandLine : コマンドラインを初期化
func (a *App) initCommandLine() {
	a.commandLine.
		SetPlaceholderStyle(tcell.StyleDefault).
		SetBackgroundColor(tcell.ColorDefault)

	a.commandLine.
		SetChangedFunc(func(text string) {
			if text == "" {
				a.app.SetFocus(a.view.pagesView)
			}
		}).
		SetFocusFunc(func() {
			a.commandLine.SetText(":")
		})

	a.commandLine.SetInputCapture(a.handleCommandLineKeyEvents)
}

// updateStatusMessage : ステータスメッセージを更新
func (a *App) updateStatusMessage(s string) {
	color := tcell.ColorDefault

	// エラーステータスなら文字色を赤に
	if strings.HasPrefix(s, "[ERR") {
		color = tcell.ColorRed
	}

	a.commandLine.
		SetPlaceholderTextColor(color).
		SetPlaceholder(s)

	a.app.Draw()
}

// handleCommandLineKeyEvents : コマンドラインのキーハンドラ
func (a *App) handleCommandLineKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()

	// フォーカスをページへ移す
	if key == tcell.KeyEsc {
		a.commandLine.SetText("")
		a.app.SetFocus(a.view.pagesView)
		return nil

	}

	return event
}
