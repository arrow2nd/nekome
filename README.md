![nekome-readme](https://user-images.githubusercontent.com/44780846/204079320-eb71727d-e7e8-4160-92f4-4bb6b9a0ea9e.png)

**nekome**: 🐈 ねこのための TUI な Twitter クライアント

[![release](https://github.com/arrow2nd/nekome/actions/workflows/release.yml/badge.svg)](https://github.com/arrow2nd/nekome/actions/workflows/release.yml)
[![test](https://github.com/arrow2nd/nekome/actions/workflows/test.yml/badge.svg)](https://github.com/arrow2nd/nekome/actions/workflows/test.yml)
[![CodeQL](https://github.com/arrow2nd/nekome/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/arrow2nd/nekome/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/arrow2nd/nekome/v2)](https://goreportcard.com/report/github.com/arrow2nd/nekome/v2)
[![GitHub all releases](https://img.shields.io/github/downloads/arrow2nd/nekome/total)](https://github.com/arrow2nd/nekome/releases)

> [English](./README_EN.md)

> **Note**
>
> [Twitter API の無料プラン廃止](https://twitter.com/TwitterDev/status/1621026986784337922)にともない、nekome の開発は 2023/02/09 をもって終了しました。
>
> 組込みの API トークンは既に失効していますが、ご自身で API トークンを用意し、設定していただくことで引き続きご利用いただくことが可能です。
> 設定個所については、[こちら](https://github.com/arrow2nd/nekome/blob/v2/docs/ja/config.md#credtoml) をご参照ください。

![nekome](https://user-images.githubusercontent.com/44780846/210126086-2be3feab-3ad9-41f5-9510-d28b947256f4.gif)

## 特徴

- Twitter API v2 対応
- マルチアカウント対応
- コマンドラインからのツイートが可能
- 柔軟な機能・外観設定

## インストール

> **Warning**
>
> 以下の方法以外でインストールした場合、コンシューマキーが内蔵されていません
>
> [Twitter Developer Portal](https://developer.twitter.com/en/portal/projects-and-apps)
> から Twitter API v2 の API キーを取得して、起動後に生成される
> [.cred.toml](./docs/ja/config.md#credtoml) に追加してください

### Homebrew

```
brew tap arrow2nd/tap
brew install nekome
```

### Scoop

```
scoop bucket add arrow2nd https://github.com/arrow2nd/scoop-bucket.git
scoop install arrow2nd/nekome
```

### バイナリ

[Releases](https://github.com/arrow2nd/nekome/releases)
からお使いの環境にあったファイルをダウンロードしてください

## 初期設定

![image](https://user-images.githubusercontent.com/44780846/177674269-2efa3342-bb1a-4be3-8133-7fc8f6e8cec0.png)

1. 初回起動時に認証ページの URL が表示されるので、ブラウザでアクセス
2. 画面の指示に沿って認証を進め、表示される PIN コードをコピー
3. PIN コードを nekome に入力
4. 完了！ 🐱

## ドキュメント

- [コマンド一覧](./docs/ja/commands.md) もしくは `nekome -h`
- [デフォルトキーバインド](./docs/ja/keybindings.md) もしくは `docs keybindings`
- [設定ファイル](./docs/ja/config.md)
- [v1 から v2 へ移行](./docs/ja/migrate-v1-v2.md)
- [スタイル定義のサンプル](./docs/sample_styles.md)

## 由来

「物事がめまぐるしく変化すること」を指す _"猫の目"_ という言葉が由来です

> https://nekojiten.com/wp/nekonome/

## Sincere thanks ❤

- [SobaCha / MateCha (@wakamesoba98)](https://www.wakamesoba98.net/sobacha/index.html)
- [Twitter Developer Community](https://twittercommunity.com/)
