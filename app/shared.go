package app

import (
	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
	"github.com/arrow2nd/nekome/log"
	"github.com/rivo/tview"
)

// Shared : 全体共有
type Shared struct {
	api                   *api.API
	conf                  *config.Config
	isCommandLineMode     bool
	chStatus              chan string
	chIndicator           chan string
	chPopupModal          chan *ModalOpt
	chExecCommand         chan string
	chInputCommand        chan string
	chFocusView           chan bool
	chFocusPrimitive      chan *tview.Primitive
	chDisableViewKeyEvent chan bool
}

var shared = Shared{
	api:                   nil,
	conf:                  nil,
	isCommandLineMode:     false,
	chStatus:              make(chan string, 1),
	chIndicator:           make(chan string, 1),
	chPopupModal:          make(chan *ModalOpt, 1),
	chExecCommand:         make(chan string, 1),
	chInputCommand:        make(chan string, 1),
	chFocusView:           make(chan bool, 1),
	chFocusPrimitive:      make(chan *tview.Primitive, 1),
	chDisableViewKeyEvent: make(chan bool, 1),
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

// SetDisableViewKeyEvent : ビューのキーイベントを無効化
func (s *Shared) SetDisableViewKeyEvent(b bool) {
	go func() {
		s.chDisableViewKeyEvent <- b
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

// RequestFocusView : ビューへのフォーカスを要求
func (s *Shared) RequestFocusView() {
	go func() {
		s.chFocusView <- true
	}()
}
