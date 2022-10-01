package config

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	t.Run("設定したキーバインド文字列を取得できるか", func(t *testing.T) {
		k := keybinding{
			"test": []string{"ctrl+t"},
		}

		s := k.GetString("test")
		assert.Equal(t, "ctrl+t", s)
	})

	t.Run("複数のキーバインドを連結して文字列にできるか", func(t *testing.T) {
		k := keybinding{
			"test": []string{"ctrl+t", "t"},
		}

		s := k.GetString("test")
		assert.Equal(t, "ctrl+t, t", s)
	})

	t.Run("割り当てがない場合に特定の文字列が返るか", func(t *testing.T) {
		k := keybinding{
			"test": []string{},
		}

		s := k.GetString("test")
		assert.Equal(t, "*No assignment*", s)
	})
}

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
