package app

import "testing"

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
