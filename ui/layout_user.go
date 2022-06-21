package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/mattn/go-runewidth"
)

const userDescColMax = 3

func createUserDesc(d string) (string, int) {
	width := getWindowWidth() - userProfilePaddingX*2

	desc := strings.ReplaceAll(d, "\n", " ")
	desc = runewidth.Truncate(desc, width*userDescColMax, "…")

	col := int(math.Ceil(float64(runewidth.StringWidth(desc)) / float64(width)))

	return desc, col
}

func createUserInfo(u *twitter.UserObj) string {
	texts := []string{}

	if u.Location != "" {
		texts = append(texts, " "+u.Location)
	}

	if u.URL != "" {
		texts = append(texts, " "+u.URL)
	}

	return "[gray:-:-]" + strings.Join(texts, " | ")
}

func createUserProfile(u *twitter.UserObj) (string, int) {
	isInfoEmpty := u.Location == "" && u.URL == ""
	profile := createUserText(u, -1)

	desc, col := createUserDesc(u.Description)
	profile += fmt.Sprintf("[white:-:-]%s\n", desc)

	if isInfoEmpty {
		return profile, col + 1
	}

	return profile + createUserInfo(u), col + 2
}
