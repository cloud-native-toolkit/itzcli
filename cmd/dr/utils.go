package dr

import (
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// GetATKHomeDir returns the home directory or the ATK command
func GetATKHomeDir() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", os.ErrNotExist
	}
	return filepath.Join(home, ".atk"), nil
}

// DefaultGetter provides a function type for handling default values of the
// configuration.
type DefaultGetter func() interface{}

// NoDefault returns a default with no value.
func NoDefault() DefaultGetter {
	return func() interface{} {
		return nil
	}
}

// TODO: Use the prompt builder here for prompting the user for the value.
func Prompter(value interface{}) DefaultGetter {
	return func() interface{} {
		return "<replace me>"
	}
}

// Static returns the static value for the default.
func Static(value interface{}) DefaultGetter {
	return func() interface{} {
		return value
	}
}

// ConfigDir returns the static value for the default.
func ConfigDir(value interface{}) DefaultGetter {
	configDir, _ := GetATKHomeDir()
	return func() interface{} {
		return filepath.Join(configDir, value.(string))
	}
}

// RandomVal returns a random string value for the default.
func RandomVal(value interface{}) DefaultGetter {
	var rlen int
	if reflect.TypeOf(value).Kind() != reflect.Int {
		rlen = 8
	} else {
		rlen = value.(int)
	}
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	randFunc := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}
	return func() interface{} {
		return randFunc(rlen)
	}
}

// ServiceURL returns the static value for the default.
func ServiceURL(scheme string, port int) DefaultGetter {
	getIP := func() net.IP {
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		localAddr := conn.LocalAddr().(*net.UDPAddr)

		return localAddr.IP
	}
	return func() interface{} {
		theUrl := url.URL{
			Scheme: scheme,
			Host:   fmt.Sprintf("%s:%d", getIP(), port),
		}
		return theUrl.String()
	}
}

// Check interface for atk doctor checks
type Check interface {
	DoCheck(tryFix bool) (string, error)
}

// ConfigCheck a check for configuration
type ConfigCheck struct {
	ConfigKey string
	Defaulter DefaultGetter
	Help      string
}

// DoCheck performs a check of the configuration value and returns the value of
// the configuration, if it exists, with no error. If the configuration value
// does not exist
func (c *ConfigCheck) DoCheck(tryFix bool) (string, error) {
	cfg := viper.Get(c.ConfigKey)
	logger.Tracef("Found configuration key:value: %s:%s", c.ConfigKey, cfg)
	if cfg == nil {
		if !tryFix {
			logger.Warnf("Configuration key %s has no value", c.ConfigKey)
		}
		if tryFix && c.Defaulter != nil {
			newValue := c.Defaulter()
			logger.Tracef("Trying to fix missing coniguration key %s by setting to value: %s", c.ConfigKey, newValue)
			viper.Set(c.ConfigKey, newValue)
			logger.Infof("%s... Fixed", c.ConfigKey)
		} else {
			return "", fmt.Errorf("%s not found", c.ConfigKey)
		}
	} else {
		logger.Infof("%s... OK", c.ConfigKey)
	}
	return viper.GetString(c.ConfigKey), nil
}

// String provides a human-readable version of the config check
func (c *ConfigCheck) String() string {
	return fmt.Sprintf("configuration key: %s", c.ConfigKey)
}

// NewConfigCheck creates a new ConfigCheck for the given key.
func NewConfigCheck(configKey string, help string, defaulter DefaultGetter) Check {
	return &ConfigCheck{
		ConfigKey: configKey,
		Defaulter: defaulter,
		Help:      help,
	}
}

// FileCheck is a check for a required file
type FileCheck struct {
	Path      string
	Name      string
	IsDir     bool
	CanCreate bool
	Perms     int
	Help      string
}

// String provides for readable logging
func (f *FileCheck) String() string {
	return fmt.Sprintf("file: %s", f.Name)
}

// DoCheck performs a file check and returns the name of the file, if it exists,
// with no error or returns a nil string with an error if the file does not
// exist.
func (f *FileCheck) DoCheck(tryFix bool) (string, error) {
	found := false
	logger.Debugf("Using path: %v", f.Path)
	for _, p := range strings.Split(f.Path, ":") {
		fn := filepath.Join(p, f.Name)
		if _, err := os.Stat(fn); errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does not exist
			continue
		}
		found = true
		logger.Infof("%s...  OK", f.Name)
		break
	}
	if !found {
		logger.Errorf("Could not find file %s", f.Name)
		return "", fmt.Errorf(f.Help, f.Name)
	}
	return f.Name, nil
}

// NewBinaryFileCheck checks for binary files, using the OS's PATH variable
// automatically as the path.
func NewBinaryFileCheck(name string, help string) Check {
	return &FileCheck{
		Path:      os.Getenv("PATH"),
		Name:      name,
		IsDir:     false,
		CanCreate: false,
		Help:      help,
	}
}

// NewConfigDirCheck checks for directories inside the ATK home directory
func NewConfigDirCheck(name string) Check {
	dir, _ := GetATKHomeDir()
	return &FileCheck{
		Path:      dir,
		Name:      name,
		IsDir:     true,
		CanCreate: true,
		Perms:     0,
	}
}

// NewConfigFileCheck checks for files inside the ATK home directory
func NewConfigFileCheck(name string) Check {
	dir, _ := GetATKHomeDir()
	return &FileCheck{
		Path:      dir,
		Name:      name,
		IsDir:     false,
		CanCreate: false,
		Perms:     0,
	}
}

// DoChecks performs the configured file checks and return
// a list of errors, if any, while checking for the files. If fix is set to
// true, then the file or directory is created if it can be.
func DoChecks(checks []Check, tryFix bool) []error {
	logger.Infof("Performing %d checks...", len(checks))
	errs := make([]error, 0)
	for _, check := range checks {
		logger.Debugf("Checking %s", check)
		if _, err := check.DoCheck(tryFix); err != nil {
			errs = append(errs, err)
		}
	}
	if tryFix {
		logger.Trace("Writing configuration...")
		err := viper.WriteConfig()
		if err != nil {
			errs = append(errs, err)
		}
	}
	logger.Info("Done.")
	return errs
}
