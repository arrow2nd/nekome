package config

type appStyle struct {
	Tab       string
	Separator string
}

type statusbarStyle struct {
	Text string
	BG   int32
}

type autocompleteStyle struct {
	NormalBG int32
	SelectBG int32
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
	TweetsMetricsBG      int32
	FollowingMetricsText string
	FollowingMetricsBG   int32
	FollowersMetricsText string
	FollowersMetricsBG   int32
}

// Style : スタイル
type Style struct {
	App          appStyle
	Statusbar    statusbarStyle
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
		Statusbar: statusbarStyle{
			Text: "black:-:-",
			BG:   0xffffff,
		},
		Autocomplete: autocompleteStyle{
			NormalBG: 0x3e4359,
			SelectBG: 0x5c6586,
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
			Name:                 "lightgray:-:b",
			UserName:             "gray:-:i",
			Verified:             "blue:-:-",
			Private:              "gray:-:-",
			Detail:               "gray:-:-",
			TweetsMetricsText:    "black:-:-",
			TweetsMetricsBG:      0xa094c7,
			FollowingMetricsText: "black:-:-",
			FollowingMetricsBG:   0x84a0c6,
			FollowersMetricsText: "black:-:-",
			FollowersMetricsBG:   0x89b8c2,
		},
	}
}
