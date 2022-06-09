package ui

import (
	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
)

// Shared 共有
type Shared struct {
	api    *api.API
	conf   *config.Config
	status *status
}

var shared = Shared{
	api:  nil,
	conf: nil,
	status: &status{
		text: "",
	},
}
