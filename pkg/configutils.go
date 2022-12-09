package pkg

import (
	"fmt"
	"github.com/spf13/cobra"
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

// FlattenCommandName returns a string, such as "command.subcommand.step"
func FlattenCommandName(cmd *cobra.Command, suffix string) string {
	if len(suffix) > 0 {
		return fmt.Sprintf("%s.%s.%s", cmd.Parent().Name(), cmd.Name(), suffix)
	}
	return fmt.Sprintf("%s.%s", cmd.Parent().Name(), cmd.Name())
}
