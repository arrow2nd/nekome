package ui

import (
	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI ユーザインターフェース
type UI struct {
	App  *tview.Application
	view *view
}

// New 生成
func New() *UI {
	return &UI{
		App:  tview.NewApplication(),
		view: newView(),
	}
}

// Init 初期化
func (u *UI) Init(a *api.API, c *config.Config) {
	// 共有
	shared.api = a
	shared.conf = c

	// 配色設定
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	tview.Styles.ContrastBackgroundColor = tcell.ColorDefault

	// NOTE: テスト用
	home := newHomeTimeline()
	home.init()
	u.view.addPage("Home", home.frame, true)
	u.view.addPage("Mention", home.frame, false)

	// 画面レイアウト
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(u.view.tabTextView, 2, 1, false).
		AddItem(u.view.pages, 0, 1, true)

	u.App.SetRoot(layout, true)

	// マウス操作有効化
	u.App.EnableMouse(true).
		SetBeforeDrawFunc(func(screen tcell.Screen) bool {
			screen.Clear()
			return false
		})

	u.setCommonKeyEvent()
}

func (u *UI) setCommonKeyEvent() {
	u.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		// タブ切り替え(左右キー)
		case tcell.KeyLeft:
			u.view.selectPrevTab()
			return nil
		case tcell.KeyRight:
			u.view.selectNextTab()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			// タブ切り替え(hl)
			case 'h':
				u.view.selectPrevTab()
				return nil
			case 'l':
				u.view.selectNextTab()
				return nil
			}
		}

		return event
	})
}
