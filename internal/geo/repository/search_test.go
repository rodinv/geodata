package repository

import (
	"reflect"
	"testing"
)

const testFilePath = "../../../data/geobase.dat"

func Test_GetLocationByIP(t *testing.T) {
	d, err := Load(testFilePath)
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name string
		ip   string
	}{
		{"test_1", "0.0.0.0"},
		{"test_2", "255.255.255.255"},
		{"test_3", "123.234.123.234"},
		{"test_4", "55.33.66.77"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want, wantErr := getLocationByIpLinear(d, tt.ip)

			got, err := d.GetLocationByIP(tt.ip)
			if err != wantErr {
				t.Errorf("GetLocationByIP() error = %v, wantErr %v", err, wantErr)
				return
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("GetLocationByIP() got = %v, want %v", got, want)
			}
		})
	}
}

func getLocationByIpLinear(d *DB, ip string) (*GeoItem, error) {
	ipRaw := ip2int(ip)
	for _, v := range d.ipRanges {
		if v.From <= ipRaw && v.To >= ipRaw {
			return d.geoItems[int(v.LocationPos)], nil
		}
	}

	return nil, ErrNotFound
}

func Test_GetLocationsByCity(t *testing.T) {
	d, err := Load(testFilePath)
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name string
		city string
	}{
		{"test_1", "cit_Ejid"},
		{"test_2", "cit_Aqis Ve"},
		{"test_3", "cit_Uwagycyguwyn"},
		{"test_4", "cit_Ululof Cirag Ce"},
		{"test_5", ""},
		{"test_6", "dummy"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want, wantErr := getLocationsByCityLinear(d, tt.city)

			got, err := d.GetLocationsByCity(tt.city)
			if err != wantErr {
				t.Errorf("GetLocationsByCity() error = %v, wantErr %v", err, wantErr)
				return
			}

			if len(got) != len(want) { // TODO check with sorted slices
				t.Errorf("GetLocationsByCity() got = %+v, want %+v", got, want)
			}
		})
	}
}

func getLocationsByCityLinear(d *DB, city string) ([]*GeoItem, error) {
	items := make([]*GeoItem, 0)
	for _, v := range d.geoItems {
		if v.City == city {
			items = append(items, v)
		}
	}

	if len(items) == 0 {
		return nil, ErrNotFound
	}

	return items, nil
}
