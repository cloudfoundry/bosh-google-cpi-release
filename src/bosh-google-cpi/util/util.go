package util

import (
	"math"
	"regexp"
	"strings"
)

func ConvertMib2Gib(size int) int {
	sizeGb := float64(size) / float64(1024)
	return int(math.Ceil(sizeGb))
}

func ResourceSplitter(resource string) string {
	splits := strings.Split(resource, "/")

	return splits[len(splits)-1]
}

var regionRe = regexp.MustCompile("(^[a-zA-Z-]+[1-9])")

// RegionFromZone extracts the region from a zone.
// For example, us-central1-a produces us-central1
func RegionFromZone(zone string) string {
	s := regionRe.FindStringSubmatch(zone)
	if len(s) == 2 {
		return s[1]
	}
	return ""
}

var zoneRe = regexp.MustCompile("/zones/([a-zA-Z1-9-]+)/?")

// ZoneFromURL extracts and returns the zone from the fully-qualified
// URL of a zonal Google Compute Engine resource. The zero value is
// returned if a zone can not be found.
func ZoneFromURL(url string) string {
	s := zoneRe.FindStringSubmatch(url)
	if len(s) == 2 {
		return s[1]
	}
	return ""
}

var regionURLRe = regexp.MustCompile("/regions/([a-zA-Z1-9-]+)/?")

// RegionFromURL extracts and returns the region from the fully-qualified
// URL of a regional Google Compute Engine resource. The zero value
// is returned if a region can not be parsed.
func RegionFromURL(url string) string {
	s := regionURLRe.FindStringSubmatch(url)
	if len(s) == 2 {
		return s[1]
	}
	return ""
}
