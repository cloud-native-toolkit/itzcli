package pkg

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var replaceRegexp = regexp.MustCompile(`[^a-zA-Z0-9]`)

// GetITZHomeDir returns the home directory or the ITZ command
func GetITZHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", os.ErrNotExist
	}
	return filepath.Join(home, ".itz"), nil
}

func GetITZDir(dir string) (string, error) {
	return filepath.Join(MustITZHomeDir(), dir), nil
}

func GetITZDirOrDefault(dir string, envvar string) (string, error) {
	envval, exists := os.LookupEnv(envvar)
	if exists {
		if _, err := os.Stat(envval); !os.IsNotExist(err) {
			return envval, nil
		} else {
			return "", fmt.Errorf("%s does not exist as a dir", envvar)
		}
	}
	return GetITZDir(dir)
}

func GetITZCacheDir() (string, error) {
	return GetITZDirOrDefault("cache", "ITZ_CACHE_DIR")
}

func GetITZWorkDir() (string, error) {
	return GetITZDirOrDefault("workspace", "ITZ_WORKSPACE_DIR")
}

func MustITZHomeDir() string {
	home, err := GetITZHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return home
}

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
