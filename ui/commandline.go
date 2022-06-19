package ui

import (
	"github.com/gdamore/tcell/v2"
)

func (u *UI) initCommandLine() {
	u.commandLine.
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetPlaceholderStyle(tcell.StyleDefault).
		SetChangedFunc(func(text string) {
			if text == "" {
				u.app.SetFocus(u.view.pages)
			}
		}).
		SetFocusFunc(func() {
			u.commandLine.SetText(":")
		}).
		SetInputCapture(u.handleCommandLineKeyEvents)
}

func (u *UI) setStatusMessage(m string, c tcell.Color) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.commandLine.
		SetPlaceholderTextColor(c).
		SetPlaceholder(m)

	u.app.Draw()
}

func (u *UI) handleCommandLineKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()

	// フォーカスをページへ移す
	if key == tcell.KeyEsc {
		u.commandLine.SetText("")
		u.app.SetFocus(u.view.pages)
		return nil

	}

	return event
}
