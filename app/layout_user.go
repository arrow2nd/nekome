package app

import (
	"fmt"
	"strings"

	"github.com/g8rswimmer/go-twitter/v2"
)

// createUserBioLayout : レイアウト済みのBIO文字列を作成し、その表示行数を返す
func createUserBioLayout(d string, w int) (string, int) {
	desc := strings.ReplaceAll(d, "\n", " ")

	maxRow := shared.conf.Settings.Apperance.UserBIOMaxRow
	desc = truncate(desc, w*maxRow)

	return desc, getStringDisplayRow(desc, w)
}

// createUserDetailLayout : レイアウト済みのユーザ詳細文字列を作成
func createUserDetailLayout(u *twitter.UserObj) string {
	texts := []string{}

	if u.Location != "" {
		icon := shared.conf.Settings.Icon.Geo
		texts = append(texts, icon+" "+u.Location)
	}

	if u.URL != "" {
		icon := shared.conf.Settings.Icon.Link
		texts = append(texts, icon+" "+u.URL)
	}

	return "[gray:-:-]" + strings.Join(texts, " | ")
}

// createProfileLayout : レイアウト済みのプロフィール文字列を作成し、その表示行数を返す
func createProfileLayout(u *twitter.UserObj, w int) (string, int) {
	padding := shared.conf.Settings.Apperance.UserProfilePaddingX
	width := w - padding*2

	desc, row := createUserBioLayout(u.Description, width)

	profile := createUserInfoLayout(u, -1, width) +
		fmt.Sprintf("[white:-:-]%s\n", desc)

	// 詳細情報が無い
	if u.Location == "" && u.URL == "" {
		return profile, row + 1
	}

	return profile + createUserDetailLayout(u), row + 2
}
