package ui

import (
	"strconv"
	"strings"
)

// getHighlightId ハイライト一覧からIDを取得
func getHighlightId(ids []string) int {
	if ids == nil {
		return -1
	}

	i := strings.Index(ids[0], "_")
	if i == -1 || i+1 >= len(ids[0]) {
		return -1
	}

	id, err := strconv.Atoi(ids[0][i+1:])
	if err != nil {
		return -1
	}

	return id
}

func setStatus(state string) {
	go func() {
		shared.stateCh <- state
	}()
}
