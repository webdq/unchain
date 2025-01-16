package client

import (
	_ "embed"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log/slog"
	"net"
)

//go:embed GeoLite2-Country.mmdb
var geoData []byte

var geoDB *geoip2.Reader

func init() {
	db, err := geoip2.FromBytes(geoData)
	if err != nil {
		slog.Error("failed to load geo database", "err", err)
	} else {
		geoDB = db
	}
}

func GeoDbClose() {
	if geoDB != nil {
		geoDB.Close()
	}
}

func Country(hostOrIp string) (isoCountryCode string, err error) {
	if geoDB == nil {
		return "", fmt.Errorf("geo databse is nil")
	}
	ip := net.ParseIP(hostOrIp)
	if ip == nil {
		ips, err := net.LookupIP(hostOrIp)
		if err != nil {
			return "", fmt.Errorf("failed to lookup IP: %w", err)
		}
		if len(ips) == 0 {
			return "", fmt.Errorf("no IP found for %s", hostOrIp)
		}
		ip = ips[0]
	}
	record, err := geoDB.Country(ip)
	if err != nil {
		return "", fmt.Errorf("failed to get Country: %w", err)
	}
	return record.Country.IsoCode, nil
}
