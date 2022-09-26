package config

type feature struct {
	// MainUser : メインで使用するユーザ
	MainUser string `toml:"main_user"`
	// LoadTweetsLimit : 1度に読み込むツイート数
	LoadTweetsLimit int `toml:"load_tweets_limit"`
	// AccmulateTweetsLimit : ツイートの最大蓄積数
	AccmulateTweetsLimit int `toml:"accmulate_tweets_limit"`
	// UseExternalEditor : ツイート編集に外部エディタを使用するか
	UseExternalEditor bool `toml:"use_external_editor"`
	// IsLocaleCJK : ロケールがCJKか
	IsLocaleCJK bool `toml:"is_locale_cjk"`
	// StartupCmds : 起動時に実行するコマンド
	StartupCmds []string `toml:"startup_cmds"`
}

type appearance struct {
	// StyleFilePath : 配色テーマファイルのパス
	StyleFilePath string `toml:"style_file"`
	// DateFormat : 日付のフォーマット
	DateFormat string `toml:"date_fmt"`
	// TimeFormat : 時刻のフォーマット
	TimeFormat string `toml:"time_fmt"`
	// UserBIOMaxRow : ユーザBIOの最大表示行数
	UserBIOMaxRow int `toml:"user_bio_max_row"`
	// UserProfilePaddingX : ユーザプロフィールの左右パディング
	UserProfilePaddingX int `toml:"user_profile_padding_x"`
	// GraphChar : 投票グラフの表示に使用する文字
	GraphChar string `toml:"graph_char"`
	// GraphMaxWidth : 投票グラフの最大表示幅
	GraphMaxWidth int `toml:"graph_max_width"`
	// TabSeparate : タブの区切り文字
	TabSeparate string `toml:"tab_separate"`
	// TabMaxWidth : タブの最大表示幅
	TabMaxWidth int `toml:"tab_max_width"`
}

type text struct {
	// Like : いいねの単位
	Like string `toml:"like"`
	// Retweet : リツイートの単位
	Retweet string `toml:"retweet"`
	// Loading : 読み込み中
	Loading string `toml:"loading"`
	// NoTweets : ツイート無し
	NoTweets string `toml:"no_tweets"`
	// TabHome : ホームタブ
	TabHome string `toml:"tab_home"`
	// TabMention : メンションタブ
	TabMention string `toml:"tab_mention"`
	// TabList : リストタブ
	TabList string `toml:"tab_list"`
	// TabUser : ユーザタブ
	TabUser string `toml:"tab_user"`
	// TabSearch : 検索タブ
	TabSearch string `toml:"tab_search"`
	// TabDocs : ドキュメントタブ
	TabDocs string `toml:"tab_docs"`
}

type icon struct {
	// Geo : 位置情報
	Geo string `toml:"geo"`
	// Link : リンク
	Link string `toml:"link"`
	// Pinned : ピン留め
	Pinned string `toml:"pinned"`
	// Verified : 認証バッジ
	Verified string `toml:"verified"`
	// Private : 非公開バッジ
	Private string `toml:"private"`
}

// Settings : 環境設定
type Settings struct {
	Feature    feature         `toml:"feature"`
	Confirm    map[string]bool `toml:"comfirm"`
	Appearance appearance      `toml:"appearance"`
	Text       text            `toml:"text"`
	Icon       icon            `toml:"icon"`
}

func defaultSettings() *Settings {
	return &Settings{
		Feature: feature{
			MainUser:             "",
			LoadTweetsLimit:      25,
			AccmulateTweetsLimit: 250,
			UseExternalEditor:    false,
			IsLocaleCJK:          true,
			StartupCmds: []string{
				"home",
				"mention --unfocus",
			},
		},
		Confirm: map[string]bool{
			"like":      true,
			"unlike":    true,
			"retweet":   true,
			"unretweet": true,
			"delete":    true,
			"follow":    true,
			"unfollow":  true,
			"block":     true,
			"unblock":   true,
			"mute":      true,
			"unmute":    true,
			"tweet":     true,
			"quit":      true,
		},
		Appearance: appearance{
			StyleFilePath:       "default.toml",
			DateFormat:          "2006/01/02",
			TimeFormat:          "15:04:05",
			UserBIOMaxRow:       3,
			UserProfilePaddingX: 4,
			GraphChar:           "\u2588",
			GraphMaxWidth:       30,
			TabSeparate:         "|",
			TabMaxWidth:         20,
		},
		Text: text{
			Like:       "Like",
			Retweet:    "RT",
			Loading:    "Loading...",
			NoTweets:   "No tweets ฅ^-ω-^ฅ",
			TabHome:    "Home",
			TabMention: "Mention",
			TabList:    "List: {name}",
			TabUser:    "User: @{name}",
			TabSearch:  "Search: {query}",
			TabDocs:    "Docs: {name}",
		},
		Icon: icon{
			Geo:      "📍",
			Link:     "🔗",
			Pinned:   "📌",
			Verified: "✅",
			Private:  "🔒",
		},
	}
}
