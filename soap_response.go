package main

import "encoding/xml"

// GetDepartureBoardResponseEnvelope contains the layout of the response from
// the GetDepartureBoard request
type GetDepartureBoardResponseEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soap    string   `xml:"soap,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Xsd     string   `xml:"xsd,attr"`
	Body    struct {
		Text                      string `xml:",chardata"`
		GetDepartureBoardResponse struct {
			Text                  string `xml:",chardata"`
			Xmlns                 string `xml:"xmlns,attr"`
			GetStationBoardResult struct {
				Text              string `xml:",chardata"`
				Lt                string `xml:"lt,attr"`
				Lt6               string `xml:"lt6,attr"`
				Lt7               string `xml:"lt7,attr"`
				Lt4               string `xml:"lt4,attr"`
				Lt5               string `xml:"lt5,attr"`
				Lt2               string `xml:"lt2,attr"`
				Lt3               string `xml:"lt3,attr"`
				GeneratedAt       string `xml:"generatedAt"`
				LocationName      string `xml:"locationName"`
				Crs               string `xml:"crs"`
				PlatformAvailable string `xml:"platformAvailable"`
				TrainServices     struct {
					Text    string `xml:",chardata"`
					Service []struct {
						Text         string `xml:",chardata"`
						Std          string `xml:"std"`
						Etd          string `xml:"etd"`
						Operator     string `xml:"operator"`
						OperatorCode string `xml:"operatorCode"`
						IsCancelled  string `xml:"isCancelled"`
						ServiceType  string `xml:"serviceType"`
						CancelReason string `xml:"cancelReason"`
						ServiceID    string `xml:"serviceID"`
						Origin       struct {
							Text     string `xml:",chardata"`
							Location struct {
								Text         string `xml:",chardata"`
								LocationName string `xml:"locationName"`
								Crs          string `xml:"crs"`
							} `xml:"location"`
						} `xml:"origin"`
						Destination struct {
							Text     string `xml:",chardata"`
							Location struct {
								Text         string `xml:",chardata"`
								LocationName string `xml:"locationName"`
								Crs          string `xml:"crs"`
							} `xml:"location"`
						} `xml:"destination"`
						Platform string `xml:"platform"`
						Length   string `xml:"length"`
						Rsid     string `xml:"rsid"`
					} `xml:"service"`
				} `xml:"trainServices"`
			} `xml:"GetStationBoardResult"`
		} `xml:"GetDepartureBoardResponse"`
	} `xml:"Body"`
}
