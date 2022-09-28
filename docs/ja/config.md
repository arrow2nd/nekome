# 設定ファイルについて

以下の様な形で作成・保存されます

```
$HOME/.config/nekome
├── .cred.toml
├── style_default.toml
└── settings.toml
```

## .cred.toml

認証情報を記録したファイルです

### 例

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

## settings.toml

環境設定のファイルです

### 日付・時刻のフォーマットについて

[time パッケージ](https://pkg.go.dev/time#pkg-constants) と同じフォーマット構文が使用できます

### 例

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
  # 読み込むスタイルファイル
  style_file = "style_default.toml"
  # 日付のフォーマット
  date_fmt = "2006/01/02"
  # 時刻のフォーマット
  time_fmt = "15:04:05"
  # ユーザページ / BIO の最大行数
  user_bio_max_row = 3
  # ユーザページ / プロフィール表示域の左右パディング
  user_profile_padding_x = 4
  # グラフ / 表示に使用する文字
  graph_char = "█"
  # グラフ / 最大表示幅
  graph_max_width = 30
  # タブ / セパレータ文字
  tab_separate = "|"
  # タブ / 最大表示幅
  tab_max_width = 20

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
```

## style_default.toml

デフォルトのスタイル定義ファイルです

`settings.toml` 内の `appearance.style_file` に指定したファイルが読み込まれます

### 設定項目の構文について

#### 末尾が `_color` で終わる項目

`#` から始まる、16 進数カラーコードが使用できます

#### それ以外の項目

[tview の Color tag](https://pkg.go.dev/github.com/rivo/tview#hdr-Colors) が使用できます

> 構文: `<前景色>:<背景色>:<フラグ>`

### 例

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
  # ユーザ詳細（BIO）
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
