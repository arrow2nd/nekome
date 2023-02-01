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
		autocomplateItems:        nil,
		isAutocompleteDisplaying: false,
		backspaceCount:           0,
	}
}

// Init : 初期化
func (c *commandLine) Init() {
	style := shared.conf.Style
	acTextColor := style.Autocomplate.TextColor.ToColor()

	c.inputField.
		SetAutocompleteStyles(
			style.Autocomplate.BackgroundColor.ToColor(),
			tcell.StyleDefault.
				Foreground(acTextColor),
			tcell.StyleDefault.
				Foreground(acTextColor).
				Background(style.Autocomplate.SelectedBackgroundColor.ToColor()),
		).
		SetLabelColor(style.App.TextColor.ToColor())

	c.inputField.
		SetAutocompleteFunc(c.getAutocompleteItems).
		SetAutocompletedFunc(c.handleAutocompleted).
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

	sort.Slice(c.autocomplateItems, func(i, j int) bool {
		return c.autocomplateItems[i] < c.autocomplateItems[j]
	})

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
func (c *commandLine) Blur() {
	c.inputField.SetLabel("").SetText("")

	shared.RequestFocusView()
	shared.SetDisableViewKeyEvent(false)
}

// getAutocompleteItems : 入力補完の候補を取得
func (c *commandLine) getAutocompleteItems(currentText string) []string {
	entries := []string{}
	c.isAutocompleteDisplaying = true

	if currentText == "" {
		// 入力中の場合のみ全ての候補を返す
		if c.inputField.GetLabel() == ":" {
			return c.autocomplateItems
		}
		return nil
	}

	for _, cmd := range c.autocomplateItems {
		if strings.HasPrefix(strings.ToLower(cmd), strings.ToLower(currentText)) {
			entries = append(entries, cmd)
		}
	}

	return entries
}

// handleAutocompleted : 補完候補選択時のイベントハンドラ
func (c *commandLine) handleAutocompleted(text string, index, source int) bool {
	// 選択中でないなら、コマンドラインの内容を変更する
	if source != tview.AutocompletedNavigate {
		c.inputField.SetText(text)
	}

	autocompleted := source == tview.AutocompletedEnter
	c.isAutocompleteDisplaying = !autocompleted

	return autocompleted
}

// handleDone : 入力確定時のイベントハンドラ
func (c *commandLine) handleDone(key tcell.Key) {
	if key != tcell.KeyEnter {
		return
	}

	text := c.inputField.GetText()

	c.Blur()

	if text != "" {
		shared.RequestExecCommand(text)
	}
}

// handleFocus : フォーカス時のイベントハンドラ
func (c *commandLine) handleFocus() {
	c.inputField.
		SetLabel(":").
		SetPlaceholder("")

	shared.SetDisableViewKeyEvent(true)
}

// handleKeyEvent : キーイベントハンドラ
func (c *commandLine) handleKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	text := c.inputField.GetText()

	// フィールドが空かつ、BSが押されたらフォーカスを外す
	if text == "" && (key == tcell.KeyBackspace || key == tcell.KeyBackspace2 || key == tcell.KeyCtrlW) {
		c.Blur()
		return nil
	}

	// ESCでフォーカスを外す
	if key == tcell.KeyEsc {
		c.Blur()
		return nil
	}

	if key == tcell.KeyTab {
		return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	}

	if key == tcell.KeyBacktab {
		return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
	}

	// 補完候補を決定
	if c.isAutocompleteDisplaying && key == tcell.KeyCtrlY {
		return tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	}

	return event
}
