package config

import "github.com/dghubble/oauth1"

type feature struct {
	// Consumer : コンシューマトークン
	Consumer oauth1.Token
	// MainUser : メインで使用するユーザ
	MainUser string
	// LoadTweetsCount : 1度に読み込むツイート数
	LoadTweetsCount int
	// TweetMaxAccumulationNum : ツイートの最大蓄積数
	TweetMaxAccumulationNum int
	// IsLocaleCJK : ロケールがCJKか
	IsLocaleCJK bool
	// Confirm : 確認ウィンドウの表示
	Confirm map[string]bool
	// RunCommands : 起動時に実行するコマンド
	RunCommands []string
}

type appearance struct {
	// StyleFile : 配色テーマファイルのパス
	StyleFile string
	// DateFormat : 日付のフォーマット
	DateFormat string
	// TimeFormat : 時刻のフォーマット
	TimeFormat string
	// UserBIOMaxRow : ユーザBIOの最大表示行数
	UserBIOMaxRow int
	// UserProfilePaddingX : ユーザプロフィールの左右パディング
	UserProfilePaddingX int
	// GraphChar : 投票グラフの表示に使用する文字
	GraphChar string
	// GraphMaxWidth : 投票グラフの最大表示幅
	GraphMaxWidth int
	// TabSeparate : タブの区切り文字
	TabSeparate string
	// TabMaxWidth : タブの最大表示幅
	TabMaxWidth int
}

type texts struct {
	// Like : いいねの単位
	Like string
	// Retweet : リツイートの単位
	Retweet string
	// Loading : 読み込み中
	Loading string
	// NoTweets : ツイート無し
	NoTweets string
	// TabHome : ホームタイムラインのタブ
	TabHome string
	// TabMention : メンションタイムラインのタブ
	TabMention string
	// TabList : リストタイムラインのタブ
	TabList string
	// TabUser : ユーザページのタブ
	TabUser string
	// TabSearch : 検索ページのタブ
	TabSearch string
	// TabDocs : テキストページのタブ
	TabDocs string
}

type icon struct {
	// Geo : 位置情報
	Geo string
	// Link : リンク
	Link string
	// Pinned : ピン留め
	Pinned string
	// Verified : 認証バッジ
	Verified string
	// Private : 非公開バッジ
	Private string
}

// Settings : 設定
type Settings struct {
	// Feature : 機能に関する設定
	Feature feature
	// Apperance : 外観に関する設定
	Apperance appearance
	// Texts : 各種文字列
	Texts texts
	// Icon : 各種アイコン
	Icon icon
}

func defaultSettings() *Settings {
	return &Settings{
		Feature: feature{
			Consumer:                oauth1.Token{Token: "", TokenSecret: ""},
			MainUser:                "",
			LoadTweetsCount:         25,
			TweetMaxAccumulationNum: 250,
			IsLocaleCJK:             true,
			Confirm: map[string]bool{
				"Like":      true,
				"Unlike":    true,
				"Retweet":   true,
				"Unretweet": true,
				"Delete":    true,
				"Follow":    true,
				"Unfollow":  true,
				"Block":     true,
				"Unblock":   true,
				"Mute":      true,
				"Unmute":    true,
				"Tweet":     true,
				"Quit":      true,
			},
			RunCommands: []string{
				"home",
				"mention --unfocus",
			},
		},
		Apperance: appearance{
			StyleFile:           "default.yml",
			DateFormat:          "2006/01/02",
			TimeFormat:          "15:04:05",
			UserBIOMaxRow:       3,
			UserProfilePaddingX: 4,
			GraphChar:           "\u2588",
			GraphMaxWidth:       30,
			TabSeparate:         "|",
			TabMaxWidth:         20,
		},
		Texts: texts{
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
