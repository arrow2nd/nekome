package app

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type commandLine struct {
	inputField        *tview.InputField
	listComplateItems map[string]string
}

func newCommandLine() *commandLine {
	c := &commandLine{
		inputField:        tview.NewInputField(),
		listComplateItems: map[string]string{},
	}

	c.inputField.
		SetAutocompleteStyles(
			tcell.NewHexColor(0x3e4359),
			tcell.StyleDefault,
			tcell.StyleDefault.Background(tcell.NewHexColor(0x5c6586)),
		).
		SetPlaceholderStyle(tcell.StyleDefault).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetBackgroundColor(tcell.ColorDefault)

	c.inputField.
		SetAutocompleteFunc(c.handleAutocomplete).
		SetDoneFunc(c.handleDone).
		SetChangedFunc(c.handleChanged).
		SetFocusFunc(c.handleFocus).
		SetInputCapture(c.handleKeyEvent)

	return c
}

// updateStatusMessage : ステータスメッセージを更新
func (c *commandLine) updateStatusMessage(s string) {
	color := tcell.ColorDefault

	// エラーステータスなら文字色を赤に
	if strings.HasPrefix(s, "[ERR") {
		color = tcell.ColorRed
	}

	c.inputField.
		SetPlaceholderTextColor(color).
		SetPlaceholder(s)
}

// blurCommandLine : コマンドラインからフォーカスを外す
func (c *commandLine) blurCommandLine() {
	c.inputField.SetText("")
}

// handleAutocomplete : コマンドの入力補完ハンドラ
func (c *commandLine) handleAutocomplete(currentText string) (entries []string) {
	if currentText == "" {
		return nil
	}

	for _, cmd := range getCommands() {
		if strings.HasPrefix(strings.ToLower(":"+cmd), strings.ToLower(currentText)) {
			entries = append(entries, cmd)
		}
	}

	if len(entries) == 0 {
		entries = nil
	}

	return
}

// handleDone : 入力確定時のイベントハンドラ
func (c *commandLine) handleDone(key tcell.Key) {
	// コマンドを実行
	if key == tcell.KeyEnter {
		text := strings.Replace(c.inputField.GetText(), ":", "", 1)
		shared.RequestExecCommand(text)

		c.blurCommandLine()
	}
}

// handleChanged : フィールド変更時のイベントハンドラ
func (c *commandLine) handleChanged(text string) {
	// フィールドが空ならフォーカスを外す
	if text == "" {
		shared.RequestFocusPageView()
		return
	}

	// 先頭に ":" が無ければ追加
	if !strings.HasPrefix(text, ":") {
		c.inputField.SetText(":" + text)
	}
}

// handleFocus : フォーカス時のイベントハンドラ
func (c *commandLine) handleFocus() {
	c.inputField.
		SetLabelColor(tcell.ColorDefault).
		// SetPlaceholder("").
		SetText(":")
}

// handleKeyEvent : キーイベントハンドラ
func (c *commandLine) handleKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()

	// フォーカスをページへ移す
	if key == tcell.KeyEsc {
		c.blurCommandLine()
		return nil
	}

	return event
}
