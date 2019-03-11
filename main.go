package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/W1lkins/go-for-train/version"
	"github.com/genuinetools/pkg/cli"
	"github.com/gregdel/pushover"
	"github.com/shibukawa/configdir"
	"github.com/sirupsen/logrus"
)

var (
	// Whether or not to enable debug logging
	debug bool
	// How often to run
	interval time.Duration
	// Whether we should run once or not
	once bool
	// NationalRail API key from (http://lite.realtime.nationalrail.co.uk/openldbws/)
	nRailAppKey string
	// Pushover app key
	pushoverAppKey string
	// Pushover client key
	pushoverClientKey string
	// Send messages between this hour (lower bound)
	initialHour int
	// Send messages between this hour (upper bound)
	finalHour int
	// Config for the service
	config Config
)

func main() {
	p := cli.NewProgram()
	p.Name = "go-for-train"
	p.Description = "A bot that checks the status of my train journey and notifies me about it."
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	p.FlagSet = flag.NewFlagSet("go-for-train", flag.ExitOnError)
	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")
	p.FlagSet.BoolVar(&debug, "debug", false, "enable debug logging")

	p.FlagSet.BoolVar(&once, "once", false, "run once and exit")
	p.FlagSet.DurationVar(&interval, "interval", 10*time.Minute, "update interval (ex. 5ms, 10s, 1m, 3h)")

	p.FlagSet.StringVar(&nRailAppKey, "national-rail-key", os.Getenv("NATIONAL_RAIL_APP_KEY"), "national rail api key")

	p.FlagSet.StringVar(&pushoverAppKey, "pushover-app-key", os.Getenv("PUSHOVER_APP_KEY"), "pushover app key")
	p.FlagSet.StringVar(&pushoverClientKey, "pushover-client-key", os.Getenv("PUSHOVER_CLIENT_KEY"), "pushover client key")

	p.FlagSet.IntVar(&initialHour, "initial-hour", 15, "Send messages between these hours (lower bound)")
	p.FlagSet.IntVar(&finalHour, "final-hour", 17, "Send messages between these hours (upper bound)")

	p.Before = func(ctx context.Context) error {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}

	p.Action = func(ctx context.Context, args []string) error {
		if nRailAppKey == "" || pushoverAppKey == "" || pushoverClientKey == "" {
			return fmt.Errorf("national-rail-key, pushover-app-key, and pushover-client-key are required")
		}

		ticker := time.NewTicker(interval)

		// On ^C, or SIGTERM handle exit.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)
		go func() {
			for sig := range c {
				logrus.Infof("received %s, exiting.", sig.String())
				os.Exit(0)
			}
		}()

		configDir := configdir.New("go-for-train", "config")
		conf := configDir.QueryFolderContainsFile("config.toml")
		if conf == nil {
			logrus.Fatalf("no config.toml found in folder: %s", configDir.LocalPath)
		}

		data, _ := conf.ReadFile("config.toml")
		_, err := toml.Decode(string(data), &config)
		if err != nil {
			logrus.Fatalf("could not decode toml data: %v", err)
		}

		if once {
			run()
			os.Exit(0)
		}

		logrus.Infof("starting bot to update every %s", interval)
		for range ticker.C {
			run()
		}

		return nil
	}

	p.Run()
}

func run() error {
	client := NewClient()
	services := client.GetNextServices()
	logrus.Infof("Using config: %+v", config)

	for _, service := range services {
		logrus.Infof("Checking service at: %s", service.Scheduled)
		if service.Late || service.Cancelled {
			logrus.Infof("Service %s has an issue, checking whether we should send a notification", service.ID)
			logrus.Info(service)
			if !client.messager.shouldSend() {
				logrus.Info("Not sending message")
				continue
			}

			if client.messager.shouldSend() {
				logrus.Infof("Sending message about service at %s", service.Scheduled)
				recipient := pushover.NewRecipient(pushoverClientKey)
				message := pushover.NewMessageWithTitle(
					fmt.Sprintf("Late: %v\nScheduled: %s\nExpected: %s\nCancelled: %v\nReason: %s", service.Late, service.Scheduled, service.Estimated, service.Cancelled, service.CancelledReason),
					fmt.Sprintf("Problem with service at %s", service.Scheduled),
				)
				client.messager.client.SendMessage(message, recipient)
			}
		}
	}

	return nil
}
