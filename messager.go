package main

import "github.com/gregdel/pushover"

// Messager defines a way tosend messages
type Messager struct {
	// TODO(jwilkins): Don't rely on pushover
	client     *pushover.Pushover
	shouldSend func() bool
}
