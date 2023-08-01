package dr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/cloud-native-toolkit/atkmod"
	"github.com/cloud-native-toolkit/itzcli/internal/prompt"
	"github.com/cloud-native-toolkit/itzcli/pkg"
	"github.com/google/uuid"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// DefaultGetter provides a function type for handling default values of the
// configuration.
type DefaultGetter func() interface{}

// NoDefault returns a default with no value.
func NoDefault() DefaultGetter {
	return func() interface{} {
		return nil
	}
}

// Prompter asks the user a question and returns the answer to the getter.
func Prompter(value interface{}) DefaultGetter {
	return func() interface{} {
		text := value.(string)
		key := uuid.New().String()
		builder := prompt.NewPromptBuilder()
		question, err := builder.Path(key).Text(text).Build()
		if err != nil {
			logger.Debugf("error when building question: %v", err)
			return "<replace me>"
		}
		err = prompt.Ask(question, os.Stdout, os.Stdin)
		if err != nil {
			logger.Debugf("error when asking question: %v", err)
			return "<replace me>"
		}
		return question.GetAnswer(key)
	}
}

// Messager asks the user a question and returns the answer to the getter.
func Messager(value interface{}) DefaultGetter {
	return func() interface{} {
		text := value.(string)
		b := bytes.NewBufferString(text + "\n")
		_, err := os.Stdout.Write(b.Bytes())
		if err != nil {
			logger.Debugf("error when writing to stdout: %v", err)
		}
		return "<replace me>"
	}
}

// Static returns the static value for the default.
func Static(value interface{}) DefaultGetter {
	return func() interface{} {
		return value
	}
}

func IifStatic(iif func() bool, tVal interface{}, fVal interface{}) DefaultGetter {
	return func() interface{} {
		if iif() {
			return tVal
		}
		return fVal
	}
}

