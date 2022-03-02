package repository

import (
	"encoding/binary"
	"os"

	bin "github.com/ghostiam/binstruct"
	"github.com/rodinv/errors"
)

func Load(path string) (*DB, error) {
	// Opening file
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "opening file")
	}
	defer f.Close()

	reader := bin.NewReader(f, binary.LittleEndian, false)

	// Parse Header
	header, err := parseHeader(reader)
	if err != nil {
		return nil, errors.Wrap(err, "reading header")
	}

	// Parse city index
	cityIndex, err := parseCityIndex(reader, header.offsetCities, header.RecordsCount)
	if err != nil {
		return nil, errors.Wrap(err, "reading city indexes")
	}

	// Parse IP ranges
	ranges, err := parseIpRanges(reader, header.offsetRanges, header.RecordsCount)
	if err != nil {
		return nil, errors.Wrap(err, "reading ip ranges")
	}

	// Parse locations
	geoItems, err := parseGeoItems(reader, header.offsetLocations, header.RecordsCount)
	if err != nil {
		return nil, errors.Wrap(err, "reading locations")
	}

	return &DB{
		Header:          header,
		ipRanges:        ranges,
		geoItems:        geoItems,
		citySortedIndex: cityIndex,
	}, nil
}

func parseHeader(r bin.Reader) (*Header, error) {
	var err error

	header := new(Header)

	header.Version, err = r.ReadInt32()
	if err != nil {
		return nil, errors.Wrap(err, "reading version")
	}

	header.Name, err = readString(r, 32)
	if err != nil {
		return nil, errors.Wrap(err, "reading name")
	}

	header.Timestamp, err = r.ReadUint64()
	if err != nil {
		return nil, errors.Wrap(err, "reading timestamp")
	}

	header.RecordsCount, err = r.ReadUint32()
	if err != nil {
		return nil, errors.Wrap(err, "reading records")
	}

	header.offsetRanges, err = r.ReadUint32()
	if err != nil {
		return nil, errors.Wrap(err, "reading offset_ranges")
	}

	header.offsetCities, err = r.ReadUint32()
	if err != nil {
		return nil, errors.Wrap(err, "reading offset_cities")
	}

	header.offsetLocations, err = r.ReadUint32()
	if err != nil {
		return nil, errors.Wrap(err, "reading offset_locations")
	}

	return header, nil
}

func parseIpRanges(r bin.Reader, offset uint32, count uint32) (ranges []*IpRange, err error) {
	_, err = r.Seek(int64(offset), 0)
	if err != nil {
		return nil, errors.Wrap(err, "changing offset")
	}

	ranges = make([]*IpRange, 0, count)

	for i := 0; i < int(count); i++ {
		row := new(IpRange)

		row.From, err = r.ReadUint32()
		if err != nil {
			return nil, errors.Wrapf(err, "row %d, ip_from", i)
		}

		row.To, err = r.ReadUint32()
		if err != nil {
			return nil, errors.Wrapf(err, "row %d, ip_to", i)
		}

		row.LocationPos, err = r.ReadUint32()
		if err != nil {
			return nil, errors.Wrapf(err, "row %d, location", i)
		}

		ranges = append(ranges, row)
	}

	return
}

func parseCityIndex(r bin.Reader, offset uint32, count uint32) ([]uint32, error) {
	const locationRecordSize = 96

	_, err := r.Seek(int64(offset), 0)
	if err != nil {
		return nil, errors.Wrap(err, "changing offset")
	}

	result := make([]uint32, 0, count)
	for i := 0; i < int(count); i++ {
		val, err := r.ReadUint32()
		if err != nil {
			return nil, errors.Wrapf(err, "row %d", i)
		}

		result = append(result, val/locationRecordSize)
	}

	return result, nil
}

func parseGeoItems(r bin.Reader, offset uint32, count uint32) (geoItems []*GeoItem, err error) {
	_, err = r.Seek(int64(offset), 0)
	if err != nil {
		return nil, errors.Wrap(err, "changing offset")
	}

	geoItems = make([]*GeoItem, 0, count)

	for i := 0; i < int(count); i++ {
		row := new(GeoItem)

		row.County, err = readString(r, 8)
		if err != nil {
			return nil, errors.Wrapf(err, "row %d: reading country", i)
		}

		row.Region, err = readString(r, 12)
		if err != nil {
			return nil, errors.Wrapf(err, "row %d: reading region", i)
		}

		row.Postal, err = readString(r, 12)
		if err != nil {
			return nil, errors.Wrapf(err, "row %d: reading postal", i)
		}

		row.City, err = readString(r, 24)
		if err != nil {
			return nil, errors.Wrapf(err, "row %d: reading city", i)
		}

		row.Organization, err = readString(r, 32)
		if err != nil {
			return nil, errors.Wrapf(err, "row %d: reading city", i)
		}

		row.Latitude, err = r.ReadFloat32()
		if err != nil {
			return nil, errors.Wrapf(err, "row %d: reading latitude", i)
		}

		row.Longitude, err = r.ReadFloat32()
		if err != nil {
			return nil, errors.Wrapf(err, "row %d: reading longitude", i)
		}

		geoItems = append(geoItems, row)
	}

	return geoItems, nil
}
