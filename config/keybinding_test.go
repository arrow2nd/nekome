package config

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestKeybinding(t *testing.T) {
	t.Run("マッピングできるか", func(t *testing.T) {
		ok := false

		k := keybinding{
			ActionQuit: {"ctrl+q"},
		}

		h := map[string]func(){
			ActionQuit: func() {
				ok = true
			},
		}

		c, err := k.MappingEventHandler(h)
		assert.NoError(t, err)

		c.Capture(tcell.NewEventKey(tcell.KeyCtrlQ, 0, tcell.ModCtrl))
		assert.True(t, ok)
	})
}
