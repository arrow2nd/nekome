package app

import (
	"testing"
	"time"

	"github.com/arrow2nd/nekome/v2/config"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/assert"
)

func TestOpenExternalEditor(t *testing.T) {
	a := &App{}

	t.Run("エディタの起動コマンドが未指定の場合にエラーが返るか", func(t *testing.T) {
		err := a.openExternalEditor("")
		assert.EqualError(t, err, "please specify which editor to use")
	})
}

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
			got := getMD5(tt.arg)
			assert.Equal(t, tt.want, got)
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
			name: "半角文字",
			s:    "morino",
			w:    1,
			want: 6,
		},
		{
			name: "全角文字",
			s:    "杜野",
			w:    10,
			want: 1,
		},
		{
			name: "複数行に渡る文字列",
			s:    "serizawa_asahi,mayuzumi_fuyuko,izumi_mei",
			w:    10,
			want: 4,
		},
		{
			name: "改行を含む文字列",
			s:    "osaki_tenka\nosaki_amana\nkuwayama_chiyuki",
			w:    11,
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStringDisplayRow(tt.s, tt.w)
			assert.Equal(t, tt.want, got)
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
			got := getHighlightId(tt.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		name  string
		s     []int
		f     func(int) bool
		want  int
		found bool
	}{
		{
			name: "条件を満たす",
			s:    []int{1, 2, 3},
			f: func(e int) bool {
				return e == 1
			},
			want:  0,
			found: true,
		},
		{
			name: "条件を満たさない",
			s:    []int{1, 2, 3},
			f: func(e int) bool {
				return e > 4
			},
			want:  -1,
			found: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, got := find(tt.s, tt.f)
			assert.Equal(t, tt.want, index)
			assert.Equal(t, tt.found, got)
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
			got := truncate(tt.s, tt.w)
			assert.Equal(t, tt.want, got)
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
			got := trimEndNewline(tt.s)
			assert.Equal(t, tt.want, got)
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

func TestReplaceLayoutTag(t *testing.T) {
	tests := []struct {
		name string
		l    string
		t    string
		s    string
		want string
	}{
		{
			name: "置換できるか",
			l:    "I am {test} man",
			t:    "{test}",
			s:    "iron",
			want: "I am iron man",
		},
		{
			name: "置換文字列がある場合タグ末尾の空白文字が残るか",
			l:    "{test}\t",
			t:    "{test}",
			s:    "neko-chan",
			want: "neko-chan\t",
		},
		{
			name: "置換文字列が空の場合タグ末尾の空白文字が消去されるか",
			l:    "{test}\n",
			t:    "{test}",
			s:    "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := replaceLayoutTag(tt.l, tt.t, tt.s)
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
			got := isSameDate(tt.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateSeparator(t *testing.T) {
	shared = Shared{
		conf: &config.Config{
			Style: &config.Style{
				Tweet: config.TweetStyle{
					Separator: "style_sep",
				},
			},
		},
	}

	t.Run("生成できるか", func(t *testing.T) {
		s := createSeparator("-", 10)
		want := "[style_sep]----------[-:-:-]"

		assert.Equal(t, want, s)
	})
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
			want:    "[pink]1Fav[-:-:-]",
		},
		{
			name:    "2いいね",
			unit:    "Fav",
			style:   "pink",
			count:   2,
			reverse: false,
			want:    "[pink]2Favs[-:-:-]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := createMetricsString(tt.unit, tt.style, tt.count)
			assert.Equal(t, tt.want, s)
		})
	}
}

func TestCreateUserSummary(t *testing.T) {
	summary := createUserSummary(&twitter.UserObj{
		Name:     "TEST",
		UserName: "test",
	})

	t.Run("作成できるか", func(t *testing.T) {
		assert.Equal(t, "TEST @test", summary)
	})
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
			s := createTweetSummary(tt.t)
			assert.Equal(t, tt.want, s)
		})
	}
}

func TestCreateTweetURL(t *testing.T) {
	url, err := createTweetUrl(&twitter.TweetDictionary{
		Tweet: twitter.TweetObj{
			ID: "0123456789",
		},
		Author: &twitter.UserObj{
			UserName: "test",
		},
	})

	t.Run("作成できるか", func(t *testing.T) {
		want := "https://twitter.com/test/status/0123456789"

		assert.NoError(t, err)
		assert.Equal(t, want, url)
	})
}
