package main

// StatusURL is the interactive_map/status endpoint from Scotrail
const StatusURL = `https://www.scotrail.co.uk/ajax/interactive_map/status`

// StatusResponse parses the response from Scortail's interactive_map/status endpoint
type StatusResponse struct {
	Routes struct {
		HomeRoute struct {
			Map    string `json:"map"`
			Status string `json:"status"`
		} `json:"route-7-strathclyde_n-helensburgh_milngavie_edinburgh_bathgate"`
	} `json:"routes"`
	Timestamp int `json:"timestamp"`
}

func (s StatusResponse) isOk() bool {
	return s.toStatus() == "good"
}

func (s StatusResponse) toStatus() string {
	return s.Routes.HomeRoute.Status
}
