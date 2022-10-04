# 設定ファイルについて

> [English](../en/config.md)

以下の様な形で作成・保存されます

```
$HOME/.config/nekome
├── .cred.toml
├── style_default.toml
└── preferences.toml
```

## .cred.toml

認証情報を記録したファイルです

```toml
# Twitter API のコンシューマキー
[consumer]
  Token = ""
  TokenSecret = ""

# 注意:
# 以下は手動で編集しないでください
# ユーザ情報の操作には `nekome account` を使用してください

# ユーザの認証情報
[user]

  [[user.accounts]]
    UserName = "user_name"
    ID = "0123456789"
    [user.accounts.Token]
      Token = "hoge"
      TokenSecret = "fuga"
```

## preferences.toml

環境設定のファイルです

```toml
[feature]
  # メインで使用するユーザ
  main_user = "user_name"
  # 1度に読み込むツイート数
  load_tweets_limit = 25
  # 1ページで蓄積する最大ツイート数
  accmulate_tweets_limit = 250
  # ツイート編集に外部エディタを使用するか
  use_external_editor = false
  # 実行環境のロケールが CJK かどうか
  is_locale_cjk = true
  # TUI 起動時に実行するコマンド
  startup_cmds = ["home", "mention --unfocus"]

# 確認モーダルを表示するか
[comfirm]
  block = true
  delete = true
  follow = true
  like = true
  mute = true
  quit = true
  retweet = true
  tweet = true
  unblock = true
  unfollow = true
  unlike = true
  unmute = true
  unretweet = true

[appearance]
  # 読み込むスタイル定義ファイル
  style_file = "style_default.toml"
  # 日付のフォーマット
  date_fmt = "2006/01/02"
  # 時刻のフォーマット
  time_fmt = "15:04:05"
  # ユーザ / BIO の最大行数
  user_bio_max_row = 3
  # ユーザ / プロフィール表示域の左右パディング
  user_profile_padding_x = 4
  # ユーザ / 詳細情報のセパレータ
  user_detail_separator = " | "
  # ツイート / ツイート間のセパレータを非表示
  hide_tweet_seperator = false
  # ツイート / 引用ツイートのセパレータを非表示
  hide_quote_tweet_separator = false
  # グラフ / 表示に使用する文字
  graph_char = "█"
  # グラフ / 最大表示幅
  graph_max_width = 30
  # タブ / セパレータ
  tab_separator = "|"
  # タブ / 最大表示幅
  tab_max_width = 20

[layout]
  # ツイート表示
  tweet = "{annotation}\n{user_info}\n{text}\n{poll}\n{detail}"
  # アノテーション
  tweet_anotation = "{text} {author_name} {author_username}"
  # ツイート詳細
  tweet_detail = "{created_at} | via {via}\n{metrics}"
  # 投票表示
  tweet_poll = "{graph}\n{detail}"
  # 投票グラフ
  tweet_poll_graph = "{label}\n{graph} {per} {votes}"
  # 投票詳細
  tweet_poll_detail = "{status} | {all_votes} votes | ends on {end_date}"
  # ユーザプロフィール表示
  user = "{user_info}\n{bio}\n{user_detail}"
  # ユーザ詳細
  user_info = "{name} {username} {badge}"

[text]
  # いいねの単位
  like = "Like"
  # リツイートの単位
  retweet = "RT"
  # 読み込み中表示
  loading = "Loading..."
  # ツイート無し表示
  no_tweets = "No tweets ฅ^-ω-^ฅ"
  # タブ表示
  tab_home = "Home"
  tab_mention = "Mention"
  tab_list = "List: {name}"
  tab_user = "User: @{name}"
  tab_search = "Search: {query}"
  tab_likes = "Likes: @{name}"
  tab_docs = "Docs: {name}"

[icon]
  # 位置情報
  geo = "📍"
  # リンク
  link = "🔗"
  # ピン留めツイート
  pinned = "📌"
  # 認証済みバッジ
  verified = "✅"
  # 非公開バッジ
  private = "🔒"

[keybinding]
  # アプリ全体
  [keybinding.global]
    quit = ["ctrl+q"]
  # メインビュー
  [keybinding.view]
    close_page = ["ctrl+w"]
    focus_cmdline = [":"]
    redraw = ["ctrl+l"]
    select_next_tab = ["l", "Right"]
    select_prev_tab = ["h", "Left"]
    show_help = ["?"]
  # 全ページ共通
  [keybinding.page]
    reload_page = ["."]
  # ホームタイムラインページ
  [keybinding.home_timeline]
    stream_mode_start = ["s"]
    stream_mode_stop = ["S"]
  # ツイートビュー
  [keybinding.tweet]
    copy_url = ["c"]
    cursor_bottom = ["G", "End"]
    cursor_down = ["j", "Down"]
    cursor_top = ["g", "Home"]
    cursor_up = ["k", "Up"]
    open_browser = ["o"]
    open_user_likes = ["I"]
    open_user_page = ["i"]
    quote = ["q"]
    reply = ["r"]
    scroll_down = ["ctrl+k", "PageDown"]
    scroll_up = ["ctrl+j", "PageUp"]
    tweet = ["n"]
    tweet_delete = ["D"]
    tweet_like = ["f"]
    tweet_retweet = ["t"]
    tweet_unlike = ["F"]
    tweet_unretweet = ["T"]
    user_block = ["x"]
    user_follow = ["w"]
    user_mute = ["u"]
    user_unblock = ["X"]
    user_unfollow = ["W"]
    user_unmute = ["U"]
```

