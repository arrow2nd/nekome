package app

import (
	"fmt"
	"strings"

	"github.com/g8rswimmer/go-twitter/v2"
)

// createProfileLayout : プロフィールのレイアウトを作成, 表示行数を返す
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

// createUserBioLayout : BIOのレイアウトを作成, 表示行数を返す
func createUserBioLayout(d string, w int) (string, int) {
	desc := strings.ReplaceAll(d, "\n", " ")

	maxRow := shared.conf.Pref.Appearance.UserBIOMaxRow
	desc = truncate(desc, w*maxRow)

	fmt.Println(w * maxRow)

	return desc, getStringDisplayRow(desc, w)
}

// createUserDetailLayout : ユーザ詳細のレイアウトを作成
func createUserDetailLayout(u *twitter.UserObj) string {
	icon := shared.conf.Pref.Icon
	texts := []string{}

	if u.Location != "" {
		texts = append(texts, icon.Geo+" "+u.Location)
	}

	if u.URL != "" {
		texts = append(texts, icon.Link+" "+u.URL)
	}

	return fmt.Sprintf(
		"[%s]%s[-:-:-]",
		shared.conf.Style.User.Detail,
		strings.Join(texts, " | "),
	)
}
