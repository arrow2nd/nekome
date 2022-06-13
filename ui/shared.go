package ui

import (
	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
)

var shared = Shared{
	api:       nil,
	conf:      nil,
	stateCh:   make(chan string, 1),
	appDrawCh: make(chan bool, 1),
}

// Shared 共有
type Shared struct {
	api       *api.API
	conf      *config.Config
	stateCh   chan string
	appDrawCh chan bool
}

// setStatus ステータスをセット
func (s *Shared) setStatus(state string) {
	go func() {
		shared.stateCh <- state
	}()
}

// reqestDrawApp アプリの描画処理をリクエスト
func (s *Shared) reqestDrawApp() {
	go func() {
		shared.appDrawCh <- true
	}()
}
