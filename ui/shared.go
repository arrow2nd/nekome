package ui

import (
	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
)

var shared = Shared{
	api:      nil,
	conf:     nil,
	chStatus: make(chan string, 1),
	chDetail: make(chan string, 1),
}

// Shared : 全体共有
type Shared struct {
	api      *api.API
	conf     *config.Config
	chStatus chan string
	chDetail chan string
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

// SetDetail : 詳細情報を設定
func (s *Shared) SetDetail(detail string) {
	go func() {
		shared.chDetail <- detail
	}()
}
