# go-for-train 🚂

A dumb bot that checks trains I care about and notifies me if they're
potentially late/cancelled.

Uses the National Rail SOAP API and Gotify.

**Table of Contents**

<!-- toc -->

- [Installation](#installation)
    + [Via Go](#via-go)
    + [Running with Docker](#running-with-docker)
- [Usage](#usage)

<!-- tocstop -->

## Installation

#### Via Go

```console
$ go get -u -v github.com/evalexpr/go-for-train
```

#### Running with Docker

**Authentication**

You'll need to sign up to the National Rail API [here](http://realtime.nationalrail.co.uk/OpenLDBWSRegistration/) to get an API key.

And set up something to receive push notifications to an endpoint with a token. I use [gotify](https://gotify.net/)

**Run it in daemon mode**

```console
$ docker run -d --restart always \
    --name go-for-train \
    -e NATIONAL_RAIL_APP_KEY=foo \
    -e NOTIFICATION_ENDPOINT=https://foo.bar/message \
    -e NOTIFICATION_TOKEN=bar \
    -v /path/to/config.toml:/home/user/.config/go-for-train/config/config.toml \
    evalexpr/go-for-train -d --interval 15m
```

## Usage

```console
go-for-train -  A bot that checks the status of my train journey and notifies me about it.

Usage: go-for-train <command>

Flags:

  -d, --debug              enable debug logging (default: false)
  --final-hour             Send messages between these hours (upper bound) (default: 17)
  --initial-hour           Send messages between these hours (lower bound) (default: 15)
  --interval               update interval (ex. 5ms, 10s, 1m, 3h) (default: 10m0s)
  --national-rail-key      national rail api key (default: <none>)
  --notification-endpoint  notification endpoint (default: <none>)
  --notification-token     notification token (default: <none>)
  --once                   run once and exit (default: false)

Commands:

  version  Show the version information.
```
