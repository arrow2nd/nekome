package app

import (
	"fmt"
	"strings"

	"github.com/g8rswimmer/go-twitter/v2"
)

// createUserBioLayout : レイアウト済みのBIO文字列を作成し、その表示行数を返す
func createUserBioLayout(d string, w int) (string, int) {
	desc := strings.ReplaceAll(d, "\n", " ")

	maxRow := shared.conf.Pref.Appearance.UserBIOMaxRow
	desc = truncate(desc, w*maxRow)

	return desc, getStringDisplayRow(desc, w)
}

// createUserDetailLayout : レイアウト済みのユーザ詳細文字列を作成
func createUserDetailLayout(u *twitter.UserObj) string {
	texts := []string{}

	if u.Location != "" {
		texts = append(texts, shared.conf.Pref.Icon.Geo+" "+u.Location)
	}

	if u.URL != "" {
		texts = append(texts, shared.conf.Pref.Icon.Link+" "+u.URL)
	}

	return fmt.Sprintf(
		"[%s]%s[-:-:-]",
		shared.conf.Style.User.Detail,
		strings.Join(texts, " | "),
	)
}

// createProfileLayout : レイアウト済みのプロフィール文字列を作成し、その表示行数を返す
func createProfileLayout(u *twitter.UserObj, w int) (string, int) {
	padding := shared.conf.Pref.Appearance.UserProfilePaddingX
	width := w - padding*2

	desc, row := createUserBioLayout(u.Description, width)
	desc = fmt.Sprintf("%s\n", desc)

	profile := createUserInfoLayout(u, -1, width) + desc

	// 詳細情報が無い
	if u.Location == "" && u.URL == "" {
		return profile, row + 1
	}

	return profile + createUserDetailLayout(u), row + 2
}
