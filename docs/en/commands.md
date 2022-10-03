# Commands List

> [日本語](../ja/commands.md)

## nekome

```
nekome [flags] [command] [global flags / command flags]
```

### Global Flags

- `-h` `--help`
  - Show help

### Flags

- `-u <userName>` `--user <userName>`
  - Specify user to use
- `-v` `--version`
  - Show version

## Commands available in common

### tweet

Posts a tweet

If you omit the tweet statement, the editor will be activated.

```
nekome tweet [flags] [text]
:tweet [flags] [text]
```

#### Flags

- `-e <editor command>` `--editor <editor command>`
  - Specify which editor to use
  - If omitted, the value of `$EDITOR` is specified
- `-i <file path>` `--image <file path>`
  - Image to be attached
  - To specify multiple images, separate them with `,`
- `-q <tweet id>` `--quote <tweet id>`
  - Specify the ID of the tweet to quote
- `-r <tweet ID>` `--reply <tweet ID>`
  - Specify the ID of the tweet to which you are replying

## Commands available from CLI

### account

Manage your account

```
nekome account [command]
```

#### add

Add account

```
nekome account add
```

#### delete

Delete account

```
nekome account delete
```

#### list

Show accounts that have been added

```
nekome account list
```

### edit

Edit configuration file

```
nekome edit [flags]
```

#### flags

- `-e <editor command>` `--editor <editor command>`
  - Specify which editor to use
  - If omitted, the value of `$EDITOR` will be specified

## Commands available from TUI

### home

Add home timeline page

```
:home [flags]
```

### mention

Add mention timeline page

```
:mention [flags]
```

### list

Add list timeline page

```
:list [flags] <list name> <list id>
```

### user

Add user timeline page

```
:user [flags] [user name]
```

### likes

Add a user's Likes page

```
:likes [flags] [user name]
```

### search

Add seaech result page

```
:search [flags] <query>
```

### docs

Show the document

```
:docs [command]
```

#### keybindings

Documentation for keybindings

```
:docs keybindings [flags]
```

### quit

Quit the application

```
:quit
```
