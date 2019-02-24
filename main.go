package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/genuinetools/pkg/cli"
	"github.com/sirupsen/logrus"
	"github.com/w1lkins/go-for-train/version"
)

var (
	// Whether or not to enable debug logging
	debug bool
	// How often to run
	interval time.Duration
	// Whether we should run once or not
	once bool
)

func main() {
	p := cli.NewProgram()
	p.Name = "go-for-train"
	p.Description = "A bot that checks the status of a train journey and alerts you if it is late or cancelled"
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	p.FlagSet = flag.NewFlagSet("go-for-train", flag.ExitOnError)
	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")
	p.FlagSet.BoolVar(&debug, "debug", false, "enable debug logging")

	p.FlagSet.BoolVar(&once, "once", true, "run once and exit")
	p.FlagSet.DurationVar(&interval, "interval", 15*time.Minute, "update interval (ex. 5ms, 10s, 1m, 3h)")

	p.Before = func(ctx context.Context) error {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}

	p.Action = func(ctx context.Context, args []string) error {
		ticker := time.NewTicker(interval)

		// On ^C, or SIGTERM handle exit.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)
		go func() {
			for sig := range c {
				logrus.Infof("Received %s, exiting.", sig.String())
				os.Exit(0)
			}
		}()

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
	client.CheckHomeRouteStatus()
	services := client.FindNextServices()
	for _, service := range services {
		logrus.Info(service)
		if !service.HasIssue {
			logrus.Infof("Service %s has an issue, checking now...", service.ID)
			client.CheckService(service)
		}
	}

	return nil
}
