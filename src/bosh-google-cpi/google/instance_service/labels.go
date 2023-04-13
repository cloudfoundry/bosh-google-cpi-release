package instance

import (
	"fmt"
	"regexp"
	"strings"
)

type Labels map[string]string

func (i *Labels) Validate() error {
	for k, v := range *i {
		if !mustMatchRe.MatchString(k) {
			return fmt.Errorf("Label key %q is invalid. Must batch regular expression %q", k, mustMatchReP)
		}
		if !mustMatchRe.MatchString(v) {
			return fmt.Errorf("Label value %q is invalid. Must batch regular expression %q", v, mustMatchReP)
		}
	}
	return nil
}

var (
	mustMatchReP = "^(?:[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?)$"
	mustMatchRe  = regexp.MustCompile(mustMatchReP)
)

// This function sanitizes an string, ensuring it is a valid label.
// It is used to clean up labels provided via BOSH metadata.
func SafeLabel(s string) (string, error) {
	maxlen := 61
	// Replace common invalid chars
	s = strings.Replace(s, "/", "-", -1)
	s = strings.Replace(s, "_", "-", -1)
	s = strings.Replace(s, ":", "-", -1)

	// Trim to max length
	if len(s) > maxlen {
		s = s[0:maxlen]
	}

	// Ensure the string doesn't begin or end in -
	s = strings.TrimSuffix(s, "-")
	s = strings.TrimPrefix(s, "-")

	// The sanitized value should pass the GCE regex
	if mustMatchRe.MatchString(s) {
		return s, nil
	}

	return "", fmt.Errorf("Label value %q did not satisfy the GCE label regexp", s)
}
