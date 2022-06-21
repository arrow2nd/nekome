package ui

import (
	"fmt"
	"strings"

	"github.com/g8rswimmer/go-twitter/v2"
)

const userDescColMax = 3

func createUserDescText(d string, w int) (string, int) {
	desc := strings.ReplaceAll(d, "\n", " ")
	desc = truncate(desc, w*userDescColMax)

	return desc, getStringDisplayColumn(desc, w)
}

func createUserInfoText(u *twitter.UserObj) string {
	texts := []string{}

	if u.Location != "" {
		texts = append(texts, "\uf276 "+u.Location)
	}

	if u.URL != "" {
		texts = append(texts, "\uf838 "+u.URL)
	}

	return "[gray:-:-]" + strings.Join(texts, " | ")
}

func createProfileLayout(u *twitter.UserObj, w int) (string, int) {
	width := w - userProfilePaddingX*2

	profile := createUserText(u, -1, width)

	desc, col := createUserDescText(u.Description, width)
	profile += fmt.Sprintf("[white:-:-]%s\n", desc)

	// 詳細情報が無い
	if u.Location == "" && u.URL == "" {
		return profile, col + 1
	}

	return profile + createUserInfoText(u), col + 2
}
