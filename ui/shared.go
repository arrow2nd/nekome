package ui

import (
	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
)

var shared = Shared{
	api:     nil,
	conf:    nil,
	stateCh: make(chan string),
}

// Shared 共有
type Shared struct {
	api     *api.API
	conf    *config.Config
	stateCh chan string
}

func (s *Shared) setStatus(state string) {
	go func() {
		shared.stateCh <- state
	}()
}
