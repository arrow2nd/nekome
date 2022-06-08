package ui

import (
	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI ユーザインターフェース
type UI struct {
	App    *tview.Application
	pages  *tview.Pages
	tabBar *tabBar
}

// New 生成
func New() *UI {
	return &UI{
		App:    tview.NewApplication(),
		pages:  tview.NewPages(),
		tabBar: newTabBar(),
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

	u.pages.AddPage("page_1", home.frame, true, true)

	u.tabBar.SetTab([]string{"Home", "Mention", "List", "Search", "User"})

	// 画面レイアウト
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(u.tabBar.textView, 2, 1, false).
		AddItem(u.pages, 0, 1, true)

	u.App.SetRoot(layout, true)

	// マウス操作有効化
	u.App.EnableMouse(true).
		SetBeforeDrawFunc(func(screen tcell.Screen) bool {
			screen.Clear()
			return false
		})
}
