// Package geoip geoip/geoip.go
package geoip

import (
	"log"
	"net"
)

// GetCountryByIP retrieves the country by IP using MaxMind
func GetCountryByIP(ip string) (string, error) {
	db, err := geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ipAddr := net.ParseIP(ip)
	record, err := db.Country(ipAddr)
	if err != nil {
		return "", err
	}

	return record.Country.IsoCode, nil
}
