# 設定ファイルについて

以下の様な形で作成・保存されます。

```
$HOME/.config/nekome
├── .cred
├── default.yml
└── settings.yml
```

## 環境設定

### 日付・時刻のフォーマット形式

[time パッケージ](https://pkg.go.dev/time#pkg-constants) と同じ書式を使用しています。

### 例

> settings.yml

```yaml
# 機能
feature:
  # コンシューマキー（空の場合、内蔵のコンシューマキーが使用されます）
  consumer:
    token: ""
    tokensecret: ""
  # メインユーザ
  mainuser: "arrow_2nd"
  # 1回で読込むツイート数
  loadtweetscount: 25
  # 1ページにおけるツイートの最大蓄積数
  tweetmaxaccumulationnum: 250
  # ツイート編集時に外部エディタを起動するか
  usetweetwhenexeditor: false
  # 実行環境のロケール設定が CJK かどうか（tviewの表示乱れ対策）
  islocalecjk: true
  # 確認モーダルを表示するか
  confirm:
    Block: true
    Delete: true
    Follow: true
    Like: true
    Mute: true
    Quit: true
    Retweet: true
    Tweet: true
    Unblock: true
    Unfollow: true
    Unlike: true
    Unmute: true
    Unretweet: true
  # 起動時に実行するコマンド
  startup:
    - home
    - mention --unfocus

# 外観
appearance:
  # 読込むスタイルファイル
  stylefile: default.yml
  # 日付のフォーマット形式
  dateformat: 2006/01/02
  # 時刻のフォーマット形式
  timeformat: "15:04:05"
  # ユーザページのBIOの最大行数
  userbiomaxrow: 3
  # ユーザページのプロフィール表示部分の左右パディング
  userprofilepaddingx: 4
  # グラフの表示に使用する文字
  graphchar: █
  # グラフの最大幅
  graphmaxwidth: 30
  # タブの区切り文字
  tabseparate: "|"
  # タブの最大幅
  tabmaxwidth: 20

# テキスト
texts:
  # いいねの単位
  like: Like
  # リツイートの単位
  retweet: RT
  # 読み込み中の表示
  loading: Loading...
  # ツイートが無い場合の表示
  notweets: No tweets ฅ^-ω-^ฅ
  # ホームタイムラインページのタブ文字列
  tabhome: Home
  # メンションタイムラインページのタブ文字列
  tabmention: Mention
  # リストタイムラインページのタブ文字列
  tablist: "List: {name}"
  # ユーザページのタブ文字列
  tabuser: "User: @{name}"
  # 検索結果ページのタブ文字列
  tabsearch: "Search: {query}"
  # ドキュメントページのタブ文字列
  tabdocs: "Docs: {name}"

# アイコン
icon:
  # プロフィールの位置情報
  geo: 📍
  # プロフィールのURL
  link: 🔗
  # ピン留めツイート
  pinned: 📌
  # 認証済みユーザ
  verified: ✅
  # 非公開ユーザ
  private: 🔒
```

## スタイル

### スタイル構文

末尾が `bg` 以外の項目については [tview の Color tag の構文](https://pkg.go.dev/github.com/rivo/tview#hdr-Colors) を使用しています。

### 末尾が `bg` の項目の構文

`#` から始まる、16 進数カラーコードのみが使用できます。

### 例

> default.yml

```yaml
# アプリ
app:
  # タブ
  tab: -:-:-
  # セパレータ
  separator: gray:-:-

# ステータスバー
statusbar:
  # 文字
  text: black:-:-
  # 背景
  bg: "#ffffff"

# 入力補完窓
autocomplete:
  # 未選択
  normalbg: "#3e4359"
  # 選択中
  selectbg: "#5c6586"

# ツイート
tweet:
  # アノテーション（RT by みたいなやつ）
  annotation: blue:-:-
  # 詳細（投稿日, via）
  detail: gray:-:-
  # いいね数
  like: pink:-:-
  # リツイート数
  rt: green:-:-
  # ハッシュタグ
  hashtag: blue:-:-
  # メンション
  mention: blue:-:-
  # アンケートグラフ
  pollgraph: blue:-:-
  # アンケート詳細（開催ステータス, 総投票数, 終了日時）
  polldetail: gray:-:-

# ユーザ
user:
  # ニックネーム
  name: lightgray:-:b
  # ユーザ名（@から始まるもの）
  username: gray:-:i
  # 認証済みバッジ
  verified: blue:-:-
  # 非公開バッジ
  private: gray:-:-
  # 詳細（位置情報, URL）
  detail: gray:-:-
  # 総ツイート数 : 文字色
  tweetsmetricstext: black:-:-
  # 総ツイート数 : 背景色
  tweetsmetricsbg: "#a094c7"
  # フォロイー数 : 文字色
  followingmetricstext: black:-:-
  # フォロイー数 : 背景色
  followingmetricsbg: "#84a0c6"
  # フォロワー数 : 文字色
  followersmetricstext: black:-:-
  # フォロワー数 : 背景色
  followersmetricsbg: "#89b8c2"
```
