package app

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"html"
	"math"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/arrow2nd/nekome/v2/log"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/mattn/go-runewidth"
	"golang.org/x/crypto/ssh/terminal"
)

// openExternalEditor : 外部エディタを開く
func (a *App) openExternalEditor(editor string, args ...string) error {
	if editor == "" {
		return errors.New("please specify which editor to use")
	}

	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	var err error

	if shared.isCommandLineMode {
		err = cmd.Run()
	} else {
		a.app.Suspend(func() {
			err = cmd.Run()
		})
	}

	if err != nil {
		return fmt.Errorf("failed to open editor (%s) : %w", editor, err)
	}

	return nil
}

// getWindowWidth : 表示領域の幅を取得
func getWindowWidth() int {
	fd := int(os.Stdout.Fd())

	w, _, err := terminal.GetSize(fd)
	if err != nil {
		log.ErrorExit(err.Error(), log.ExitCodeErrTerm)
	}

	return w - 2
}

// getMD5 : MD5のハッシュ値を取得
func getMD5(s string) string {
	h := md5.New()
	defer h.Reset()

	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// getStringDisplayRow : 文字列の表示行数を取得
func getStringDisplayRow(s string, w int) int {
	row := 0

	for _, s := range strings.Split(s, "\n") {
		row += int(math.Ceil(float64(runewidth.StringWidth(s)) / float64(w)))
	}

	return row
}

// getHighlightId : ハイライト一覧からIDを取得（見つからない場合 -1 が返る）
func getHighlightId(ids []string) int {
	if ids == nil {
		return -1
	}

	i := strings.Index(ids[0], "_")
	if i == -1 || i+1 >= len(ids[0]) {
		return -1
	}

	id, err := strconv.Atoi(ids[0][i+1:])
	if err != nil {
		return -1
	}

	return id
}

// find : スライス内から任意の条件を満たす値を探す
func find[T any](s []T, f func(T) bool) (int, bool) {
	for i, v := range s {
		if f(v) {
			return i, true
		}
	}

	return -1, false
}

// truncate : 文字列を指定幅で丸める
func truncate(s string, w int) string {
	return runewidth.Truncate(s, w, "…")
}

// trimEndNewline : 末尾の改行を削除
func trimEndNewline(s string) string {
	s = strings.TrimRight(s, "\n")

	if strings.HasSuffix(s, "\r") {
		s = strings.TrimRight(s, "\r")
	}

	return s
}

// split : 文字列をスペースで分割（ダブルクオートで囲まれた部分は残す）
func split(s string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = ' '
	return r.Read()
}

// replaceAll : 正規表現にマッチした文字列を一斉置換
func replaceAll(str, reg, rep string) string {
	replace := regexp.MustCompile(reg)
	return replace.ReplaceAllString(str, rep)
}

// replaceLayoutTag : レイアウトタグを置換
func replaceLayoutTag(l, tag, newStr string) string {
	// newStrが空なら、タグと後ろにある空白文字を削除
	// NOTE: 後ろに改行がある場合に無駄な空白行ができるのを防止
	if newStr == "" {
		return replaceAll(l, tag+"\\s?", "")
	}

	return strings.ReplaceAll(l, tag, newStr)
}

// isSameDate : 同じ日付かどうか
func isSameDate(t time.Time) bool {
	now := time.Now()
	location := now.Location()
	fixedTime := t.In(location)

	t1 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	t2 := time.Date(fixedTime.Year(), fixedTime.Month(), fixedTime.Day(), 0, 0, 0, 0, location)

	return t1.Equal(t2)
}

// convertDateString : 日付文字列を変換
func convertDateString(createAt string) string {
	pref := shared.conf.Pref.Appearance

	t, _ := time.Parse(time.RFC3339, createAt)
	format := ""

	// 今日の日付なら時刻のみを表示
	if isSameDate(t) {
		format = pref.TimeFormat
	} else {
		format = fmt.Sprintf("%s %s", pref.DateFormat, pref.TimeFormat)
	}

	return t.Local().Format(format)
}

// createStyledText : スタイル適応済みの文字列を作成
func createStyledText(style, text string) string {
	return fmt.Sprintf("[%s]%s[-:-:-]", style, text)
}

// createSeparator : 指定幅のセパレータ文字列を作成
func createSeparator(s string, width int) string {
	return createStyledText(
		shared.conf.Style.Tweet.Separator,
		strings.Repeat(s, width),
	)
}

// createMetricsString : ツイートのリアクション数文字列を作成
func createMetricsString(unit, style string, count int) string {
	if count <= 0 {
		return ""
	} else if count > 1 {
		unit += "s"
	}

	return createStyledText(style, strconv.Itoa(count)+unit)
}

// createUserSummary : ユーザの要約文を作成
func createUserSummary(u *twitter.UserObj) string {
	return fmt.Sprintf("%s @%s", u.Name, u.UserName)
}

// createTweetSummary : ツイートの要約文を作成
func createTweetSummary(t *twitter.TweetDictionary) string {
	return fmt.Sprintf("%s | %s", createUserSummary(t.Author), html.UnescapeString(t.Tweet.Text))
}

// createTweetUrl : ツイートのURLを作成
func createTweetUrl(t *twitter.TweetDictionary) (string, error) {
	return url.JoinPath("https://twitter.com/", t.Author.UserName, "status", t.Tweet.ID)
}

// createStatusMessage : ラベル付きステータスメッセージを作成
func createStatusMessage(label, status string) string {
	width := getWindowWidth()
	status = strings.ReplaceAll(status, "\n", " ")

	return truncate(fmt.Sprintf("[%s] %s", label, status), width)
}
