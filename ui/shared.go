package ui

import (
	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
)

var shared = Shared{
	api:           nil,
	conf:          nil,
	chNormalState: make(chan string, 1),
	chErrorState:  make(chan string, 1),
}

// Shared : 全体共有
type Shared struct {
	api           *api.API
	conf          *config.Config
	chNormalState chan string
	chErrorState  chan string
}

// SetStatus : ステータスメッセージを設定
func (s *Shared) SetStatus(label, status string) {
	go func() {
		shared.chNormalState <- createStatusMessage(label, status)
	}()
}

// SetErrorStatus : エラーメッセージを設定
func (s *Shared) SetErrorStatus(label, errStatus string) {
	go func() {
		shared.chErrorState <- createStatusMessage(label, errStatus)
	}()
}
