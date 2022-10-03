package app

import (
	"strings"

	"github.com/g8rswimmer/go-twitter/v2"
)

// createProfileLayout : プロフィールのレイアウトを作成, 表示行数を返す
func createProfileLayout(u *twitter.UserObj, w int) (string, int) {
	padding := shared.conf.Pref.Appearance.UserProfilePaddingX
	width := w - padding*2

	layout := shared.conf.Pref.Layout.User
	layout = replaceLayoutTag(layout, "{user_info}", createUserInfoLayout(u, -1, width))
	layout = replaceLayoutTag(layout, "{bio}", createUserBioLayout(u.Description, width))
	layout = replaceLayoutTag(layout, "{user_detail}", createUserDetailLayout(u))

	return layout, getStringDisplayRow(layout, width)
}

// createUserBioLayout : BIOのレイアウトを作成
func createUserBioLayout(d string, w int) string {
	desc := strings.ReplaceAll(d, "\n", " ")
	maxRow := shared.conf.Pref.Appearance.UserBIOMaxRow

	return truncate(desc, w*maxRow)
}

// createUserDetailLayout : ユーザ詳細のレイアウトを作成
func createUserDetailLayout(u *twitter.UserObj) string {
	icon := shared.conf.Pref.Icon
	details := []string{}

	if u.Location != "" {
		details = append(details, icon.Geo+" "+u.Location)
	}

	if u.URL != "" {
		details = append(details, icon.Link+" "+u.URL)
	}

	if len(details) == 0 {
		return ""
	}

	style := shared.conf.Style.User.Detail
	separator := shared.conf.Pref.Appearance.UserDetailSeparator

	return createStyledText(style, strings.Join(details, separator))
}
