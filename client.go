package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

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
	pushover := pushover.New(pushoverAppKey)
	nine := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), initialHour, 0, 0, 0, time.UTC)
	five := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), finalHour, 0, 0, 0, time.UTC)
	shouldSend := func() bool {
		now := time.Now()
		return now.Weekday() != 6 && now.Weekday() != 7 && now.After(nine) && now.Before(five)
	}
	messager := &Messager{client: pushover, shouldSend: shouldSend}
	return Client{http, messager}
}

// GetNextServices gets the next services for a certain station
// these results are then filtered by the rules in the config.toml
func (c Client) GetNextServices() []Service {
	payload := []byte(strings.TrimSpace(fmt.Sprintf(`
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope" xmlns:typ="http://thalesgroup.com/RTTI/2013-11-28/Token/types" xmlns:ldb="http://thalesgroup.com/RTTI/2017-10-01/ldb/">
   <soap:Header>
      <typ:AccessToken>
         <typ:TokenValue>%s</typ:TokenValue>
      </typ:AccessToken>
   </soap:Header>
   <soap:Body>
      <ldb:GetDepartureBoardRequest>
         <ldb:numRows>150</ldb:numRows>
         <ldb:crs>%s</ldb:crs>
         <ldb:filterCrs></ldb:filterCrs>
         <ldb:filterType>from</ldb:filterType>
         <ldb:timeOffset>0</ldb:timeOffset>
         <ldb:timeWindow>120</ldb:timeWindow>
      </ldb:GetDepartureBoardRequest>
   </soap:Body>
</soap:Envelope>
`, nRailAppKey, config.Service.Code,
	)))

	const soapURL = `https://lite.realtime.nationalrail.co.uk/OpenLDBWS/ldb11.asmx`
	res, err := c.http.Post(soapURL, "text/xml", bytes.NewReader(payload))
	if err != nil {
		logrus.Errorf("Could not complete SOAP request: %v", err)
		return make([]Service, 0)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Errorf("Could not read bytes from body: %v", err)
		return make([]Service, 0)
	}

	var statusRes GetDepartureBoardResponseEnvelope
	err = xml.Unmarshal(b, &statusRes)
	if err != nil {
		logrus.Errorf("Could not unmarshal XML: %v", err)
		return make([]Service, 0)
	}

	services := make([]Service, 0)
	for _, s := range statusRes.Body.GetDepartureBoardResponse.GetStationBoardResult.TrainServices.Service {
		origin := s.Origin.Location.LocationName
		destination := s.Destination.Location.LocationName
		if !config.OriginContains(origin) || !config.DestinationContains(destination) {
			continue
		}

		service := Service{
			ID:              s.ServiceID,
			Origin:          origin,
			Destination:     destination,
			Late:            s.Etd != "On time",
			Scheduled:       s.Std,
			Estimated:       s.Etd,
			Cancelled:       s.IsCancelled != "",
			CancelledReason: s.CancelReason,
			CheckedAt:       time.Now(),
		}
		services = append(services, service)
	}

	logrus.Debug(services)
	return services
}
