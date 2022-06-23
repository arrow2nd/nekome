package app

import "github.com/gdamore/tcell/v2"

// initModal：モーダル初期化
func (a *App) initModal() {
	a.modal.
		AddButtons([]string{"No", "Yes"}).
		SetTextColor(tcell.ColorDefault).
		SetBackgroundColor(tcell.ColorDefault)
}

// popupModal : モーダルを表示
func (a *App) popupModal(s string, f func()) {
	a.modal.
		SetText(s).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				f()
			}
			a.view.pagesView.RemovePage("modal")
		})

	a.view.pagesView.AddPage("modal", a.modal, true, true)

	a.app.SetFocus(a.modal)
}
