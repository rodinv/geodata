package repository

import (
	"sort"
)

// GetLocationByIP gets location items by ip address string
func (d *DB) GetLocationByIP(ip string) (*GeoItem, error) {
	ipRaw := ip2int(ip)

	// binary search ip ranges
	indx := sort.Search(len(d.ipRanges), func(i int) bool {
		return d.ipRanges[i].From >= ipRaw || d.ipRanges[i].To >= ipRaw
	})

	// if found
	if indx < len(d.ipRanges) &&
		d.ipRanges[indx].From <= ipRaw && d.ipRanges[indx].To >= ipRaw {
		return d.geoItems[int(d.ipRanges[indx].LocationPos)], nil
	}

	return nil, ErrNotFound
}

// GetLocationsByCity gets list of locations by city name
func (d *DB) GetLocationsByCity(city string) ([]*GeoItem, error) {
	// find first item in sorted index (binary search)
	indx := sort.Search(len(d.citySortedIndex), func(i int) bool {
		item := d.geoItems[d.citySortedIndex[i]]

		return item.City >= city
	})

	// if not found
	if indx >= len(d.citySortedIndex) || indx < 0 || d.geoItems[d.citySortedIndex[indx]].City != city {
		return nil, ErrNotFound
	}

	// get items by sorted index
	items := make([]*GeoItem, 0)
	for i := indx; ; i++ {
		if d.geoItems[d.citySortedIndex[i]].City != city {
			break
		}

		items = append(items, d.geoItems[d.citySortedIndex[i]])
	}

	return items, nil
}
