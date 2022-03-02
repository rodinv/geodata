package repository

import (
	"bytes"
	"encoding/binary"
	"net"

	bin "github.com/ghostiam/binstruct"
)

func readString(r bin.Reader, n int) (string, error) {
	_, rawString, err := r.ReadBytes(n)
	if err != nil {
		return "", nil
	}

	return string(bytes.Trim(rawString, "\x00")), nil
}

func ip2int(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}
