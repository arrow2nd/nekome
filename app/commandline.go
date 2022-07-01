package app

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

type commandLine struct {
	inputField        *tview.InputField
	autoComplateItems []string
	backspaceCount    int
	mu                sync.Mutex
}

func newCommandLine() *commandLine {
	return &commandLine{
		inputField:        tview.NewInputField(),
		backspaceCount:    0,
		autoComplateItems: []string{},
	}
}

// Init : 初期化
func (c *commandLine) Init() {
	normalBG := tcell.NewHexColor(shared.conf.Theme.Autocomplete.NormalBG)
	selectBG := tcell.NewHexColor(shared.conf.Theme.Autocomplete.SelectBG)

	c.inputField.
		SetAutocompleteStyles(
			normalBG,
			tcell.StyleDefault,
			tcell.StyleDefault.Background(selectBG),
		).
		SetPlaceholderStyle(tcell.StyleDefault).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetBackgroundColor(tcell.ColorDefault)

	c.inputField.
		SetAutocompleteFunc(c.handleAutocomplete).
		SetDoneFunc(c.handleDone).
		SetFocusFunc(c.handleFocus).
		SetInputCapture(c.handleKeyEvent)
}

// SetText : テキストを設定
func (c *commandLine) SetText(s string) {
	c.inputField.SetText(s)
}

// SetListCompleteItems : リストの補完要素を設定
func (c *commandLine) SetListCompleteItems(cmds []*cobra.Command) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.autoComplateItems = []string{}

	// 基本のコマンドを追加
	for _, cmd := range cmds {
		c.autoComplateItems = append(c.autoComplateItems, cmd.Name())
	}

	// ユーザが所有しているリストを取得
	lists, err := shared.api.FetchOwnedLists(shared.api.CurrentUser.ID)
	if err != nil {
		return err
	}

	// フラグ指定済みのリストコマンドを追加
	for _, l := range lists {
		cmd := fmt.Sprintf("list %s %s", l.Name, l.ID)
		c.autoComplateItems = append(c.autoComplateItems, cmd)
	}

	return nil
}

// UpdateStatusMessage : ステータスメッセージを更新
func (c *commandLine) UpdateStatusMessage(s string) {
	color := tcell.ColorDefault

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
	c.inputField.
		SetLabel("").
		SetText("")

	shared.RequestFocusPageView()
}

// handleAutocomplete : コマンドの入力補完ハンドラ
func (c *commandLine) handleAutocomplete(currentText string) []string {
	var entries []string = nil

	if currentText == "" {
		return nil
	}

	for _, cmd := range c.autoComplateItems {
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

		c.Blur()
	}
}

// handleFocus : フォーカス時のイベントハンドラ
func (c *commandLine) handleFocus() {
	c.inputField.
		SetLabelColor(tcell.ColorDefault).
		SetLabel(":").
		SetPlaceholder("")
}

// handleKeyEvent : キーイベントハンドラ
func (c *commandLine) handleKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	text := c.inputField.GetText()

	// フィールドが空かつ、BSが押されたらフォーカスを外す
	if text == "" && (key == tcell.KeyBackspace || key == tcell.KeyBackspace2) {
		c.backspaceCount++
		if c.backspaceCount >= 2 {
			c.Blur()
		}
		return nil
	}

	// フォーカスを外す
	if key == tcell.KeyEsc {
		c.Blur()
		return nil
	}

	// Tabキーを上キーの入力に変換
	// NOTE: デフォルトだとTabキーで補完候補の選択ができない
	if key == tcell.KeyTAB {
		return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	}

	return event
}
