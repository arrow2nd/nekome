package app

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/arrow2nd/nekome/config"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/mattn/go-runewidth"
	"golang.org/x/crypto/ssh/terminal"
)

// getWindowWidth : 表示領域の幅を取得
func getWindowWidth() int {
	fd := int(os.Stdout.Fd())

	w, _, err := terminal.GetSize(fd)
	if err != nil {
		log.Fatal(err)
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

// getStringDisplayRow 文字列の表示行数を取得
func getStringDisplayRow(s string, w int) int {
	return int(math.Ceil(float64(runewidth.StringWidth(s)) / float64(w)))
}

// getHighlightId : ハイライト一覧からIDを取得
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

// truncate : 文字列を指定幅で丸める
func truncate(s string, w int) string {
	return runewidth.Truncate(s, w, "…")
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
	t, _ := time.Parse(time.RFC3339, createAt)
	format := ""

	// 今日の日付なら時刻のみを表示
	if isSameDate(t) {
		format = shared.conf.Settings.Apperance.TimeFormat
	} else {
		format = fmt.Sprintf(
			"%s %s",
			shared.conf.Settings.Apperance.DateFormat,
			shared.conf.Settings.Apperance.TimeFormat,
		)
	}

	return t.Local().Format(format)
}

// createSeparator : 指定幅のセパレータ文字列を作成
func createSeparator(s string, width int) string {
	return fmt.Sprintf("[gray:-:-]%s[-:-:-]", strings.Repeat(s, width))
}

// createMetricsString : ツイートのリアクション数文字列を作成
func createMetricsString(unit, color string, count int, reverse bool) string {
	if count <= 0 {
		return ""
	} else if count > 1 {
		unit += "s"
	}

	if reverse {
		return fmt.Sprintf("[%s:-:r] %d%s [-:-:-]", color, count, unit)
	}

	return fmt.Sprintf("[%s]%d%s[-:-:-] ", color, count, unit)
}

// createUserSummary : ユーザの要約文を作成
func createUserSummary(u *twitter.UserObj) string {
	return fmt.Sprintf("%s @%s", u.Name, u.UserName)
}

// createTweetSummary : ツイートの要約文を作成
func createTweetSummary(t *twitter.TweetDictionary) string {
	return fmt.Sprintf("%s | %s", createUserSummary(t.Author), html.EscapeString(t.Tweet.Text))
}

// createTweetURL : ツイートのURLを作成
func createTweetURL(t *twitter.TweetDictionary) string {
	return path.Join("https://twitter.com", t.Author.UserName, "status", t.Tweet.ID)
}

// createStatusMessage : ラベル付きステータスメッセージを作成
func createStatusMessage(label, status string) string {
	width := getWindowWidth()
	return truncate(fmt.Sprintf("[%s] %s", label, status), width)
}

// editTextInEditor テキストをエディタで編集
func editTextInEditor(editor string) (string, error) {
	dirPath, err := config.GetConfigDir()
	if err != nil {
		return "", err
	}

	tmpFilePath := path.Join(dirPath, ".tmp")

	cmd := exec.Command(editor, tmpFilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	// エディタを起動
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to start editor (%s) : %w", editor, err)
	}

	// 一時ファイル読み込み
	bytes, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		return "", err
	}

	// 一時ファイル削除
	if err := os.Remove(tmpFilePath); err != nil {
		return "", err
	}

	return string(bytes), nil
}
