package repository

import "github.com/rodinv/errors"

var ErrNotFound = errors.New("not found")

type DB struct {
	Header          *Header    // database Header
	ipRanges        []*IpRange // ip ranges
	geoItems        []*GeoItem // locations info
	citySortedIndex []uint32   // sorted index by city name
}

type Header struct {
	Version         int32
	Name            string
	Timestamp       uint64
	RecordsCount    uint32
	offsetRanges    uint32
	offsetCities    uint32
	offsetLocations uint32
}

type IpRange struct {
	From        uint32
	To          uint32
	LocationPos uint32
}

type GeoItem struct {
	County       string  `json:"county,omitempty"`
	Region       string  `json:"region,omitempty"`
	Postal       string  `json:"postal,omitempty"`
	City         string  `json:"city,omitempty"`
	Organization string  `json:"organization,omitempty"`
	Latitude     float32 `json:"latitude,omitempty"`
	Longitude    float32 `json:"longitude,omitempty"`
}
