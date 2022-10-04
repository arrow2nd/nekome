package log_test

import (
	"testing"

	"github.com/arrow2nd/nekome/v2/log"
	"github.com/stretchr/testify/assert"
)

func TestExit(t *testing.T) {
	t.Run("終了コードが取得できるか", func(t *testing.T) {
		code := log.ExitCodeOK.GetInt()
		want := 0

		assert.Equal(t, want, code)
	})
}
