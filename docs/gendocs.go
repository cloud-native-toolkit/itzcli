package main

import (
	"fmt"
	"github.com/spf13/cobra/doc"
	"github.com/cloud-native-toolkit/itzcli/cmd"
)

func main() {
	cmd.RootCmd.DisableAutoGenTag = true
	err := doc.GenMarkdownTree(cmd.RootCmd, "./docs")
	if err != nil {
		fmt.Println(err)
	}
}
