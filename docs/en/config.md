# About Configuration File

> [日本語](../ja/config.md)

The configuration file is created and saved as follows

```
$HOME/.config/nekome
├── .cred.toml
├── style_default.toml
└── preferences.toml
```

## .cred.toml

Credential data file

```toml
# Twitter API consumer key
[consumer]
  Token = ""
  TokenSecret = ""

# CAUTION :
# DO NOT EDIT THE FOLLOWING MANUALLY
# Use `nekome account` to manipulate user information

# User credentials
[user]

  [[user.accounts]]
    UserName = "user_name"
    ID = "0123456789"
    [user.accounts.Token]
      Token = "hoge"
      TokenSecret = "fuga"
```

## preferences.toml

Preferences file

```toml
[feature]
  # User used by default
  main_user = "user_name"
  # Number of tweets read at one time
  load_tweets_limit = 25
  # Maximum number of tweets accumulated on a page
  accmulate_tweets_limit = 250
  # Whether to launch an external editor when editing tweets
  use_external_editor = false
  # Whether the locale setting of the execution environment is CJK or not
  # (countermeasure for display disorder of tview)
  is_locale_cjk = true
  # Commands to be executed at startup
  startup_cmds = ["home", "mention --unfocus"]

# Whether to display a confirmation modal or not
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
  # Style files to load
  style_file = "style_default.toml"
  # Date format
  date_fmt = "2006/01/02"
  # Time format
  time_fmt = "15:04:05"
  # Maximum number of lines of BIO on user page
  user_bio_max_row = 3
  # Padding on left and right side of user page profile display area
  user_profile_padding_x = 4
  # Separator for user details
  user_detail_separator = " | "
  # Hide separator between tweets
  hide_tweet_separator = false
  # Hide separator for quoted tweets
  hide_quote_tweet_separator = false
  # Characters used to display the graph
  graph_char = "█"
  # Maximum width of graph
  graph_max_width = 30
  # Separator of Tab
  tab_separator = "|"
  # Tab maximum width
  tab_max_width = 20

[layout]
  # Tweet
  tweet = "{annotation}\n{user_info}\n{text}\n{poll}\n{detail}"
  # Annotate about tweet
  tweet_anotation = "{text} {author_name} {author_username}"
  # Tweet Details
  tweet_detail = "{created_at} | via {via}\n{metrics}"
  # Polls
  tweet_poll = "{graph}\n{detail}"
  # Polling Graphs
  tweet_poll_graph = "{label}\n{graph} {per} {votes}"
  # Polling Details
  tweet_poll_detail = "{status} | {all_votes} votes | ends on {end_date}"
  # User profile
  user = "{user_info}\n{bio}\n{user_detail}"
  # User information
  user_info = "{name} {username} {badge}"

[text]
  # Unit of likes
  like = "Like"
  # Unit of retweets
  retweet = "RT"
  # Display loading
  loading = "Loading..."
  # Display when there are no tweets
  no_tweets = "No tweets ฅ^-ω-^ฅ"
  # Display tab text
  tab_home = "Home"
  tab_mention = "Mention"
  tab_list = "List: {name}"
  tab_user = "User: @{name}"
  tab_search = "Search: {query}"
  tab_likes = "Likes: @{name}"
  tab_docs = "Docs: {name}"

[icon]
  # Profile location
  geo = "📍"
  # Profile website
  link = "🔗"
  # Pinned tweet
  pinned = "📌"
  # Verified badge
  verified = "✅"
  # Private badge
  private = "🔒"

[keybinding]
  # Global
  [keybinding.global]
    quit = ["ctrl+q"]
  # Main view
  [keybinding.view]
    close_page = ["ctrl+w"]
    focus_cmdline = [":"]
    redraw = ["ctrl+l"]
    select_next_tab = ["l", "Right"]
    select_prev_tab = ["h", "Left"]
    show_help = ["?"]
  # All pages in common
  [keybinding.page]
    reload_page = ["."]
  # Home timeline page
  [keybinding.home_timeline]
    stream_mode_start = ["s"]
    stream_mode_stop = ["S"]
  # Tweet view
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

### About date and time formats

Uses the same format as [time package](https://pkg.go.dev/time#pkg-constants)

### About layout customization

- You can customize the layout by combining valid tags `{tag}` within an item
- If there is no content to replace for a given tag, **that tag + one trailing space/line break** will be removed
  - Example: If a tweet has no annotation `{annotation} {name}`, `{annotation} ` will be removed and only `name` will be displayed
- Only for `tweet` `user`, if all tags are replaced and end with a newline, the newline will be removed

#### tweet

Layout of overall tweet

##### Default

```
{annotation}\n{user_info}\n{text}\n{poll}\n{detail}
```

##### Correspondence between display and settings

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

Layout of annotations in tweets

##### Default

```
{text} {author_name} {author_username}
```

##### Correspondence between display and settings

```
RT by arrow2nd @arrow_2nd
~~~~~ ~~~~~~~~ ~~~~~~~~~~
  │      │         └── {author_username}
  │      └──────────── {author_name}
  └─────────────────── {text}
