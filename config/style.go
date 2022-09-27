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
	BackgroundColor color  `toml:"background_color"`
	BorderColor     color  `toml:"border_color"`
	TextColor       color  `toml:"text_color"`
	EmphasisText    string `toml:"emphasis_text"`
}

type tabStyle struct {
	Text            string `toml:"text"`
	BackgroundColor color  `toml:"background_color"`
}

type autocompleteStyle struct {
	BackgroundColor         color `toml:"background_color"`
	SelectedBackgroundColor color `toml:"selected_background_color"`
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
	Separator  string `toml:"separator"`
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

// Style : スタイル定義
type Style struct {
	App          appStyle          `toml:"app"`
	Tab          tabStyle          `toml:"tab"`
	Autocomplate autocompleteStyle `toml:"autocomplete"`
	StatusBar    statusBarStyle    `toml:"statusbar"`
	Tweet        tweetStyle        `toml:"tweet"`
	User         userStyle         `toml:"user"`
	Metrics      metricsStyle      `toml:"metrics"`
}

func defaultStyle() *Style {
	return &Style{
		App: appStyle{
			BackgroundColor: "#000000",
			BorderColor:     "#ffffff",
			TextColor:       "#f9f9f9",
			EmphasisText:    "maroon:-:bi",
		},
		Tab: tabStyle{
			Text:            "white:-:-",
			BackgroundColor: "#000000",
		},
		Autocomplate: autocompleteStyle{
			BackgroundColor:         "#808080",
			SelectedBackgroundColor: "#C0C0C0",
		},
		StatusBar: statusBarStyle{
			Text:            "black:-:-",
			BackgroundColor: "#ffffff",
		},
		Tweet: tweetStyle{
			Annotation: "teal:-:-",
			Detail:     "gray:-:-",
			Like:       "pink:-:-",
			Retweet:    "lime:-:-",
			HashTag:    "aqua:-:-",
			Mention:    "aqua:-:-",
			PollGraph:  "aqua:-:-",
			PollDetail: "gray:-:-",
			Separator:  "gray:-:-",
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
