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
	chDraw:        make(chan bool, 1),
}

// Shared 共有
type Shared struct {
	api           *api.API
	conf          *config.Config
	chNormalState chan string
	chErrorState  chan string
	chDraw        chan bool
}

// setStatus ステータスメッセージを設定
func (s *Shared) setStatus(label, status string) {
	go func() {
		shared.chNormalState <- createStatusMessage(label, status)
	}()
}

// setErrorStatus エラーメッセージを設定
func (s *Shared) setErrorStatus(label, errStatus string) {
	go func() {
		shared.chErrorState <- createStatusMessage(label, errStatus)
	}()
}

// reqestDrawApp アプリの描画処理をリクエスト
func (s *Shared) reqestDrawApp() {
	go func() {
		shared.chDraw <- true
	}()
}
