package config

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestColor(t *testing.T) {
	tests := []struct {
		name  string
		color color
		want  tcell.Color
	}{
		{
			name:  "正しく変換できるか",
			color: "#D162CB",
			want:  tcell.NewHexColor(0xD162CB),
		},
		{
			name:  "W3Cの色名を解釈できるか",
			color: "red",
			want:  tcell.ColorRed,
		},
		{
			name:  "空文字ならデフォルト色が返るか",
			color: "",
			want:  tcell.ColorDefault,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.ToColor(); got != tt.want {
				t.Errorf("color.ToColor() = %v, want %v", got, tt.want)
			}
		})
	}
}
