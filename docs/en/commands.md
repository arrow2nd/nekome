# List of commands

### nekome

```
nekome [flags] [command]
```

#### Flags

- `-u <userName>` `--user <userName>`
  - Specify user to use
- `-v` `--version`
  - Show version

#### Global Flags

- `-h` `--help`
  - Show help

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

## Commands available from command line

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

## Commands available from within the app

### home

Add home timeline page

```
:home [flags]
```

#### flags

- `-u` `--unfocus`
  - Do not focus when adding page

### mention

Add mention timeline page

```
:mention [flags]
```

#### flags

- `-u` `--unfocus`
  - Do not focus when adding page

### list

Add list timeline page

```
:list [flags] <list name> <list id>
```

#### flags

- `-u` `--unfocus`
  - Do not focus when adding page

### user

Add user timeline page

```
:user [flags] [user name]
```

#### flags

- `-u` `--unfocus`
  - Do not focus when adding page

### search

Add seaech result page

```
:search [flags] <query>
```

#### flags

- `-u` `--unfocus`
  - Do not focus when adding page

### docs

Show the document

```
:docs [command]
```

#### shortcuts

Documentation for shortcut keys

```
:docs shortcuts [flags]
```

#### flags

- `-u` `--unfocus`
  - Do not focus when adding page

### quit

Quit the application

```
:quit
```
