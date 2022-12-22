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
  - Specify the editor to use for editing
  - If omitted, the value of `$EDITOR` is specified
- `-i <file path>` `--image <file path>`
  - Attach the image
  - To specify multiple images, separate them with `,`
- `-c` `--clipboard`
  - Attach the image in the clipboard
  - If the --image is specified, it takes precedence
- `-q <tweet id>` `--quote <tweet id>`
  - Quotes the tweet with the specified ID
- `-r <tweet ID>` `--reply <tweet ID>`
  - Send a reply to the tweet with the specified ID

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

#### flags

- `-s` `--stream`
  - Start stream mode

### mention

Add mention timeline page

```
:mention [flags]
```

#### flags

- `-s` `--stream`
  - Start stream mode

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
