package app

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type commandLine struct {
	inputField               *tview.InputField
	backspaceCount           int
	autocomplateItems        []string
	isAutocompleteDisplaying bool
	mu                       sync.Mutex
}

func newCommandLine() *commandLine {
	return &commandLine{
		inputField:               tview.NewInputField(),
		autocomplateItems:        []string{},
		isAutocompleteDisplaying: false,
		backspaceCount:           0,
	}
}

// Init : 初期化
func (c *commandLine) Init() {
	style := shared.conf.Style.Autocomplate
	textColor := style.TextColor.ToColor()

	c.inputField.
		SetAutocompleteStyles(
			style.BackgroundColor.ToColor(),
			tcell.StyleDefault.
				Foreground(textColor),
			tcell.StyleDefault.
				Foreground(textColor).
				Background(style.SelectedBackgroundColor.ToColor()),
		).
		SetLabelColor(textColor).
		SetAutocompleteFunc(c.handleAutocomplete).
		SetDoneFunc(c.handleDone).
		SetFocusFunc(c.handleFocus).
		SetInputCapture(c.handleKeyEvent)
}

// SetText : テキストを設定
func (c *commandLine) SetText(s string) {
	c.inputField.SetText(s)
}

// SetAutocompleteItems : 補完要素を設定
func (c *commandLine) SetAutocompleteItems(cmds []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.autocomplateItems = []string{}

	// 基本のコマンドを追加
	for _, cmd := range cmds {
		c.autocomplateItems = append(c.autocomplateItems, cmd)
	}

	// ユーザが所有しているリストを取得
	lists, err := shared.api.FetchOwnedLists(shared.api.CurrentUser.ID)
	if err != nil {
		return err
	}

	// フラグ指定済みのlistコマンドを追加
	for _, l := range lists {
		cmd := fmt.Sprintf("list %s %s", l.Name, l.ID)
		c.autocomplateItems = append(c.autocomplateItems, cmd)
	}

	// ユーザ指定済みのaccount switchコマンドを追加
	for _, u := range shared.conf.Cred.GetAllNames() {
		cmd := fmt.Sprintf("account switch %s", u)
		c.autocomplateItems = append(c.autocomplateItems, cmd)
	}

	sort.Slice(c.autocomplateItems, func(i, j int) bool { return c.autocomplateItems[i] < c.autocomplateItems[j] })

	return nil
}

// UpdateStatusMessage : ステータスメッセージを更新
func (c *commandLine) UpdateStatusMessage(s string) {
	color := tview.Styles.PrimaryTextColor

	// エラーステータスなら文字色を赤に
	if strings.HasPrefix(s, "[ERR") {
		color = tcell.ColorRed
	}

	c.inputField.
		SetPlaceholderTextColor(color).
		SetPlaceholder(s)
}

// Blur : コマンドラインからフォーカスを外す
func (c *commandLine) Blur(closeAutocompleteList bool) {
	c.inputField.SetLabel("").SetText("")

	// 補完リストを閉じる
	// NOTE: 補完リストが表示された状態でフォーカスを外すと、一部が表示されたままになる
	if closeAutocompleteList {
		c.inputField.Autocomplete()
	}

	shared.RequestFocusView()
}

// handleAutocomplete : コマンドの入力補完ハンドラ
func (c *commandLine) handleAutocomplete(currentText string) []string {
	entries := []string{}
	c.isAutocompleteDisplaying = true

	if currentText == "" {
		// NOTE: 外部から補完リストを初期化できないので、空リストを返すことで補完リストを削除し表示をクリアする
		// LINK: https://github.com/rivo/tview/blob/2e69b7385a37df55e0c2ef4d1c0054898bed05a1/inputfield.go#L286-L292
		if c.inputField.GetLabel() == "" {
			c.isAutocompleteDisplaying = false
			return []string{}
		}

		return c.autocomplateItems
	}

	for _, cmd := range c.autocomplateItems {
		if strings.HasPrefix(strings.ToLower(cmd), strings.ToLower(currentText)) {
			entries = append(entries, cmd)
		}
	}

	return entries
}

// handleDone : 入力確定時のイベントハンドラ
func (c *commandLine) handleDone(key tcell.Key) {
	if key == tcell.KeyEnter {
		// コマンドを実行
		if text := c.inputField.GetText(); text != "" {
			shared.RequestExecCommand(text)
		}

		c.Blur(false)
	}
}

// handleFocus : フォーカス時のイベントハンドラ
func (c *commandLine) handleFocus() {
	c.inputField.
		SetLabel(":").
		SetPlaceholder("")
}

// handleKeyEvent : キーイベントハンドラ
func (c *commandLine) handleKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	text := c.inputField.GetText()

	// フィールドが空かつ、BSが押されたらフォーカスを外す
	if text == "" && (key == tcell.KeyBackspace || key == tcell.KeyBackspace2 || key == tcell.KeyCtrlW) {
		c.Blur(true)
		return nil
	}

	// フォーカスを外す
	if key == tcell.KeyEsc {
		c.Blur(true)
		return nil
	}

	if key == tcell.KeyTab {
		if c.isAutocompleteDisplaying {
			// 補完リスト表示中なら、上キーの入力に変換
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		} else {
			// 補完リストを表示
			c.inputField.Autocomplete()
			return nil
		}
	}

	// 項目を決定
	if key == tcell.KeyCtrlY {
		return tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	}

	return event
}
