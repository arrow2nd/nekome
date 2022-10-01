package config

import (
	"fmt"
	"strings"

	"code.rocketnine.space/tslocum/cbind"
	"github.com/gdamore/tcell/v2"
)

const (
	// アプリ全体のアクション
	ActionQuit = "quit"

	// ページビューのアクション
	ActionSelectPrevTab = "select_prev_tab"
	ActionSelectNextTab = "select_next_tab"
	ActionClosePage     = "close_page"
	ActionRedraw        = "redraw"
	ActionFocusCmdLine  = "focus_cmdline"
	ActionShowHelp      = "show_help"

	// ページ共通のアクション
	ActionReloadPage = "reload_page"

	// ホームタイムラインページのアクション
	ActionStreamModeStart = "stream_mode_start"
	ActionStreamModeStop  = "stream_mode_stop"

	// ツイートビューのアクション
	ActionScrollUp       = "scroll_up"
	ActionScrollDown     = "scroll_down"
	ActionCursorUp       = "cursor_up"
	ActionCursorDown     = "cursor_down"
	ActionCursorTop      = "cursor_top"
	ActionCursorBottom   = "cursor_bottom"
	ActionTweetLike      = "tweet_like"
	ActionTweetUnlike    = "tweet_unlike"
	ActionTweetRetweet   = "tweet_retweet"
	ActionTweetUnretweet = "tweet_unretweet"
	ActionTweetDelete    = "tweet_delete"
	ActionUserFollow     = "user_follow"
	ActionUserUnfollow   = "user_unfollow"
	ActionUserBlock      = "user_block"
	ActionUserUnblock    = "user_unblock"
	ActionUserMute       = "user_mute"
	ActionUserUnmute     = "user_unmute"
	ActionOpenUserPage   = "open_user_page"
	ActionOpenUserLikes  = "open_user_likes"
	ActionTweet          = "tweet"
	ActionQuote          = "quote"
	ActionReply          = "reply"
	ActionOpenBrowser    = "open_browser"
	ActionCopyUrl        = "copy_url"
)

type keybinding map[string][]string

// GetString : キーバインド文字列を取得
func (k keybinding) GetString(key string) string {
	s := strings.Join(k[key], ", ")

	if s == "" {
		return "*No assignment*"
	}

	return s
}

// MappingEventHandler : キーバインドにイベントハンドラをマッピング
func (k keybinding) MappingEventHandler(handlers map[string]func()) (*cbind.Configuration, error) {
	c := cbind.NewConfiguration()

	for action, keys := range k {
		f, ok := handlers[action]
		if !ok {
			return nil, fmt.Errorf("unknown action: %s", action)
		}

		handler := func(_ *tcell.EventKey) *tcell.EventKey {
			f()
			return nil
		}

		for _, key := range keys {
			key = strings.TrimSpace(key)

			if key == "" {
				continue
			}

			if err := c.Set(key, handler); err != nil {
				return nil, fmt.Errorf("failed to set key bindings: %w", err)
			}
		}

	}

	return c, nil
}
