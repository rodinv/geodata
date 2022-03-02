# Geodata
REST API to get the coordinates of the user by his IP address and get a list of locations by the name of the city.

## Technical description:

The database is stored in the geobase.dat file, which is contained in the archive attached to the letter. The database will not change, it is read-only.
The database has a binary format. The file is sequentially stored:
```
60 byte header

int version; // version database
sbyte name[32]; // name/prefix for the database
ulong timestamp; // time of creation of the database
int records; // total number of records
uint offset_ranges; // offset relative to the beginning of the file to the beginning of the list of records with geoinformation
uint offset_cities; // offset relative to the beginning of the file to the beginning of the index sorted by city name
uint offset_locations; // offset relative to the beginning of the file to the beginning of the list of location records
```
```
12 bytes * Header.records (number of records) — a list of records with information about IP address intervals, sorted by ip_from and ip_to fields

uint ip_from; // start of IP address range
uint ip_to; // end of IP address range
uint location_index; // index of the location record
```
```
96 bytes * Header.records —number of records) - a list of records with location information with coordinates (longitude and latitude)

sbyte country[8]; // country name (a random string with the prefix "cou_")
sbyte region[12]; // name of the field (random string prefixed with "reg_")
postal sbyte[12]; // postal code (a random string with the prefix "pos_")
city sbyte[24]; // the name of the city (random string with the prefix "cit_")
sbyte organization[32]; // name of the company (a random string with the prefix "org_")
float latitude; // latitude
float longitude; // longitude
```
```
4 bytes * Header.records —number of records) - a list of indexes of location records sorted by city name, each index is the address of an entry in the file relative to Header.offset_locations
```

The database is loaded completely into memory when the application starts.

Two HTTP API methods must be implemented in the application:
GET /ip/location?ip=123.234.123.234
GET /city/locations?city=cit_Gbqw4

### install
```
make mod
make build
```
### run
```
bin/geodata --config=config/geodata.yaml
```

### documentation
[swagger file](docs/swagger.yaml)