### 日付・時刻のフォーマットについて

[time パッケージ](https://pkg.go.dev/time#pkg-constants) と同じフォーマット構文が使用できます

### レイアウトのカスタマイズについて

- 項目内で有効なタグ `{tag}` を組み合わせてレイアウトをカスタマイズできます
- 指定したタグに置換する内容が無かった場合、**そのタグ + 後ろにある 1 文字分の空白・改行** が削除されます
  - 例: アノテーションが無い場合 `{annotation} {name}` は `{annotation} ` が削除され `name` のみが表示される
- `tweet` `user` のみ、全てのタグを置換後に末尾が改行で終わる場合、その改行は削除されます

#### tweet

ツイート全体のレイアウト

##### デフォルト

```
{annotation}\n{user_info}\n{text}\n{poll}\n{detail}
```

##### 表示との対応

```
RT by arrow2nd @arrow_2nd                        ── {annotation}
arrow2nd @arrow_2nd                              ── {user_info}
This is test tweet                               ── {text}
test1                                            ┐
███ 0.1% (2)                                     │
test2                                            │
███████ 0.2% (4)                                 │
test3                                            ├─ {poll}
████████████ 0.4% (7)                            │
test4                                            │
███████ 0.2% (4)                                 │
closed | 17 votes | ends on 2022/06/28 10:07:56  ┘
2022/06/21 10:07:56 | via Twitter for Android    ┬─ {detail}
1Like 2RTs                                       ┘
```

#### tweet_anotation

ツイートのアノテーションのレイアウト

##### デフォルト

```
{text} {author_name} {author_username}
```

##### 表示との対応

```
RT by arrow2nd @arrow_2nd
~~~~~ ~~~~~~~~ ~~~~~~~~~~
  │      │         └── {author_username}
  │      └──────────── {author_name}
  └─────────────────── {text}
```

#### tweet_detail

ツイート詳細のレイアウト

##### デフォルト

```
{created_at} | via {via}\n{metrics}
```

##### 表示との対応

```

2022/06/21 10:07:56 | via Twitter for Android
~~~~~~~~~~~~~~~~~~~       ~~~~~~~~~~~~~~~~~~~
1Like 2RTs    │                    └───────── {via}
~~~~~~~~~~    └────────────────────────────── {created_at}
    └──────────────────────────────────────── {metrics}
```

#### tweet_poll

投票全体のレイアウト

##### デフォルト

```
{graph}\n{detail}
```

##### 表示との対応

```
test1                                            ┐
███ 0.1% (2)                                     │
test2                                            │
███████ 0.2% (4)                                 │
test3                                            ├─ {graph}
████████████ 0.4% (7)                            │
test4                                            │
███████ 0.2% (4)                                 ┘
closed | 17 votes | ends on 2022/06/28 10:07:56  ── {detail}
```

#### tweet_poll_graph

投票グラフのレイアウト

##### デフォルト

```
{label}\n{graph} {per} {votes}
```

##### 表示との対応

```
test1 ─────── {label}
███ 0.1% (2)
~~~ ~~~~ ~~~
 │    │   └── {votes}
 │    └────── {per}
 └─────────── {graph}
```

#### tweet_poll_detail

投票詳細のレイアウト

##### デフォルト

```
{status} | {all_votes} votes | ends on {end_date}
```

##### 表示との対応

```
closed | 17 votes | ends on 2022/06/28 10:07:56
~~~~~~   ~~                 ~~~~~~~~~~~~~~~~~~~
  │       │                         └────────── {end_date}
  │       └──────────────────────────────────── {all_votes}
  └──────────────────────────────────────────── {status}
```

#### user

ユーザプロフィールのレイアウト

##### デフォルト

```
{user_info}\n{bio}\n{user_detail}
```

##### 表示との対応

```
         arrow2nd @arrow_2nd           ── {user_info}
            I am super cat             ── {bio}
📍 Japan | 🔗 https://t.co/PTi91DHh5Q  ── {user_detail}
```

#### user_info

ユーザ詳細のレイアウト

- `tweet` と `user` 内共通のレイアウトです

##### デフォルト

```
{name} {username} {badge}
```

##### 表示との対応

```
arrow2nd @arrow_2nd ✅ 🔒
~~~~~~~~ ~~~~~~~~~~ ~~~~~
   │         │        └──── {badge}
   │         └───────────── {username}
   └─────────────────────── {name}
```

## style_default.toml

デフォルトのスタイル定義ファイルです

### 設定項目の構文について

```toml
# アプリ全体
[app]
  # 背景色
  background_color = "#000000"
  # 罫線
  border_color = "#ffffff"
  # 文字
  text_color = "#f9f9f9"
  # プレースホルダ文字
  sub_text_color = "#979797"
  # 注意・警告文字
  emphasis_text = "maroon:-:bi"

# タブバー
[tab]
  # 文字
  text = "white:-:-"
  # 背景
  background_color = "#000000"

# 補完候補
[autocomplete]
  # 文字
  text_color = "#000000"
  # 背景
  background_color = "#808080"
  # 選択中の補完候補背景
  selected_background_color = "#C0C0C0"

# ステータスバー
[statusbar]
  # 文字
  text = "black:-:-"
  # 背景
  background_color = "#ffffff"

# ツイート
[tweet]
  # アノテーション（RT by みたいなやつ）
  annotation = "teal:-:-"
  # ツイート詳細（投稿日時, via）
  detail = "gray:-:-"
  # いいね数
  like = "pink:-:-"
  # リツイート数
  retweet = "lime:-:-"
  # ハッシュタグ
  hashtag = "aqua:-:-"
  # メンション
  mention = "aqua:-:-"
  # 投票グラフ
  poll_graph = "aqua:-:-"
  # 投票詳細（開催ステータス, 総投票数, 終了日時）
  poll_detail = "gray:-:-"
  # セパレータ
  separator = "gray:-:-"

# ユーザ
[user]
  # 表示名
  name = "white:-:b"
  # ユーザ名（@arrow_2nd みたいなの）
  user_name = "gray:-:i"
  # ユーザ詳細（位置情報, URL）
  detaill = "gray:-:-"
  # 認証済みバッジ
  verified = "blue:-:-"
  # 非公開バッジ
  private = "gray:-:-"

# ユーザメトリクス
[metrics]
  # ツイート数 / 文字
  tweets_text = "black:-:-"
  # ツイート数 / 背景
  tweets_background_color = "#a094c7"
  # フォロイー数 / 文字
  following_text = "black:-:-"
  # フォロイー数 / 背景
  following_background_color = "#84a0c6"
  # フォロワー数 / 文字
  followers_text = "black:-:-"
  # フォロワー数 / 背景
  followers_background_color = "#89b8c2"
```

#### 末尾が `_color` で終わる項目

`#` から始まる、16 進数カラーコードが使用できます

#### それ以外の項目

[tview の Color tag](https://pkg.go.dev/github.com/rivo/tview#hdr-Colors) が使用できます

> 構文: `<前景色>:<背景色>:<フラグ>`
