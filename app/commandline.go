package app

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

// initCommandLine : コマンドラインを初期化
func (a *App) initCommandLine() {
	a.commandLine.
		SetPlaceholderStyle(tcell.StyleDefault).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetBackgroundColor(tcell.ColorDefault)

	a.commandLine.
		SetChangedFunc(func(text string) {
			if text == "" {
				a.app.SetFocus(a.view.pageView)
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

// blurCommandLine : コマンドラインからフォーカスを外す
func (a *App) blurCommandLine() {
	a.commandLine.SetText("")
	a.app.SetFocus(a.view.pageView)
}

// handleCommandLineKeyEvents : コマンドラインのキーハンドラ
func (a *App) handleCommandLineKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()

	// フォーカスをページへ移す
	if key == tcell.KeyEsc {
		a.blurCommandLine()
		return nil

	}

	return event
}
