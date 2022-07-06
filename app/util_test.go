package app

import (
	"testing"
	"time"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/assert"
)

func TestMD5(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "ハッシュ値が計算できる(A)",
			arg:  "HotaHota",
			want: "40f593f4a80645bb25e82c3fcb06c304",
		},
		{
			name: "ハッシュ値が計算できる(B)",
			arg:  "ほたほた",
			want: "b220b84063e671b9a7a4f622c6a6364f",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMD5(tt.arg); got != tt.want {
				t.Errorf("getMD5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStringDisplayRow(t *testing.T) {
	tests := []struct {
		name string
		s    string
		w    int
		want int
	}{
		{
			name: "半角文字の表示行数が取得できるか",
			s:    "morino",
			w:    1,
			want: 6,
		},
		{
			name: "全角文字の表示行数が取得できるか",
			s:    "杜野",
			w:    10,
			want: 1,
		},
		{
			name: "複数行に渡る文字列の表示行数が取得できるか",
			s:    "serizawa_asahi,mayuzumi_fuyuko,izumi_mei",
			w:    10,
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStringDisplayRow(tt.s, tt.w); got != tt.want {
				t.Errorf("getStringDisplayRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHighlightId(t *testing.T) {
	tests := []struct {
		name string
		arg  []string
		want int
	}{
		{
			name: "1桁のIDを抽出",
			arg:  []string{"page_0"},
			want: 0,
		},
		{
			name: "2桁のIDを抽出",
			arg:  []string{"tweet_10"},
			want: 10,
		},
		{
			name: "3桁のIDを抽出",
			arg:  []string{"tweet_100"},
			want: 100,
		},
		{
			name: "nilが渡された",
			arg:  nil,
			want: -1,
		},
		{
			name: "IDがない",
			arg:  []string{"page_"},
			want: -1,
		},
		{
			name: "_がない",
			arg:  []string{"rinze"},
			want: -1,
		},
		{
			name: "解析できない形式",
			arg:  []string{"asahi_serizawa"},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHighlightId(tt.arg); got != tt.want {
				t.Errorf("getHighlightId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		f    func(int) bool
		want bool
	}{
		{
			name: "条件を満たす",
			s:    []int{1, 2, 3},
			f: func(e int) bool {
				return e == 1
			},
			want: true,
		},
		{
			name: "条件を満たさない",
			s:    []int{1, 2, 3},
			f: func(e int) bool {
				return e > 4
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := find(tt.s, tt.f); got != tt.want {
				t.Errorf("find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	runewidth.DefaultCondition.EastAsianWidth = false

	tests := []struct {
		name string
		s    string
		w    int
		want string
	}{
		{
			name: "そのままの文字列が返る",
			s:    "komiya_kaho",
			w:    20,
			want: "komiya_kaho",
		},
		{
			name: "丸められた文字列が返る",
			s:    "shirase_sakuyasan",
			w:    15,
			want: "shirase_sakuya…",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncate(tt.s, tt.w); got != tt.want {
				t.Errorf("truncate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrimEndNewline(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "改行がない",
			s:    "komiya_kaho",
			want: "komiya_kaho",
		},
		{
			name: "改行を削除(CR)",
			s:    "tsukioka_kogane\r",
			want: "tsukioka_kogane",
		},
		{
			name: "改行を削除(LF)",
			s:    "mitsumine_yuika\n",
			want: "mitsumine_yuika",
		},
		{
			name: "改行を削除(CRLF)",
			s:    "tanaka_mamimi\r\n",
			want: "tanaka_mamimi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimEndNewline(tt.s); got != tt.want {
				t.Errorf("trimEndNewline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []string
	}{
		{
			name: "正しく分割できるか（英語）",
			s:    "komiya kaho",
			want: []string{"komiya", "kaho"},
		},
		{
			name: "正しく分割できるか（日本語）",
			s:    "小宮 果穂",
			want: []string{"小宮", "果穂"},
		},
		{
			name: "ダブルクオートで囲んだ部分が残るか",
			s:    `aketa mikoto "ikaruga luca" nanakusa nichika`,
			want: []string{"aketa", "mikoto", `ikaruga luca`, "nanakusa", "nichika"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := split(tt.s)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsSameDate(t *testing.T) {
	tests := []struct {
		name string
		arg  time.Time
		want bool
	}{
		{
			name: "現在の日時",
			arg:  time.Now(),
			want: true,
		},
		{
			name: "今日の日付",
			arg:  time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
			want: true,
		},
		{
			name: "過去の日付",
			arg:  time.Date(2018, 4, 24, 0, 0, 0, 0, time.Local),
			want: false,
		},
		{
			name: "未来の日付",
			arg:  time.Now().Add(time.Hour * 24),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSameDate(tt.arg); got != tt.want {
				t.Errorf("isSameDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateMetricsString(t *testing.T) {
	tests := []struct {
		name    string
		unit    string
		style   string
		count   int
		reverse bool
		want    string
	}{
		{
			name:    "いいね無し",
			unit:    "Fav",
			style:   "pink",
			count:   0,
			reverse: false,
			want:    "",
		},
		{
			name:    "1いいね",
			unit:    "Fav",
			style:   "pink",
			count:   1,
			reverse: false,
			want:    "[pink]1Fav[-:-:-] ",
		},
		{
			name:    "2いいね",
			unit:    "Fav",
			style:   "pink",
			count:   2,
			reverse: false,
			want:    "[pink]2Favs[-:-:-] ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createMetricsString(tt.unit, tt.style, tt.count, tt.reverse); got != tt.want {
				t.Errorf("createMetricsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateUserSummary(t *testing.T) {
	summary := createUserSummary(&twitter.UserObj{
		Name:     "TEST",
		UserName: "test",
	})

	assert.Equal(t, "TEST @test", summary, "作成できるか")
}

func TestCreateTweetSummary(t *testing.T) {
	tests := []struct {
		name string
		t    *twitter.TweetDictionary
		want string
	}{
		{
			name: "作成できるか",
			t: &twitter.TweetDictionary{
				Tweet: twitter.TweetObj{
					Text: "morikubo_nono",
				},
				Author: &twitter.UserObj{
					Name:     "TEST",
					UserName: "test",
				},
			},
			want: "TEST @test | morikubo_nono",
		},
		{
			name: "HTML文字列がアンエスケープされるか",
			t: &twitter.TweetDictionary{
				Tweet: twitter.TweetObj{
					Text: "&lt;&gt;&amp;&quot;&#x27;&#x60;",
				},
				Author: &twitter.UserObj{
					Name:     "TEST",
					UserName: "test",
				},
			},
			want: "TEST @test | <>&\"'`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createTweetSummary(tt.t); got != tt.want {
				t.Errorf("createTweetSummary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTweetURL(t *testing.T) {
	url := createTweetURL(&twitter.TweetDictionary{
		Tweet: twitter.TweetObj{
			ID: "0123456789",
		},
		Author: &twitter.UserObj{
			UserName: "test",
		},
	})

	assert.Contains(t, url, "/test/status/0123456789", "URLが作成できるか")
}
