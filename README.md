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

Create a Pushover app and grab your app key and client key from a device you
want to send messages to.

**Run it in daemon mode with twitter keys/tokens**

```console
# You need to either have environment variables that are
# PUSHOVER_APP_KEY and PUSHOVER_CLIENT_KEY
# or pass them into the container.
$ docker run -d --restart always \
    --name go-for-train \
    -e PUSHOVER_APP_KEY=foo
    -e PUSHOVER_CLIENT_KEY=bar
    w1lkins/go-for-train -d --interval 15m
```

## Usage

```console
go-for-train -  A bot that checks the status of my train journey and notifies me
about it.

Usage: go-for-train <command>

Flags:

  --app-key     pushover app key (default: foo)
  --client-key  pushover client key (default: bar)
  -d, --debug   enable debug logging (default: false)
  --interval    update interval (ex. 5ms, 10s, 1m, 3h) (default: 15m0s)
  --once        run once and exit (default: false)

Commands:

  version  Show the version information.
```
