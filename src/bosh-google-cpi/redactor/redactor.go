package redactor

import (
	"regexp"
)

func RedactSecrets(sourceString string) string {
	re := regexp.MustCompile(`(?si)\\*"(account_key|json_key|password|private_key|secret_access_key)\\*"\\*: ?\\*".*?\\*"`)
	return re.ReplaceAllString(sourceString, `$1:"REDACTED"`)
}
