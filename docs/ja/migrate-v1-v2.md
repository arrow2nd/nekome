# v1 から v2 へ移行する

> [English](../en/migrate-v1-v2.md)

nekome v2 では主に、設定ファイル周りに破壊的な変更がなされています。

詳しくは [設定ファイルについて](./config.md) をご覧ください。

## 1. カスタムコンシューマキーを再設定

コンシューマキーの設定箇所が `settings.yml` 内から `.cred.toml` 内へ変更されました。

また、これに伴い、認証情報のファイル名が `.cred` から `.cred.toml` に変更されています。

### 設定の対応表

| v1                           | v2                   |
| ---------------------------- | -------------------- |
| feature.consumer.token       | consumer.Token       |
| feature.consumer.tokensecret | consumer.TokenSecret |

## 2. 環境設定を移行

ファイル名が `settings.yml` から`preferences.toml` に変更されました。

### 設定の対応表

| v1                              | v2                                |
| ------------------------------- | --------------------------------- |
| feature.mainuser                | feature.main_user                 |
| feature.loadtweetscount         | feature.load_tweets_limit         |
| feature.tweetmaxaccumulationnum | feature.accmulate_tweets_limit    |
| feature.islocalecjk             | feature.is_locale_cjk             |
| feature.confirm.\*              | confirm.\*                        |
| feature.startup                 | feature.startup_cmds              |
| appearance.stylefile            | stylefile.style_file              |
| appearance.dateformat           | appearance.date_fmt               |
| appearance.timeformat           | appearance.time_fmt               |
| appearance.userbiomaxrow        | appearance.user_bio_max_row       |
| appearance.userprofilepaddingx  | appearance.user_profile_padding_x |
| appearance.graphchar            | appearance.graph_char             |
| appearance.graphmaxwidth        | appearance.graph_max_width        |
| appearance.tabseparate          | appearance.tab_separator          |
| appearance.tabmaxwidth          | appearance.tab_max_width          |
| texts.like                      | text.like                         |
| texts.retweet                   | text.retweet                      |
| texts.loading                   | text.loading                      |
| texts.notweets                  | text.no_tweets                    |
| texts.tabhome                   | text.tab_home                     |
| texts.tabmention                | text.tab_mention                  |
| texts.tablist                   | text.tab_list                     |
| texts.tabuser                   | text.tab_user                     |
| texts.tabsearch                 | text.tab_search                   |
| texts.tabdocs                   | text.tab_docs                     |
| icon.\*                         | icon.\*                           |

## 3. スタイル定義ファイルを移行

スタイル定義ファイルの形式も `.toml` へ変更されました。

サンプルは [Sample Styles](../sample_styles.md) をご覧ください。

### 破壊的な変更点

- 背景色の設定が可能になりました
  - v1 では端末の背景色でしたが v2 からはデフォルトで黒色が使用されます
  - 端末の背景色を使用する場合、空文字（`""`）を設定してください

### 設定の対応表

| v1                        | v2                                     |
| ------------------------- | -------------------------------------- |
| app.tab                   | tab.text                               |
| app.separator             | tweet.separator                        |
| statusbar.text            | statusbar.text                         |
| statusbar.bg              | statusbar.background_color             |
| autocomplete.normalbg     | autocomplete.background_color          |
| autocomplete.selectbg     | autocomplete.selected_background_color |
| tweet.annotation          | tweet.annotation                       |
| tweet.detail              | tweet.detail                           |
| tweet.like                | tweet.like                             |
| tweet.rt                  | tweet.retweet                          |
| tweet.hashtag             | tweet.hashtag                          |
| tweet.mention             | tweet.mention                          |
| tweet.pollgraph           | tweet.poll_graph                       |
| tweet.polldetail          | tweet.poll_detail                      |
| user.name                 | user.name                              |
| user.username             | user.user_name                         |
| user.verified             | user.verified                          |
| user.private              | user.private                           |
| user.detail               | user.detail                            |
| user.tweetsmetricstext    | metrics.tweets_text                    |
| user.tweetsmetricsbg      | metrics.tweets_background_color        |
| user.followingmetricstext | metrics.following_text                 |
| user.followingmetricsbg   | metrics.following_background_color     |
| user.followersmetricstext | metrics.followers_text                 |
| user.followersmetricsbg   | metrics.followers_background_color     |

## 4. 古い設定ファイルを削除

```sh
rm ~/.config/nekome/{.cred,*.yml}
```
