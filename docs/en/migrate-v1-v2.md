# Migrate from v1 to v2

> [日本語](../ja/migrate-v1-v2.md)

nekome v2 mainly introduces destructive changes around the configuration file.

For more information, see [About Configuration File](. /config.md) for details.

## 1. re-configure custom consumer key

The location of the consumer key has been changed from `settings.yml` to `.cred.toml`.

Also, the file name of the credentials has been changed from `.cred` to `.cred.toml`.

### Configuration correspondence table

| v1                           | v2                   |
| ---------------------------- | -------------------- |
| feature.consumer.token       | consumer.Token       |
| feature.consumer.tokensecret | consumer.TokenSecret |

## 2. Migrating preferences

File name changed from `settings.yml` to `preferences.toml`.

#### Configuration correspondence table

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

## 3. Migrating style definition file

The style definition file format has also been changed to `.toml`.

See [Sample Styles](../sample_styles.md) for a sample.

### Destructive changes

- Background color can now be set.
  - In v1, it was the terminal background color. v2 uses black by default.
  - If you want to use the terminal background color, please set an empty character (`""`)

#### Configuration correspondence table

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
| user.detail               | user.detaill                           |
| user.tweetsmetricstext    | metrics.tweets_text                    |
| user.tweetsmetricsbg      | metrics.tweets_background_color        |
| user.followingmetricstext | metrics.following_text                 |
| user.followingmetricsbg   | metrics.following_background_color     |
| user.followersmetricstext | metrics.followers_text                 |
| user.followersmetricsbg   | metrics.followers_background_color     |

## 4. Delete old configuration files

```sh
rm ~/.config/nekome/{.cred,*.yml}
```
