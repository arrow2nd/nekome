package app

import (
	"fmt"
	"sync"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tweets struct {
	view     *tview.TextView
	pinned   *twitter.TweetDictionary
	contents []*twitter.TweetDictionary
	mu       sync.Mutex
}

func newTweets() *tweets {
	t := &tweets{
		view:     tview.NewTextView(),
		pinned:   nil,
		contents: []*twitter.TweetDictionary{},
	}

	t.view.
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true).
		SetRegions(true).
		SetBackgroundColor(tcell.ColorDefault)

	t.view.
		SetHighlightedFunc(func(added, removed, remaining []string) {
			t.view.ScrollToHighlight()
		}).
		SetInputCapture(t.handleKeyEvents)

	return t
}

// GetSinceID : 一番新しいツイートIDを取得
func (t *tweets) GetSinceID() string {
	if len(t.contents) == 0 {
		return ""
	}

	return t.contents[0].Tweet.ID
}

// GetTweetsCount : 表示中のツイート数を取得
func (t *tweets) GetTweetsCount() int {
	c := len(t.contents)

	if t.pinned != nil {
		c++
	}

	return c
}

// getSelectTweet : 選択中のツイートを取得
func (t *tweets) getSelectTweet() *twitter.TweetDictionary {
	id := getHighlightId(t.view.GetHighlights())
	if id == -1 {
		return nil
	}

	c := &twitter.TweetDictionary{}

	if t.pinned == nil {
		// ピン留めツイートが無い
		c = t.contents[id]
	} else if id == 0 {
		// ピン留めツイートが選択されている
		c = t.pinned
	} else {
		// ピン留めツイート意外が選択されている
		c = t.contents[id-1]
	}

	// リツイートなら参照先に置き換え
	for _, rc := range c.ReferencedTweets {
		if rc.Reference.Type == "retweeted" {
			c = c.ReferencedTweets[0].TweetDictionary
		}
	}

	return c
}

// Register : ツイートを登録
func (t *tweets) Register(tweets []*twitter.TweetDictionary) {
	t.mu.Lock()
	defer t.mu.Unlock()

	size := len(t.contents)
	allSize := size + len(tweets)
	maxSize := shared.conf.Settings.Feature.TweetMaxAccumulationNum

	// 最大蓄積数を超えていたら古いものから削除
	if allSize > maxSize {
		size -= allSize - maxSize
	}

	t.contents = append(tweets, t.contents[:size]...)
}

// RegisterPinned : ピン留めツイートを登録
func (t *tweets) RegisterPinned(tweet *twitter.TweetDictionary) {
	t.pinned = tweet
}

// Draw : 描画（表示幅はターミナルのウィンドウ幅に依存）
func (t *tweets) Draw() {
	width := getWindowWidth()

	// ビューの初期化
	t.view.
		SetTextAlign(tview.AlignLeft).
		Clear()

	// 表示するツイートが無いなら描画を中断
	if t.GetTweetsCount() == 0 {
		t.DrawMessage(shared.conf.Settings.Texts.NoTweets)
		return
	}

	contents := t.contents

	// ピン留めツイートがある場合、先頭に追加
	if t.pinned != nil {
		contents = append([]*twitter.TweetDictionary{t.pinned}, t.contents...)
	}

	for i, content := range contents {
		var quotedTweet *twitter.TweetDictionary = nil

		// 参照ツイートを確認
		for _, rc := range content.ReferencedTweets {
			switch rc.Reference.Type {
			case "retweeted":
				fmt.Fprintln(t.view, createAnnotation("RT by", content.Author))
				content = content.ReferencedTweets[0].TweetDictionary
			case "replied_to":
				fmt.Fprintln(t.view, createAnnotation("Reply to", rc.TweetDictionary.Author))
			case "quoted":
				quotedTweet = rc.TweetDictionary
			}
		}

		// ピン留めツイート
		if i == 0 && t.pinned != nil {
			icon := shared.conf.Settings.Icon.Pinned
			fmt.Fprintf(t.view, "[gray:-:-]%s Pinned Tweet[-:-:-]\n", icon)
		}

		// 表示部分を作成
		layout := createTweetLayout(content, i, width)
		fmt.Fprintln(t.view, layout)

		// 引用元ツイートを表示
		if quotedTweet != nil {
			fmt.Fprintln(t.view, createSeparator("-", width))
			fmt.Fprintln(t.view, createTweetLayout(quotedTweet, -1, width))
		}

		// 末尾のツイートでないならセパレータを挿入
		lastIndex := t.GetTweetsCount() - 1
		if i < lastIndex {
			fmt.Fprintln(t.view, createSeparator("─", width))
		}
	}

	t.scrollToTweet(0)
}

