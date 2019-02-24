package main

import (
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gregdel/pushover"
	"github.com/sirupsen/logrus"
)

// Client is an HTTP client that can be used to make requests
type Client struct {
	http     *http.Client
	messager *Messager
}

// NewClient returns an instance of Client
func NewClient() Client {
	http := &http.Client{Timeout: 10 * time.Second}
	key, ok := os.LookupEnv("PUSHOVER_APP_KEY")
	if !ok {
		logrus.Fatal("PUSHOVER_APP_KEY not found")
	}
	pushover := pushover.New(key)

	nine := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 9, 0, 0, 0, time.UTC)
	five := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 17, 0, 0, 0, time.UTC)
	shouldSend := func() bool {
		return time.Now().After(nine) && time.Now().Before(five)
	}
	messager := &Messager{client: pushover, shouldSend: shouldSend}
	return Client{http, messager}
}

// CheckHomeRouteStatus checks the status of
// route-7-strathclyde_n-helensburgh_milngavie_edinburgh_bathgate
func (c Client) CheckHomeRouteStatus() bool {
	res, err := c.http.Get(StatusURL)
	if err != nil {
		logrus.Fatalf("Could not get status of routes: %v", err)
	}
	if res.StatusCode != 200 {
		logrus.Fatalf("Got %d response code from route status check", res.StatusCode)
	}
	defer res.Body.Close()

	var statusRes StatusResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&statusRes)
	if err != nil {
		logrus.Fatalf("Could not decode JSON: %v", err)
	}

	logrus.Debugf("%+v", statusRes)
	logrus.Infof("Got status for home route: %s", statusRes.toStatus())

	return statusRes.isOk()
}

// LiveURL is the URL of the live-boards endpoint
const liveURL = `https://www.scotrail.co.uk/nre/live-boards/EDP/lazy`

// FindNextServices will generate a list of the next services using Scotrails live-boards endpoint
func (c Client) FindNextServices() map[string]Service {
	res, err := c.http.Get(liveURL)
	if err != nil {
		logrus.Fatalf("Could not get live URL: %v", err)
	}
	if res.StatusCode != 200 {
		logrus.Fatalf("Got %d response code from live departure check", res.StatusCode)
	}
	defer res.Body.Close()

	// Parse HTML
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	services := make(map[string]Service, 0)
	doc.Find("tr.service").Each(func(i int, s *goquery.Selection) {
		care := true
		s.ChildrenFiltered("td").Each(func(i int, s *goquery.Selection) {
			if s.AttrOr("data-label", "") == "Destination" {
				if s.Text() != "Helensburgh Central" && s.Text() != "Milngavie" {
					care = false
				}
			}
		})
		if !care {
			return
		}

		var service = new(Service)
		ID := s.AttrOr("data-id", "unknown")
		service.ID = ID
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			switch s.AttrOr("data-label", "") {
			case "Arrives":
				// Our departure time is when it's due to arrive at EDP
				service.Departs = s.Text()
			case "Destination":
				service.Destination = s.Text()
			case "Expected":
				service.Status = s.Text()
				if s.Text() != "On time" {
					logrus.Debugf("Not 'On time'. Got '%s' setting HasIssue to true", s.Text())
					service.HasIssue = true
				}
			case "Origin":
				service.Origin = s.Text()
			}
		})
		service.CheckedAt = time.Now()
		services[ID] = *service
	})

	return services
}

const serviceURL = `https://www.scotrail.co.uk/nre/service-details/`

// CheckService will check the status of an individual service using the service-details endpoint
func (c Client) CheckService(sr Service) {
	res, err := c.http.Get(serviceURL + sr.ID)
	if err != nil {
		logrus.Fatalf("Could not get live URL: %v", err)
	}
	if res.StatusCode != 200 {
		logrus.Fatalf("Got %d response code from live departure check", res.StatusCode)
	}
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)
	doc.Find("ul").ChildrenFiltered("li").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "Edinburgh Park") {
			r := regexp.MustCompile("\\d\\d:\\d\\d")
			time := r.Find([]byte(s.Text()))
			logrus.Infof(
				"Service expected at Edinburgh Park: %s, actually arriving at: %s",
				sr.Departs,
				string(time),
			)
			key, ok := os.LookupEnv("PUSHOVER_CLIENT_KEY")
			if !ok {
				logrus.Warning("Would try to send notification but PUSHOVER_CLIENT_KEY not defined")
				return
			}

			if c.messager.shouldSend() {
				logrus.Debugf("Should send is %v, so sending message", c.messager.shouldSend())
				recipient := pushover.NewRecipient(key)
				message := pushover.NewMessageWithTitle("Actually arriving at "+string(time), "Issue with service at "+sr.Departs)
				c.messager.client.SendMessage(message, recipient)
			}
		}
	})
}
