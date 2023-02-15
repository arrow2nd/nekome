![nekome-readme](https://user-images.githubusercontent.com/44780846/204079320-eb71727d-e7e8-4160-92f4-4bb6b9a0ea9e.png)

**nekome**: 🐈 TUI Twitter client for cats

[![release](https://github.com/arrow2nd/nekome/actions/workflows/release.yml/badge.svg)](https://github.com/arrow2nd/nekome/actions/workflows/release.yml)
[![test](https://github.com/arrow2nd/nekome/actions/workflows/test.yml/badge.svg)](https://github.com/arrow2nd/nekome/actions/workflows/test.yml)
[![CodeQL](https://github.com/arrow2nd/nekome/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/arrow2nd/nekome/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/arrow2nd/nekome)](https://goreportcard.com/report/github.com/arrow2nd/nekome)
[![GitHub all releases](https://img.shields.io/github/downloads/arrow2nd/nekome/total)](https://github.com/arrow2nd/nekome/releases)

> [日本語](./README.md)

> **Note**
>
> Due to the [Discontinuation of the free Twitter API plan](https://twitter.com/TwitterDev/status/1621026986784337922), the development of nekome will end on 2023/02/09.
>
> The built-in API token may have already expired. In that case, you can continue to use the API by preparing and setting up your own API token.
> Please refer to [this page](https://github.com/arrow2nd/nekome/blob/v2/docs/en/config.md#credtoml) for the setting locations.

![nekome](https://user-images.githubusercontent.com/44780846/210126086-2be3feab-3ad9-41f5-9510-d28b947256f4.gif)

## Features

- Twitter API v2 support
- Multi-account support
- Tweeting from the command line is possible
- Flexible feature and appearance settings

## Installation

> **Warning**
>
> If you install the software in a manner other than the following, the consumer
> key is not built into the software
>
> Get your Twitter API v2 API key from
> [Twitter Developer Portal](https://developer.twitter.com/en/portal/projects-and-apps)
> and add it to the [.cred.toml](./docs/en/config.md#credtoml) generated after
> startup

### Homebrew

```
brew tap arrow2nd/tap
brew install nekome
```

### Scoop

```
scoop bucket add arrow2nd https://github.com/arrow2nd/scoop-bucket.git
scoop install arrow2nd/nekome
```

### Go install

```
go install github.com/arrow2nd/nekome/v2@latest
```

### Binary

Download the appropriate file for your environment from
[Releases](https://github.com/arrow2nd/nekome/releases)

## Initialization

![image](https://user-images.githubusercontent.com/44780846/177674269-2efa3342-bb1a-4be3-8133-7fc8f6e8cec0.png)

1. The URL of the authentication page is displayed at the first startup, so
   access it with a browser
2. Follow the on-screen instructions for authentication and copy the PIN code
   displayed
3. Enter PIN code into nekome
4. Done! 🐱

## Documents

- [Commands List](./docs/en/commands.md) or `nekome -h`
- [Default Keybindings](./docs/en/keybindings.md) or `docs keybindings`
- [About Configuration File](./docs/en/config.md)
- [Migrate from v1 to v2](./docs/en/migrate-v1-v2.md)
- [Sample Styles](./docs/sample_styles.md)

## Origin of name

The name comes from the Japanese word "猫の目 (neko no me)" which means "the eye of
a cat" and refers to things changing at a dizzying pace

> https://nekojiten.com/wp/nekonome/

## Sincere thanks ❤

- [SobaCha / MateCha (@wakamesoba98)](https://www.wakamesoba98.net/sobacha/index.html)
- [Twitter Developer Community](https://twittercommunity.com/)
