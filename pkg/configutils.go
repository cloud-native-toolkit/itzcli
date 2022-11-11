package pkg

import (
	"regexp"
	"strings"
)

var replaceRegexp = regexp.MustCompile(`[^a-zA-Z0-9]`)

// Keyify returns a string value that is suitable for use as a YAML configuration
// key.
func Keyify(name string) string {
	// Just remove any non-alphanumeric characters
	return strings.ToLower(replaceRegexp.ReplaceAllString(name, ""))
}
