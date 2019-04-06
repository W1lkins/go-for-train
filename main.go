package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/W1lkins/go-for-train/version"
	"github.com/genuinetools/pkg/cli"
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
	// Endpoint to post notifications to
	notificationEndpoint string
	// Token required to post messages to endpoint
	notificationToken string
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

	p.FlagSet.StringVar(&notificationEndpoint, "notification-endpoint", os.Getenv("NOTIFICATION_ENDPOINT"), "notification endpoint")
	p.FlagSet.StringVar(&notificationToken, "notification-token", os.Getenv("NOTIFICATION_TOKEN"), "notification token")

	p.FlagSet.IntVar(&initialHour, "initial-hour", 15, "Send messages between these hours (lower bound)")
	p.FlagSet.IntVar(&finalHour, "final-hour", 17, "Send messages between these hours (upper bound)")

	p.Before = func(ctx context.Context) error {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}

	p.Action = func(ctx context.Context, args []string) error {
		if nRailAppKey == "" || notificationEndpoint == "" || notificationToken == "" {
			return fmt.Errorf("national-rail-key, notification-endpoint, and notificationToken are required")
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
	endpoint := notificationEndpoint + "?token=" + notificationToken
	client := NewClient()
	services := client.GetNextServices()
	logrus.Infof("Using config: %+v", config)
	logrus.Debugf("Using endpoint: %s", endpoint)

	for _, service := range services {
		logrus.Infof("Checking service at: %s", service.Scheduled)
		if service.Late || service.Cancelled {
			logrus.Infof("Service %s has an issue, checking whether we should send a notification", service.ID)
			logrus.Info(service)
			if !client.shouldNotify() {
				logrus.Info("Not sending message")
				continue
			}

			if client.shouldNotify() {
				logrus.Infof("Sending message about service at %s", service.Scheduled)
				status := ""
				extra := ""
				if service.Late {
					status = "late"
					extra = fmt.Sprintf("Expected at %s", service.Estimated)
				}
				if service.Cancelled {
					status = "cancelled"
					extra = fmt.Sprintf("Reason: %s", service.CancelledReason)
				}
				message := fmt.Sprintf("Train scheduled for %s is %s. %s", service.Scheduled, status, extra)
				title := fmt.Sprintf("Problem with service at %s", service.Scheduled)
				http.PostForm(endpoint, url.Values{"message": {message}, "title": {title}})
			}
		}
	}

	return nil
}
