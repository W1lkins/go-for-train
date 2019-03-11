package main

// Config represents the config for the service
type Config struct {
	Service struct {
		Code         string   `toml:"code"`
		Origins      []string `toml:"origins"`
		Destinations []string `toml:"destinations"`
	} `toml:"service"`
}

// OriginContains ensures that the value to check belongs to the Origin slice
func (c Config) OriginContains(check string) bool {
	for _, o := range c.Service.Origins {
		if o == check {
			return true
		}
	}
	return false
}

// DestinationContains ensures that the value to check belongs to the Destination slice
func (c Config) DestinationContains(check string) bool {
	for _, o := range c.Service.Destinations {
		if o == check {
			return true
		}
	}
	return false
}
