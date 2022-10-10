# Change Log

## [Unreleased]

### Fixed

- ユーザプロフィールが見切れる・空行が表示されることがある

## [v2.0.5] - 2022-10-09

### Fixed

- 内部エディタが起動しない
- 外部エディタ起動後に画面が再描画されない

### Changed

- エラー出力の形式を変更

## [v2.0.4] - 2022-10-09

### Fixed

- サブコマンドのフラグがパースできない
- `nekome tweet` で外部エディタが起動できない

### Changed

- `tweet`, `search` でスペースを含む文字列をそのまま受け取れるよう変更

## [v2.0.3] - 2022-10-04

### Fixed

- 設定項目の typo

## [v2.0.2] - 2022-10-04

### Fixed

- 投票グラフのパーセント表示がおかしい

## [v2.0.1] - 2022-10-04

### Fixed

- 自動リリースに使用しているモジュールのパスを修正
- Notice メッセージの URL の typo を修正

## [v2.0.0] - 2022-10-04

### Added

- nekome 内での複数行のツイート編集機能を追加
- キーバインドのカスタマイズを追加
- レイアウトのカスタマイズを追加
- カスタマイズできる配色・スタイル設定の項目を追加
- カスタマイズできる外観設定の項目を追加
- キーバインドを追加
  - 新規ツイート作成
  - ユーザのいいねリストを開く

### Fixed

- アプリ終了モーダル表示中にドキュメントを開くことができる
- テーマが白ベースの端末で表示が見にくい
- 補完中に `<ESC>` 等で入力を抜けると補完リストが画面に残り続ける場合がある
- コマンドライン入力中に `ctrl+q` でモーダルを表示させると補完リストが画面に残り続ける

### Changed

- 設定ファイル群の形式を yaml から toml に変更
- カスタムコンシューマキーを `.cred.toml` から読み込むよう変更
- ショートカットの呼称をキーバインドに変更
- モーダルのボーダースタイルを変更
- カーソル移動の上下ループを廃止
- `:<TAB>` でコマンド補完を開始するよう変更

## [v1.2.0] - 2022-08-18

### Fixed

- アンフォロー時の確認モーダルが表示されない

### Changed

- ページ追加時に重複するページがあった場合、既にあるページに移動するよう変更
- 確認モーダルの表示を変更

## [v1.1.0] - 2022-08-13

### Added

- `tweet` コマンドを標準入力に対応
- ホームタイムラインページに Stream Mode を追加

### Changed

- リロードした際にカーソルの位置を保持するよう変更

### Fixed

- 非アクティブ状態のページがインジケータを更新してしまう

## [v1.0.3] - 2022-08-02

### Fixed

- ハッシュタグのハイライト処理に失敗する

## [v1.0.2] - 2022-07-23

### Fixed

- コピーした URL の形式が不正 #24

## [v1.0.1] - 2022-07-08

### Changed

- コンシューマキーが無い状態でも `edit` コマンドが実行できるよう変更

## [v1.0.0] - 2022-07-07

- リリースしました！ 😸

[unreleased]: https://github.com/arrow2nd/nekome/compare/v2.0.5...HEAD
[v2.0.5]: https://github.com/arrow2nd/nekome/compare/v2.0.4...v2.0.5
[v2.0.4]: https://github.com/arrow2nd/nekome/compare/v2.0.3...v2.0.4
[v2.0.3]: https://github.com/arrow2nd/nekome/compare/v2.0.2...v2.0.3
[v2.0.2]: https://github.com/arrow2nd/nekome/compare/v2.0.1...v2.0.2
[v2.0.1]: https://github.com/arrow2nd/nekome/compare/v2.0.0...v2.0.1
[v2.0.0]: https://github.com/arrow2nd/nekome/compare/v1.1.0...v2.0.0
[v1.2.0]: https://github.com/arrow2nd/nekome/compare/v1.1.0...v1.2.0
[v1.1.0]: https://github.com/arrow2nd/nekome/compare/v1.0.3...v1.1.0
[v1.0.3]: https://github.com/arrow2nd/nekome/compare/v1.0.2...v1.0.3
[v1.0.2]: https://github.com/arrow2nd/nekome/compare/v1.0.1...v1.0.2
[v1.0.1]: https://github.com/arrow2nd/nekome/compare/v1.0.0...v1.0.1
[v1.0.0]: https://github.com/arrow2nd/nekome/compare/v0.0.0...v1.0.0
