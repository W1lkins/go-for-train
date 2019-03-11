# go-for-train 🚂

A dumb bot that checks train status and sends messages, specific only to me

**Table of Contents**

<!-- toc -->

- [Installation](#installation)
    + [Binaries](#binaries)
    + [Via Go](#via-go)
    + [Running with Docker](#running-with-docker)
- [Usage](#usage)

<!-- tocstop -->

## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/w1lkins/go-for-train/releases).

#### Via Go

```console
$ go get -u -v github.com/w1lkins/go-for-train
```

#### Running with Docker

**Authentication**

Create a Pushover app and grab your app key and client key from a device you want to send messages to.
You'll need to sign up to the National Rail API [here](http://realtime.nationalrail.co.uk/OpenLDBWSRegistration/) to get an API key.

**Run it in daemon mode with Pushover key/token + National Rail API key**

```console
# You need to either have environment variables that are
# PUSHOVER_APP_KEY and PUSHOVER_CLIENT_KEY
# or pass them into the container.
$ docker run -d --restart always \
    --name go-for-train \
    -e PUSHOVER_APP_KEY=foo \
    -e PUSHOVER_CLIENT_KEY=bar \
    -e NATIONAL_RAIL_APP_KEY=baz \
    -v /path/to/config.toml:/home/user/.config/go-for-train/config/config.toml \
    w1lkins/go-for-train -d --interval 15m
```

## Usage

```console
go-for-train -  A bot that checks the status of my train journey and notifies me about it.

Usage: go-for-train <command>

Flags:

  -d, --debug            enable debug logging (default: false)
  --final-hour           Send messages between these hours (upper bound) (default: 17)
  --initial-hour         Send messages between these hours (lower bound) (default: 15)
  --interval             update interval (ex. 5ms, 10s, 1m, 3h) (default: 10m0s)
  --national-rail-key    national rail api key (default: <none>)
  --once                 run once and exit (default: false)
  --pushover-app-key     pushover app key (default: <none>)
  --pushover-client-key  pushover client key (default: <none>)

Commands:

  version  Show the version information.
```
