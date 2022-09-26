package config

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type color string

// ToColor : tcell.Colorに変換
func (c color) ToColor() tcell.Color {
	s := strings.Replace(string(c), "#", "", 1)
	i, _ := strconv.ParseInt(s, 16, 32)
	return tcell.NewHexColor(int32(i))
}

type appStyle struct {
	BackgroundColor             color  `toml:"background_color"`
	ContrastBackgroundColor     color  `toml:"contrast_background_color"`
	MoreContrastBackgroundColor color  `toml:"more_contrast_background_color"`
	BorderColor                 color  `toml:"border_color"`
	TitleColor                  color  `toml:"title_color"`
	GraphicsColor               color  `toml:"graphics_color"`
	TextColor                   color  `toml:"text_color"`
	SecondaryTextColor          color  `toml:"secondary_text_color"`
	TertiaryTextColor           color  `toml:"tertiary_text_color"`
	InverseTextColor            color  `toml:"inverse_text_color"`
	ContrastSecondaryTextColor  color  `toml:"contrast_secondary_text_color"`
	TabText                     string `toml:"tab_text"`
	Separator                   string `toml:"separator"`
}

type statusBarStyle struct {
	Text            string `toml:"text"`
	BackgroundColor color  `toml:"background_color"`
}

type tweetStyle struct {
	Annotation string `toml:"annotation"`
	Detail     string `toml:"detail"`
	Like       string `toml:"like"`
	Retweet    string `toml:"retweet"`
	HashTag    string `toml:"hashtag"`
	Mention    string `toml:"mention"`
	PollGraph  string `toml:"poll_graph"`
	PollDetail string `toml:"poll_detail"`
}

type userStyle struct {
	Name     string `toml:"name"`
	UserName string `toml:"user_name"`
	Detail   string `toml:"detaill"`
	Verified string `toml:"verified"`
	Private  string `toml:"private"`
}

type metricsStyle struct {
	TweetsText               string `toml:"tweets_text"`
	TweetsBackgroundColor    color  `toml:"tweets_background_color"`
	FollowingText            string `toml:"following_text"`
	FollowingBackgroundColor color  `toml:"following_background_color"`
	FollowersText            string `toml:"followers_text"`
	FollowersBackgroundColor color  `toml:"followers_background_color"`
}

// Style : スタイル
type Style struct {
	App       appStyle       `toml:"app"`
	StatusBar statusBarStyle `toml:"statusbar"`
	Tweet     tweetStyle     `toml:"tweet"`
	User      userStyle      `toml:"user"`
	Metrics   metricsStyle   `toml:"metrics"`
}

func defaultStyle() *Style {
	return &Style{
		App: appStyle{
			// TODO: どこに影響するのか調査する
			BackgroundColor:             "#000000",
			ContrastBackgroundColor:     "#606CB2",
			MoreContrastBackgroundColor: "#01A860",
			BorderColor:                 "#ffffff",
			TitleColor:                  "#ffffff",
			GraphicsColor:               "#ffffff",
			TextColor:                   "#ffffff",
			SecondaryTextColor:          "#FEE806",
			TertiaryTextColor:           "#01A860",
			InverseTextColor:            "#606CB2",
			ContrastSecondaryTextColor:  "#144384",
			TabText:                     "-:-:-",
			Separator:                   "gray:-:-",
		},
		StatusBar: statusBarStyle{
			Text:            "black:-:-",
			BackgroundColor: "#ffffff",
		},
		Tweet: tweetStyle{
			Annotation: "blue:-:-",
			Detail:     "gray:-:-",
			Like:       "pink:-:-",
			Retweet:    "green:-:-",
			HashTag:    "blue:-:-",
			Mention:    "blue:-:-",
			PollGraph:  "blue:-:-",
			PollDetail: "gray:-:-",
		},
		User: userStyle{
			Name:     "white:-:b",
			UserName: "gray:-:i",
			Detail:   "gray:-:-",
			Verified: "blue:-:-",
			Private:  "gray:-:-",
		},
		Metrics: metricsStyle{
			TweetsText:               "black:-:-",
			TweetsBackgroundColor:    "#a094c7",
			FollowingText:            "black:-:-",
			FollowingBackgroundColor: "#84a0c6",
			FollowersText:            "black:-:-",
			FollowersBackgroundColor: "#89b8c2",
		},
	}
}
