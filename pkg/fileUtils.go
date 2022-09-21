package pkg

import (
	"os"
	"regexp"
)

var re = regexp.MustCompile(`^file:\/{0,2}`)

// ReadFile reads the given file into the byte array
func ReadFile(path string) ([]byte, error) {
	// Yes, for now this is just a wrapper around os.ReadFile
	p := re.ReplaceAllString(path, ``)
	return os.ReadFile(p)
}
