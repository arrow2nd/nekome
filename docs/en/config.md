# About Configuration File

The configuration file is created and saved as follows

```
$HOME/.config/nekome
├── .cred
├── default.yml
└── settings.yml
```

## Preferences

### Date and time formatting format

Uses the same format as [time package](https://pkg.go.dev/time#pkg-constants).

### Example

> settings.yml

```yaml
feature:
  # Consumer key (if empty, the built-in consumer key is used)
  consumer:
    token: ""
    tokensecret: ""
  # User used by default
  mainuser: "arrow_2nd"
  # Number of tweets read at one time
  loadtweetscount: 25
  # Maximum number of tweets accumulated on a page
  tweetmaxaccumulationnum: 250
  # Whether the locale setting of the execution environment is CJK or not (countermeasure for display disorder of tview)
  islocalecjk: true
  # Whether to display a confirmation modal or not
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
  # Commands to be executed at startup
  startup:
    - home
    - mention --unfocus

appearance:
  # Style files to load
  stylefile: default.yml
  # Date Format
  dateformat: 2006/01/02
  # Time Format
  timeformat: "15:04:05"
  # Maximum number of lines of BIO on user page
  userbiomaxrow: 3
  # Padding on left and right side of user page profile display area
  userprofilepaddingx: 4
  # Characters used to display the graph
  graphchar: █
  # Maximum width of graph
  graphmaxwidth: 30
  # Tab separator
  tabseparate: "|"
  # Tab maximum width
  tabmaxwidth: 20

texts:
  # Unit of Like
  like: Like
  # Unit of retweet
  retweet: RT
  # Loading display
  loading: Loading...
  # Display when there are no tweets
  notweets: No tweets ฅ^-ω-^ฅ
  # Tab String on Home Timeline Page
  tabhome: Home
  # Tab String on Mention Timeline Page
  tabmention: Mention
  # Tab String on List Timeline Page
  tablist: "List: {name}"
  # Tab String on User Page
  tabuser: "User: @{name}"
  # Tab String on Search Result Page
  tabsearch: "Search: {query}"
  # Tab String on Documentation Page
  tabdocs: "Docs: {name}"

icon:
  # Profile Location
  geo: 📍
  # Profile Website
  link: 🔗
  # Pinned Tweet
  pinned: 📌
  # Verified User
  verified: ✅
  # Private User
  private: 🔒
```

## Style

### Style syntax

For items not ending in `bg`, use [syntax for Color tag in tview](https://pkg.go.dev/github.com/rivo/tview#hdr-Colors).

### Syntax for items ending in `bg`

Only hexadecimal color codes beginning with `#` are allowed.

### Example

> default.yml

```yaml
app:
  # Page Tabs
  tab: -:-:-
  # Separator
  separator: gray:-:-

statusbar:
  # Text
  text: black:-:-
  # Background
  bg: "#ffffff"

autocomplete:
  # Unselected Item
  normalbg: "#3e4359"
  # Selecting Item
  selectbg: "#5c6586"

tweet:
  # Annotation (RT by ...)
  annotation: blue:-:-
  # Detail (date, via)
  detail: gray:-:-
  # Likes
  like: pink:-:-
  # Retweets
  rt: green:-:-
  # Hashtag
  hashtag: blue:-:-
  # Mention
  mention: blue:-:-
  # Graph
  pollgraph: blue:-:-
  # Poll Detail (status, total votes, end date)
  polldetail: gray:-:-

user:
  # Nickname
  name: lightgray:-:b
  # Username (starting with @)
  username: gray:-:i
  # Verified Badge
  verified: blue:-:-
  # Private Badge
  private: gray:-:-
  # Detail (location, website)
  detail: gray:-:-
  # Tweets Metrics : Text
  tweetsmetricstext: black:-:-
  # Tweets Metrics : Background
  tweetsmetricsbg: "#a094c7"
  # Following Metrics : Text
  followingmetricstext: black:-:-
  # Following Metrics : Background
  followingmetricsbg: "#84a0c6"
  # Followers Metrics : Text
  followersmetricstext: black:-:-
  # Followers Metrics : Background
  followersmetricsbg: "#89b8c2"
```
