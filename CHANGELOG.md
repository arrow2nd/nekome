# Change Log

## [Unreleased]

### Added

- nekome 内での複数行のツイート編集機能を追加
- カスタムできる配色・スタイルの項目を追加
- キーバインドのカスタム機能を追加
- キーバインドを追加
  - 新規ツイート作成
  - ユーザのいいねリストを開く
- `:<TAB>` でのコマンド補完を追加

### Fixed

- アプリ終了モーダル表示中にドキュメントを開くことができる
- テーマが白ベースの端末で表示が見にくい

### Changed

- 設定ファイル群の形式を yaml から toml に変更
- カスタムコンシューマキーを `.cred.toml` から読み込むよう変更
- ショートカットの呼称をキーバインドに変更
- モーダルのボーダースタイルを変更
- カーソル移動の上下ループを廃止

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

[unreleased]: https://github.com/arrow2nd/nekome/compare/v1.2.0...HEAD
[v1.2.0]: https://github.com/arrow2nd/nekome/compare/v1.1.0...v1.2.0
[v1.1.0]: https://github.com/arrow2nd/nekome/compare/v1.0.3...v1.1.0
[v1.0.3]: https://github.com/arrow2nd/nekome/compare/v1.0.2...v1.0.3
[v1.0.2]: https://github.com/arrow2nd/nekome/compare/v1.0.1...v1.0.2
[v1.0.1]: https://github.com/arrow2nd/nekome/compare/v1.0.0...v1.0.1
[v1.0.0]: https://github.com/arrow2nd/nekome/compare/v0.0.0...v1.0.0