// DrawMessage : ビューにメッセージを表示
func (t *tweets) DrawMessage(s string) {
	t.view.Clear().
		SetTextAlign(tview.AlignCenter).
		SetText(s)
}

// scrollToTweet : 指定ツイートまでスクロール
func (t *tweets) scrollToTweet(i int) {
	t.view.Highlight(createTweetTag(i))
}

// cursorUp : カーソルを上に移動
func (t *tweets) cursorUp() {
	idx := getHighlightId(t.view.GetHighlights())
	if idx == -1 {
		return
	}

	if idx--; idx < 0 {
		idx = t.GetTweetsCount() - 1
	}

	t.scrollToTweet(idx)
}

// cursorDown : カーソルを下に移動
func (t *tweets) cursorDown() {
	idx := getHighlightId(t.view.GetHighlights())
	if idx == -1 {
		return
	}

	idx = (idx + 1) % t.GetTweetsCount()

	t.scrollToTweet(idx)
}

// handleKeyEvents : ツイートビューのキーイベントハンドラ
func (t *tweets) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	keyRune := event.Rune()

	// 上にスクロール
	if key == tcell.KeyPgUp || key == tcell.KeyCtrlJ {
		r, c := t.view.GetScrollOffset()
		t.view.ScrollTo(r+1, c)
		return nil
	}

	// 下にスクロール
	if key == tcell.KeyPgDn || key == tcell.KeyCtrlK {
		r, c := t.view.GetScrollOffset()
		t.view.ScrollTo(r-1, c)
		return nil
	}

	// カーソルを上に移動
	if key == tcell.KeyUp || keyRune == 'k' {
		t.cursorUp()
		return nil
	}

	// カーソルを下に移動
	if key == tcell.KeyDown || keyRune == 'j' {
		t.cursorDown()
		return nil
	}

	// カーソルを先頭に移動
	if key == tcell.KeyHome || keyRune == 'g' {
		t.scrollToTweet(0)
		return nil
	}

	// カーソルを末尾に移動
	if key == tcell.KeyEnd || keyRune == 'G' {
		lastIndex := t.GetTweetsCount() - 1
		t.scrollToTweet(lastIndex)
		return nil
	}

	// いいね
	if keyRune == 'f' {
		t.actionForTweet(tweetLike)
		return nil
	}

	// いいね解除
	if keyRune == 'F' {
		t.actionForTweet(tweetUnlike)
		return nil
	}

	// リツイート
	if keyRune == 't' {
		t.actionForTweet(tweetRetweet)
		return nil
	}

	// リツイート解除
	if keyRune == 'T' {
		t.actionForTweet(tweetUnretweet)
		return nil
	}

	// ツイートを削除
	if keyRune == 'D' {
		t.actionForTweet(tweetDelete)
		return nil
	}

	// フォロー
	if keyRune == 'w' {
		t.actionForUser(userFollow)
		return nil
	}

	// フォロー解除
	if keyRune == 'W' {
		t.actionForUser(userUnfollow)
		return nil
	}

	// ブロック
	if keyRune == 'x' {
		t.actionForUser(userBlock)
		return nil
	}

	// ブロック解除
	if keyRune == 'X' {
		t.actionForUser(userUnblock)
		return nil
	}

	// ミュート
	if keyRune == 'u' {
		t.actionForUser(userMute)
		return nil
	}

	// ミュート解除
	if keyRune == 'U' {
		t.actionForUser(userUnmute)
		return nil
	}

	// ユーザページを開く
	if keyRune == 'i' {
		t.openUserPage()
		return nil
	}

	// 引用ツイート
	if keyRune == 'q' {
		t.postQuoteTweet()
		return nil
	}

	// リプライ
	if keyRune == 'r' {
		t.postReply()
		return nil
	}

	// ブラウザで開く
	if keyRune == 'o' {
		t.openBrower()
		return nil
	}

	// リンクをコピー
	if keyRune == 'c' {
		t.copyLinkToClipBoard()
		return nil
	}

	return event
}
