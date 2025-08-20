package instance

import (
	"fmt"
	"regexp"
	"strings"
)

type Labels map[string]string

func (i *Labels) Validate() error {
	for k, v := range *i {
		if !mustMatchReKey.MatchString(k) {
			return fmt.Errorf("Label key %q is invalid. Must batch regular expression %q", k, mustMatchReKeyP) //nolint:staticcheck
		}
		if !mustMatchReValue.MatchString(v) {
			return fmt.Errorf("Label value %q is invalid. Must batch regular expression %q", v, mustMatchReValueP) //nolint:staticcheck
		}
	}
	return nil
}

var (
	numFirstRe        = regexp.MustCompile("^[0-9]")
	mustMatchReKeyP   = "^(?:[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?)$"
	mustMatchReValueP = "^(?:[a-z0-9](?:[-_a-z0-9]{0,61}[a-z0-9])?)$"
	mustMatchReKey    = regexp.MustCompile(mustMatchReKeyP)
	mustMatchReValue  = regexp.MustCompile(mustMatchReValueP)
)

// This function sanitizes a string, ensuring it is a valid label.
// It is used to clean up labels provided via BOSH metadata.

func SafeLabel(s string, isKey bool) (string, error) {
	maxlen := 61
	// Replace common invalid chars
	s = strings.Replace(s, "/", "-", -1) //nolint:staticcheck
	s = strings.Replace(s, "_", "-", -1) //nolint:staticcheck
	s = strings.Replace(s, ":", "-", -1) //nolint:staticcheck

	// Trim to max length
	if len(s) > maxlen {
		s = s[0:maxlen]
	}

	// Ensure the string doesn't begin or end in -
	s = strings.TrimSuffix(s, "-")
	s = strings.TrimPrefix(s, "-")

	// Ensure the string doesn't begin with a number
	if isKey && numFirstRe.MatchString(s) {
		s = "n" + s
	}

	// The sanitized value should pass the GCE regex
	if isKey {
		if mustMatchReKey.MatchString(s) {
			return s, nil
		} else {
			return "", fmt.Errorf("Label value %q did not satisfy the GCE label regexp", s) //nolint:staticcheck
		}
	} else if mustMatchReValue.MatchString(s) {
		return s, nil
	}

	return "", fmt.Errorf("Label value %q did not satisfy the GCE label regexp", s) //nolint:staticcheck
}
