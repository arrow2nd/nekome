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
		SetAutocompleteFunc(a.handleCommandLineAutocomplete).
		SetChangedFunc(func(text string) {
			if text == "" {
				a.commandLine.SetLabel("")
				a.app.SetFocus(a.view.pageView)
			}
		}).
		SetFocusFunc(func() {
			a.commandLine.
				SetLabelColor(tcell.ColorDefault).
				SetPlaceholder("").
				SetLabel(":")
		}).
		SetInputCapture(a.handleCommandLineKeyEvents)
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
}

// handleCommandLineKeyEvents : コマンドラインのキーイベントハンドラ
func (a *App) handleCommandLineKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	text := a.commandLine.GetText()

	// コマンドを実行
	if key == tcell.KeyEnter {
		args := strings.Split(text, " ")

		if err := a.ExecCmd(args); err != nil {
			shared.SetErrorStatus("Command", err.Error())
		}

		a.blurCommandLine()

		return nil
	}

	// フォーカスをページへ移す
	if key == tcell.KeyEsc {
		a.blurCommandLine()
		return nil
	}

	return event
}

// handleCommandLineAutocomplete : コマンドの入力補完ハンドラ
func (a *App) handleCommandLineAutocomplete(currentText string) (entries []string) {
	currentText = strings.Replace(currentText, ":", "", 1)

	if currentText == "" {
		return nil
	}

	for _, cmd := range a.getCommands() {
		if strings.HasPrefix(strings.ToLower(cmd), strings.ToLower(currentText)) {
			entries = append(entries, cmd)
		}
	}

	if len(entries) == 0 {
		entries = nil
	}

	return
}
