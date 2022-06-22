package ui

import (
	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
)

var shared = Shared{
	api:         nil,
	conf:        nil,
	chStatus:    make(chan string, 1),
	chIndicator: make(chan string, 1),
}

// Shared : 全体共有
type Shared struct {
	api         *api.API
	conf        *config.Config
	chStatus    chan string
	chIndicator chan string
}

// SetStatus : ステータスメッセージを設定
func (s *Shared) SetStatus(label, status string) {
	go func() {
		shared.chStatus <- createStatusMessage(label, status)
	}()
}

// SetErrorStatus : エラーメッセージを設定
func (s *Shared) SetErrorStatus(label, errStatus string) {
	s.SetStatus("ERR: "+label, errStatus)
}

// SetIndicator : インジケータを設定
func (s *Shared) SetIndicator(indicator string) {
	go func() {
		shared.chIndicator <- indicator
	}()
}