// ConfigDir returns the static value for the default.
func ConfigDir(value interface{}) DefaultGetter {
	configDir, _ := pkg.GetITZHomeDir()
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

func getLocalIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

func getPodmanDefaultSystemIP() (net.IP, error) {
	cfg := &atkmod.CliParts{
		Path: viper.GetString("podman.path"),
		Cmd:  "system connection list --format \"{{.URI}}\"",
	}
	cmd := atkmod.NewPodmanCliCommandBuilder(cfg)
	stdOut := new(bytes.Buffer)
	stdErr := new(bytes.Buffer)
	localCtx := &atkmod.RunContext{
		Out: stdOut,
		Err: stdErr,
	}
	runner := &atkmod.CliModuleRunner{PodmanCliCommandBuilder: *cmd}
	err := runner.Run(localCtx)
	if err != nil {
		return nil, err
	}
	remoteURI := ""
	for _, line := range strings.Split(stdOut.String(), "\n") {
		if !strings.Contains(line, "localhost") {
			logger.Tracef("Found URL: %s", line)
			remoteURI = strings.Trim(strings.TrimSpace(line), "\"")
			break
		}
	}
	if remoteURI == "" {
		return nil, fmt.Errorf("unable to find non-localhost address")
	}
	uri, err := url.Parse(remoteURI)
	if err != nil {
		logger.Warnf("Could not parse URL: %s; %v", remoteURI, err)
		return nil, err
	}
	host := strings.Split(uri.Host, ":")
	if len(host) < 1 {
		return nil, fmt.Errorf("host in un-expected format: %s", uri.Host)
	}
	return net.ParseIP(host[0]), nil
}

// ServiceURL returns the static value for the default.
func ServiceURL(scheme string, port int) DefaultGetter {
	getIP := func() net.IP {
		ipAddr, err := getPodmanDefaultSystemIP()
		if err == nil {
			return ipAddr
		}
		// Fall back to the current IP address of the machine.
		ipAddr, err = getLocalIP()
		return ipAddr
	}
	return func() interface{} {
		theUrl := url.URL{
			Scheme: scheme,
			Host:   fmt.Sprintf("%s:%d", getIP(), port),
		}
		return theUrl.String()
	}
}

// Check interface for itz doctor checks
type Check interface {
	DoCheck(tryFix bool) (string, error)
}

type PreChecker func() bool
type ActionRunner func() (string, error)

// ActionCheck is
type ActionCheck struct {
	Message  string
	PreCheck PreChecker
	Cmd      ActionRunner
}

func (c *ActionCheck) DoCheck(tryFix bool) (string, error) {
	if c.PreCheck != nil && !c.PreCheck() {
		logger.Warnf("%s...  Skipped", c.Message)
		return "skipping action as precheck is false", nil
	}
	if c.Cmd != nil && tryFix {
		msg, err := c.Cmd()
		if err == nil {
			logger.Infof("%s...  OK", c.Message)
		}
		return msg, err
	}
	if c.Cmd == nil {
		return "", fmt.Errorf("no cmd runner")
	}
	return "", nil
}

func NewCmdActionCheck(msg string, preCheck PreChecker, cmd ActionRunner) Check {
	return &ActionCheck{
		Message:  msg,
		PreCheck: preCheck,
		Cmd:      cmd,
	}
}

type PodmanMachine struct {
	MachineState string
}

type podmanMachineOutput struct {
	Host PodmanMachine
}

func PodmanMachineExists() PreChecker {
	return func() bool {
		podmanPath, err := exec.LookPath("podman")
		if err != nil {
			return false
		}

		outputJson, err := exec.Command(podmanPath, "machine", "info", "--format", "json").Output()
		// unmarshell the json output to an object...
		var podmanInfo podmanMachineOutput
		err = json.Unmarshal(outputJson, &podmanInfo)
		if err != nil {
			return false
		}

		return strings.ToLower(podmanInfo.Host.MachineState) == "running"
	}
}

// UpdatePodmanMachineDate updates the date on the podman machine to be the
// current date on the host.
func UpdatePodmanMachineDate() ActionRunner {
	return func() (string, error) {
		// Get the current date formatted in 2023-04-26T14:45:26 format
		dateNow := time.Now().Format(time.RFC3339)
		_, err := exec.Command("podman", "machine", "ssh", "sudo", "date", "--set", dateNow).Output()
		if err != nil {
			return "", err
		}

		return "OK", nil
	}
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
			logger.Tracef("Trying to fix missing configuration key %s by setting to value: %s", c.ConfigKey, newValue)
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

type CheckerFunc func() (string, string, bool)

type FileAutoFixFunc func(path string) (string, error)

// FileCheck is a check for a required file
type FileCheck struct {
	PathCheckFunc CheckerFunc
	Path          string
	Name          string
	IsDir         bool
	Help          string
	FixerFunc     FileAutoFixFunc
	UpdaterFunc   FileAutoFixFunc
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
	foundPath := f.Path

	if f.PathCheckFunc != nil {
		if foundPath, name, found := f.PathCheckFunc(); found {
			f.Path = foundPath
			f.Name = name
		} else {
			return "", fmt.Errorf("%s not found", f.Path)
		}
	}

	for _, p := range strings.Split(f.Path, ":") {
		fn := filepath.Join(p, f.Name)
		if _, err := os.Stat(fn); errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does not exist
			continue
		}
		found = true
		foundPath = p
		logger.Infof("%s...  OK", f.Name)
		break
	}
	if !found {
		if tryFix && f.FixerFunc != nil {
			logger.Tracef("Could not find %s, attempting to fix.", f.Name)
			fixedPath, err := f.FixerFunc(filepath.Join(f.Path, f.Name))
			if err != nil {
				return "", err
			}
			logger.Infof("Did not find %s... Fixed", f.Name)
			return fixedPath, nil
		}
		logger.Warnf("%s not found", f.Name)
		return "", fmt.Errorf(f.Help, f.Name)
	} else {
		// The updater runs only if the file was found. This can be used to
		// save or record the file path, touch the file, update the file with
		// some other contents, etc.
		if f.UpdaterFunc != nil && len(foundPath) > 0 {
			_, err := f.UpdaterFunc(filepath.Join(foundPath, f.Name))
			if err != nil {
				return "", err
			}
		}
	}
	return filepath.Join(foundPath, f.Name), nil
}

// NewResourceFileCheck checks for any files, using the OS's PATH variable
// automatically as the path.
func NewResourceFileCheck(c CheckerFunc, help string, f FileAutoFixFunc) Check {
	return &FileCheck{
		PathCheckFunc: c,
		Path:          os.Getenv("PATH"),
		Name:          "",
		IsDir:         false,
		Help:          help,
		UpdaterFunc:   f,
	}
}

// ExistsOnPath checks if a binary file exists on the path, using the OS's PATH variable
// automatically as the path.
func ExistsOnPath(name string) CheckerFunc {
	return func() (string, string, bool) {
		foundPath, err := exec.LookPath(name)
		if err != nil {
			logger.Infof("%s...  Not found on Path", name)
			return "", "", false
		} else {
			foundPath = filepath.Dir(foundPath)
			name = filepath.Base(name)
			return foundPath, name, true
		}
	}
}

// OneExistsOnPath checks if one of the binary files in a list exists on the path, using the OS's PATH variable
// automatically as the path. It returns the first path that exists.
func OneExistsOnPath(names ...string) CheckerFunc {
	return func() (string, string, bool) {
		for _, name := range names {
			if foundPath, binName, found := ExistsOnPath(name)(); found == true {
				return foundPath, binName, found
			}
		}
		return "", "", false
	}
}

// NewReqConfigDirCheck checks for directories inside the ITZ home directory
func NewReqConfigDirCheck(name string) Check {
	dir, _ := pkg.GetITZHomeDir()
	return &FileCheck{
		Path:      dir,
		Name:      name,
		IsDir:     true,
		FixerFunc: CreateDir,
	}
}

// CreateDir creates a directory if it does not exist
func CreateDir(name string) (string, error) {
	return name, os.MkdirAll(name, os.ModePerm)
}

func EmptyFileCreator(path string) (string, error) {
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	_, err = f.WriteString("")
	if err != nil {
		return "", err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer f.Close()
	return f.Name(), nil
}

func UpdateConfig(configPath string) FileAutoFixFunc {
	return func(path string) (string, error) {
		logger.Tracef("Updating configuration <%s> with file path %s", configPath, path)
		viper.Set(configPath, path)
		err := viper.WriteConfig()
		return path, err
	}
}

func UpdateConfigIfMissing(configPath string) FileAutoFixFunc {
	return func(path string) (string, error) {
		existing := viper.GetString(configPath)
		if len(existing) > 0 {
			logger.Debugf("Configuration <%s> found; not updating", configPath)
			return path, nil
		}
		logger.Tracef("Updating configuration <%s> with file path %s", configPath, path)
		viper.Set(configPath, path)
		err := viper.WriteConfig()
		return path, err
	}
}

func TemplatedFileCreator(template string) FileAutoFixFunc {
	return func(path string) (string, error) {
		f, err := os.Create(path)
		if err != nil {
			return "", err
		}
		// And then write the template to the file
		_, err = f.WriteString(template)
		if err != nil {
			return "", err
		}
		//goland:noinspection GoUnhandledErrorResult
		defer f.Close()
		return f.Name(), nil
	}
}

// NewConfigFileCheck checks for files inside the ITZ home directory
func NewConfigFileCheck(name string) Check {
	dir, _ := pkg.GetITZHomeDir()
	return &FileCheck{
		Path:  dir,
		Name:  name,
		IsDir: false,
	}
}

func NewFixableConfigFileCheck(name string, fixFunc FileAutoFixFunc) Check {
	dir, _ := pkg.GetITZHomeDir()
	return &FileCheck{
		Path:      dir,
		Name:      name,
		IsDir:     false,
		FixerFunc: fixFunc,
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
