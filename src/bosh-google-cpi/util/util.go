package util

import (
	"math"
	"regexp"
	"strings"
)

var zoneRe *regexp.Regexp = regexp.MustCompile("/zones/([a-zA-Z1-9-]+)/")
var regionRe *regexp.Regexp = regexp.MustCompile("(^[a-zA-Z-]+[1-9])")

func ConvertMib2Gib(size int) int {
	sizeGb := float64(size) / float64(1024)
	return int(math.Ceil(sizeGb))
}

func ResourceSplitter(resource string) string {
	splits := strings.Split(resource, "/")

	return splits[len(splits)-1]
}

// RegionFromZone extracts the region from a zone.
// For example, us-central1-a produces us-central1
func RegionFromZone(zone string) string {
	s := regionRe.FindStringSubmatch(zone)
	if len(s) == 2 {
		return s[1]
	}
	return ""
}

// ZoneFromURL extracts and returns the zone from the fully-qualified
// URL of a zonal Google Compute Engine resource. The nil value is
// returned if a zone can not be found.
func ZoneFromURL(url string) string {
	s := zoneRe.FindStringSubmatch(url)
	if len(s) == 2 {
		return s[1]
	}
	return ""
}
