# nekome

🐈 ねこのための TUI な Twitter クライアント

[![release](https://github.com/arrow2nd/nekome/actions/workflows/release.yml/badge.svg)](https://github.com/arrow2nd/nekome/actions/workflows/release.yml)
[![test](https://github.com/arrow2nd/nekome/actions/workflows/test.yml/badge.svg)](https://github.com/arrow2nd/nekome/actions/workflows/test.yml)
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

![image](https://user-images.githubusercontent.com/44780846/177674269-2efa3342-bb1a-4be3-8133-7fc8f6e8cec0.png)

1. 初回起動時に認証ページの URL が表示されるので、ブラウザでアクセス
2. 画面の指示に沿って認証を進め、表示される PIN コードをコピー
3. PIN コードを nekome に入力
4. 完了！ 🐱

### コマンド

[コマンド一覧](./docs/commands.md)、もしくは `nekome -h` をご覧ください。

### ショートカット

[ショートカット一覧](./docs/shortcuts.md)、もしくは アプリ内で `?` を入力しヘルプをご覧ください。

### 設定

[設定ファイルについて](./docs/config.md)をご覧ください。

## 由来

`物事がめまぐるしく変化すること` を指す 猫の目 という言葉が由来です。

> https://nekojiten.com/wp/nekonome/
