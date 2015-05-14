package util

import (
	"math"
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
