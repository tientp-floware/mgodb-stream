package util

import (
	"log"
	"regexp"

	"github.com/golang/geo/s2"
	"googlemaps.github.io/maps"
)

// RemoveSpecialChar example like this: #GoLangCode!$!  to GoLangCode
func RemoveSpecialChar(text string) string {
	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(text, "")
}

// PolylineFromPoint add point to encode
func PolylineFromPoint(points []maps.LatLng) *s2.Polyline {
	s2Points := []s2.LatLng{}
	for _, vs2point := range points {
		s2Points = append(s2Points, s2.LatLngFromDegrees(vs2point.Lat, vs2point.Lng))
	}
	return s2.PolylineFromLatLngs(s2Points)
}
