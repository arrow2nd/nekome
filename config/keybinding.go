package config

import (
	"fmt"

	"code.rocketnine.space/tslocum/cbind"
	"github.com/gdamore/tcell/v2"
)

const (
	// アプリ全体のアクション
	ActionQuit = "quit"

	// ページビューのアクション
	ActionSelectPrevTab = "select_prev_tab"
	ActionSelectNextTab = "select_next_tab"
	ActionRedraw        = "redraw"
	ActionFocusCmdLine  = "focus_cmdline"
	ActionShowHelp      = "show_help"
	ActionRemovePage    = "remove_page"

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
	ActionTweetRemove    = "tweet_remove"
	ActionUserFollow     = "user_follow"
	ActionUserUnfollow   = "user_unfollow"
	ActionUserBlock      = "user_block"
	ActionUserUnblock    = "user_unblock"
	ActionUserMute       = "user_mute"
	ActionUserUnmute     = "user_unmute"
	ActionOpenUserPage   = "open_user_page"
	ActionQuote          = "quote"
	ActionReply          = "reply"
	ActionOpenBrowser    = "open_browser"
	ActionCopyUrl        = "copy_url"
)

type keybinding map[string][]string

// MappingEventHandler : キーバインドにイベントハンドラをマッピング
func (k *keybinding) MappingEventHandler(handlers map[string]func()) (*cbind.Configuration, error) {
	c := cbind.NewConfiguration()

	for action, keys := range *k {
		handler, ok := handlers[action]
		if !ok {
			return nil, fmt.Errorf("unknown action: %s", action)
		}

		for _, key := range keys {
			if err := c.Set(key, wrapEventHandler(handler)); err != nil {
				return nil, fmt.Errorf("failed to set key bindings: %w", err)
			}
		}

	}

	return c, nil
}

// wrapEventHandler : 戻り値無しの関数を設定するためのラップ関数
func wrapEventHandler(f func()) func(*tcell.EventKey) *tcell.EventKey {
	return func(_ *tcell.EventKey) *tcell.EventKey {
		f()
		return nil
	}
}
