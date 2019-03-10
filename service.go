package main

import (
	"fmt"
	"time"
)

// Service represents a single train from Origin --> Dest
type Service struct {
	ID              string
	Origin          string
	Destination     string
	Scheduled       string
	Estimated       string
	Late            bool
	Cancelled       bool
	CancelledReason string
	CheckedAt       time.Time
}

func (s Service) String() string {
	return fmt.Sprintf(
		"\n(%s) - Service [%s]\nFrom: %s\nTo: %s\nDeparting at: %s\nEstimated: %s\n Cancelled: %v\n CancelledReason: %s\n",
		s.CheckedAt,
		s.ID,
		s.Origin,
		s.Destination,
		s.Scheduled,
		s.Estimated,
		s.Cancelled,
		s.CancelledReason,
	)
}
