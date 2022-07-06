# nekome

🐈 ねこのための TUI な Twitter クライアント

[![Go Report Card](https://goreportcard.com/badge/github.com/arrow2nd/nekome)](https://goreportcard.com/report/github.com/arrow2nd/nekome)
[![GitHub license](https://img.shields.io/github/license/arrow2nd/nekome)](https://github.com/arrow2nd/nekome/blob/main/LICENSE)

![nekome](https://user-images.githubusercontent.com/44780846/177174791-d5fb9db2-2a83-490a-8ed0-7d08fe16f89c.gif)

## 特徴

- Twitter API v2 対応
- マルチアカウント対応
- コマンドラインからのツイートが可能
- 柔軟な機能・外観設定

## インストール

> **Warning**
>
> 以下の方法以外でインストールした場合、コンシューマキーが内蔵されていません。
>
> [Twitter Developer Portal](https://developer.twitter.com/en/portal/projects-and-apps) から Twitter API v2 の API キーを取得して、`settings.yml` に設定してください。

### Homebrew

```sh
brew tap arrow2nd/tap
brew install nekome
```

### Scoop

```
scoop bucket add arrow2nd https://github.com/arrow2nd/scoop-bucket.git
scoop install arrow2nd/nekome
```

### バイナリ

[Releases](https://github.com/arrow2nd/nekome/releases) からお使いの環境にあったファイルをダウンロードしてください。

## 使い方

### 初期設定

### コマンド

[コマンド一覧](./docs/commands.md)、もしくは `nekome -h` をご覧ください。

### ショートカット

[ショートカット一覧](./docs/shortcuts.md)、もしくは アプリ内で `?` を入力しヘルプをご覧ください。

### 設定

[設定ファイル](./docs/config.md)をご覧ください。

## 由来

`物事がめまぐるしく変化すること` を指す 猫の目 という言葉が由来です。

> https://nekojiten.com/wp/nekonome/
