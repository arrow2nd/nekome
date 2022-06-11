package ui

import "github.com/gdamore/tcell/v2"

func (u *UI) initCommandLine() {
	u.commandLine.SetFieldBackgroundColor(tcell.ColorDefault).
		SetPlaceholderStyle(tcell.StyleDefault)

	u.commandLine.SetFocusFunc(func() {
		u.commandLine.SetText(":")
	})

	u.commandLine.SetChangedFunc(func(text string) {
		if text == "" {
			u.app.SetFocus(u.view.pages)
		}
	})

	u.setCommandLineKeyEvent()
}

func (u *UI) setCommandLineKeyEvent() {
	u.commandLine.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			u.commandLine.SetText("")
			u.app.SetFocus(u.view.pages)
			return nil
		}

		return event
	})
}