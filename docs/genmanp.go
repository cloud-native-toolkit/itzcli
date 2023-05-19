package main

import (
	"fmt"
	"github.com/cloud-native-toolkit/itzcli/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	cmd.RootCmd.DisableAutoGenTag = true
	mh := &doc.GenManHeader{
		Title:   "ITZ",
		Section: "1",
	}
	err := doc.GenManTree(cmd.RootCmd, mh, "./contrib/manpages")
	if err != nil {
		fmt.Println(err)
	}
}
