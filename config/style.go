package config

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type hex string

// ToColor : tcell.Colorに変換
func (h hex) ToColor() tcell.Color {
	s := strings.Replace(string(h), "#", "", 1)
	i, _ := strconv.ParseInt(s, 16, 32)
	return tcell.NewHexColor(int32(i))
}

type appStyle struct {
	Tab       string
	Separator string
}

type statusBarStyle struct {
	Text string
	BG   hex
}

type autocompleteStyle struct {
	NormalBG hex
	SelectBG hex
}

type tweetStyle struct {
	Annotation string
	Detail     string
	Like       string
	RT         string
	HashTag    string
	Mention    string
	PollGraph  string
	PollDetail string
}

type userStyle struct {
	Name                 string
	UserName             string
	Verified             string
	Private              string
	Detail               string
	TweetsMetricsText    string
	TweetsMetricsBG      hex
	FollowingMetricsText string
	FollowingMetricsBG   hex
	FollowersMetricsText string
	FollowersMetricsBG   hex
}

// Style : スタイル
type Style struct {
	App          appStyle
	StatusBar    statusBarStyle
	Autocomplete autocompleteStyle
	Tweet        tweetStyle
	User         userStyle
}

func defaultStyle() *Style {
	return &Style{
		App: appStyle{
			Tab:       "-:-:-",
			Separator: "gray:-:-",
		},
		StatusBar: statusBarStyle{
			Text: "black:-:-",
			BG:   "#ffffff",
		},
		Autocomplete: autocompleteStyle{
			NormalBG: "#3e4359",
			SelectBG: "#5c6586",
		},
		Tweet: tweetStyle{
			Annotation: "blue:-:-",
			Detail:     "gray:-:-",
			Like:       "pink:-:-",
			RT:         "green:-:-",
			HashTag:    "blue:-:-",
			Mention:    "blue:-:-",
			PollGraph:  "blue:-:-",
			PollDetail: "gray:-:-",
		},
		User: userStyle{
			Name:                 "white:-:b",
			UserName:             "gray:-:i",
			Verified:             "blue:-:-",
			Private:              "gray:-:-",
			Detail:               "gray:-:-",
			TweetsMetricsText:    "black:-:-",
			TweetsMetricsBG:      "#a094c7",
			FollowingMetricsText: "black:-:-",
			FollowingMetricsBG:   "#84a0c6",
			FollowersMetricsText: "black:-:-",
			FollowersMetricsBG:   "#89b8c2",
		},
	}
}
