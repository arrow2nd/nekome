package app

import (
	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
	"github.com/arrow2nd/nekome/log"
	"github.com/rivo/tview"
)

// Shared : 全体共有
type Shared struct {
	isCommandLineMode     bool
	api                   *api.API
	conf                  *config.Config
	chStatus              chan string
	chIndicator           chan string
	chPopupModal          chan *ModalOpt
	chExecCommand         chan string
	chInputCommand        chan string
	chFocusMainView       chan bool
	chFocusPrimitive      chan *tview.Primitive
	chDisablePageKeyEvent chan bool
}

var shared = Shared{
	isCommandLineMode:     false,
	api:                   nil,
	conf:                  nil,
	chStatus:              make(chan string, 1),
	chIndicator:           make(chan string, 1),
	chPopupModal:          make(chan *ModalOpt, 1),
	chExecCommand:         make(chan string, 1),
	chInputCommand:        make(chan string, 1),
	chFocusMainView:       make(chan bool, 1),
	chFocusPrimitive:      make(chan *tview.Primitive, 1),
	chDisablePageKeyEvent: make(chan bool, 1),
}

// SetStatus : ステータスメッセージを設定
func (s *Shared) SetStatus(label, status string) {
	message := createStatusMessage(label, status)

	if s.isCommandLineMode {
		log.Exit(message)
	}

	go func() {
		s.chStatus <- message
	}()
}

// SetErrorStatus : エラーメッセージを設定
func (s *Shared) SetErrorStatus(label, errStatus string) {
	if s.isCommandLineMode {
		log.ErrorExit(createStatusMessage(label, errStatus), log.ExitCodeErrApp)
	}

	s.SetStatus("ERR: "+label, errStatus)
}

// SetIndicator : インジケータを設定
func (s *Shared) SetIndicator(indicator string) {
	go func() {
		s.chIndicator <- indicator
	}()
}

// SetDisablePageKeyEvent : ページの共通キーハンドラを無効化
func (s *Shared) SetDisablePageKeyEvent(b bool) {
	go func() {
		s.chDisablePageKeyEvent <- b
	}()
}

// ReqestPopupModal : モーダルの表示をリクエスト
func (s *Shared) ReqestPopupModal(o *ModalOpt) {
	go func() {
		s.chPopupModal <- o
	}()
}

// RequestExecCommand : コマンドの実行をリクエスト
func (s *Shared) RequestExecCommand(c string) {
	go func() {
		s.chExecCommand <- c
	}()
}

// RequestInputCommand : コマンドの入力をリクエスト
func (s *Shared) RequestInputCommand(c string) {
	go func() {
		s.chInputCommand <- c
	}()
}

// RequestFocusPrimitive : 指定したプリミティブへのフォーカスを要求
func (s *Shared) RequestFocusPrimitive(p tview.Primitive) {
	go func() {
		s.chFocusPrimitive <- &p
	}()
}

// RequestFocusMainView : MainViewへのフォーカスを要求
func (s *Shared) RequestFocusMainView() {
	go func() {
		s.chFocusMainView <- true
	}()
}