```

#### tweet_detail

Layout of tweet detail

##### Default

```
{created_at} | via {via}\n{metrics}
```

##### Correspondence between display and settings

```

2022/06/21 10:07:56 | via Twitter for Android
~~~~~~~~~~~~~~~~~~~       ~~~~~~~~~~~~~~~~~~~
1Like 2RTs    │                    └───────── {via}
~~~~~~~~~~    └────────────────────────────── {created_at}
    └──────────────────────────────────────── {metrics}
```

#### tweet_poll

Layout of overall poll

##### Default

```
{graph}\n{detail}
```

##### Correspondence between display and settings

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

Layout of poll graph

##### Default

```
{label}\n{graph} {per} {votes}
```

##### Correspondence between display and settings

```
test1 ─────── {label}
███ 0.1% (2)
~~~ ~~~~ ~~~
 │    │   └── {votes}
 │    └────── {per}
 └─────────── {graph}
```

#### tweet_poll_detail

Layout of poll detail

##### Default

```
{status} | {all_votes} votes | ends on {end_date}
```

##### Correspondence between display and settings

```
closed | 17 votes | ends on 2022/06/28 10:07:56
~~~~~~   ~~                 ~~~~~~~~~~~~~~~~~~~
  │       │                         └────────── {end_date}
  │       └──────────────────────────────────── {all_votes}
  └──────────────────────────────────────────── {status}
```

#### user

Layout of user profile

##### Default

```
{user_info}\n{bio}\n{user_detail}
```

##### Correspondence between display and settings

```
         arrow2nd @arrow_2nd           ── {user_info}
            I am super cat             ── {bio}
📍 Japan | 🔗 https://t.co/PTi91DHh5Q  ── {user_detail}
```

#### user_info

Layout of user info

- Common layout in `tweet` and `user`

##### Default

```
{name} {username} {badge}
```

##### Correspondence between display and settings

```
arrow2nd @arrow_2nd ✅ 🔒
~~~~~~~~ ~~~~~~~~~~ ~~~~~
   │         │        └──── {badge}
   │         └───────────── {username}
   └─────────────────────── {name}
```

## style_default.toml

The default style definition file.

```toml
[app]
  # Background color of the entire app
  background_color = "#000000"
  # Border
  border_color = "#ffffff"
  # Text
  text_color = "#f9f9f9"
  # Placeholder
  sub_text_color = "#979797"
  # Caution and warning text
  emphasis_text = "maroon:-:bi"

[tab]
  # Text
  text = "white:-:-"
  # Background
  background_color = "#000000"

[autocomplete]
  # Item Text
  text_color = "#000000"
  # Unselected item background
  background_color = "#808080"
  # Selecting item background
  selected_background_color = "#C0C0C0"

[statusbar]
  # Text
  text = "black:-:-"
  # Background
  background_color = "#ffffff"

[tweet]
  # Annotation (RT by ...)
  annotation = "teal:-:-"
  # Detail (date, via)
  detail = "gray:-:-"
  # Likes
  like = "pink:-:-"
  # Retweets
  retweet = "lime:-:-"
  # Hashtag
  hashtag = "aqua:-:-"
  # Mention
  mention = "aqua:-:-"
  # Graph
  poll_graph = "aqua:-:-"
  # Poll detail (status, total votes, end date)
  poll_detail = "gray:-:-"
  # Separator
  separator = "gray:-:-"

[user]
  # Nickname
  name = "white:-:b"
  # Username (starting with @)
  user_name = "gray:-:i"
  # User detail (Geo, URL)
  detail = "gray:-:-"
  # Verified badge
  verified = "blue:-:-"
  # Private badge
  private = "gray:-:-"

[metrics]
  # Tweets / Text
  tweets_text = "black:-:-"
  # Tweets / Background
  tweets_background_color = "#a094c7"
  # Following / Text
  following_text = "black:-:-"
  # Following / Background
  following_background_color = "#84a0c6"
  # Followers / Text
  followers_text = "black:-:-"
  # Followers / Background
  followers_background_color = "#89b8c2"
```

### Syntax of configuration items

#### Items ending with `_color`

Hexadecimal color codes starting with `#` or W3C color names can be used

#### Other items

[Color tag for tview](https://pkg.go.dev/github.com/rivo/tview#hdr-Colors) can be used

> Syntax: `<foreground>:<background>:<flags>`
