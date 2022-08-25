package outputters

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	logger "github.com/sirupsen/logrus"
	"github.ibm.com/Nathan-Good/atkdep"
	"github.ibm.com/Nathan-Good/atkmod"
	"gopkg.in/yaml.v2"
)

func WriteEntries(index []atkdep.IndexInfo, out io.Writer) error {
	for _, item := range index {
		io.WriteString(out, item.Id)
		io.WriteString(out, "\n")
	}
	return nil
}

func WriteModuleInfo(info *atkmod.ModuleInfo, out io.Writer) error {
	if info != nil {
		io.WriteString(out, fmt.Sprintf("Name: %s\n", info.Id))
		io.WriteString(out, fmt.Sprintf("Version: %s\n", info.Version))
		io.WriteString(out, fmt.Sprintf("Description: %s\n", info.Name))
		io.WriteString(out, fmt.Sprintf("Dependencies: %s\n", info.Dependencies))
	}
	return nil
}

func WriteToFile(info *atkdep.IndexInfo, toFile string) error {
	// If the file already exists, we load the file up
	var entries []atkdep.IndexInfo = []atkdep.IndexInfo{}
	if _, err := os.Stat(toFile); err == nil {
		logger.Info("adding module %s to existing file at %s", info.Id, toFile)
		yamlFile, err := ioutil.ReadFile(toFile)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(yamlFile, &entries)

		if err != nil {
			return err
		}
	} else if errors.Is(err, os.ErrNotExist) {
		// just write the single info to the file...
		entries = append(entries, *info)
	}

	data, err := yaml.Marshal(&entries)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(toFile, data, 0644)

	return err
}
