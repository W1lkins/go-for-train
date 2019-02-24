package main

import (
	"fmt"
	"time"
)

// Service represents a single train from Origin --> Dest
type Service struct {
	ID          string
	Origin      string
	Destination string
	Departs     string
	Status      string
	CheckedAt   time.Time
	HasIssue    bool
}

func (s Service) String() string {
	return fmt.Sprintf(
		"Service [%s] (%s)\n\tDestination: %s\n\tDeparture time: %s\n\tStatus: %s\n\tIssue: %v",
		s.ID,
		s.CheckedAt.Format("3:04PM"),
		s.Destination,
		s.Departs,
		s.Status,
		s.HasIssue,
	)
}
