package shared

import "net/url"

type WeatherStationEvent struct {
	OutsideTemperature string
	IndoorTemperature  string
	RequestFormData    url.Values
}
