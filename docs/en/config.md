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

### Example

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

### About date and time formats

Uses the same format as [time package](https://pkg.go.dev/time#pkg-constants)

### Example

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
  # Characters used to display the graph
  graph_char = "█"
  # Maximum width of graph
  graph_max_width = 30
  # Tab separator
  tab_separate = "|"
  # Tab maximum width
  tab_max_width = 20

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

## style_default.toml

The default style definition file.

The file specified in `appearance.style_file` in `preferences.toml` is loaded.

### Syntax of configuration items

#### Items ending with `_color`

Hexadecimal color codes beginning with `#` can be used

#### Other items

[Color tag for tview](https://pkg.go.dev/github.com/rivo/tview#hdr-Colors) can be used

> Syntax: `<foreground>:<background>:<flags>`

### Example

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
  detaill = "gray:-:-"
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
