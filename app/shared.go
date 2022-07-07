package app

import (
	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
	"github.com/arrow2nd/nekome/log"
)

var shared = Shared{
	isCommandLineMode: false,
	api:               nil,
	conf:              nil,
	chStatus:          make(chan string, 1),
	chIndicator:       make(chan string, 1),
	chPopupModal:      make(chan *ModalOpt, 1),
	chExecCommand:     make(chan string, 1),
	chInputCommand:    make(chan string, 1),
	chFocusPageView:   make(chan bool, 1),
}

// Shared : 全体共有
type Shared struct {
	isCommandLineMode bool
	api               *api.API
	conf              *config.Config
	chStatus          chan string
	chIndicator       chan string
	chPopupModal      chan *ModalOpt
	chExecCommand     chan string
	chInputCommand    chan string
	chFocusPageView   chan bool
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

// RequestFocusPageView : PageViewへのフォーカスをリクエスト
func (s *Shared) RequestFocusPageView() {
	go func() {
		s.chFocusPageView <- true
	}()
}
