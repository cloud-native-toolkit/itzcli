package main

import (
	"fmt"
	"github.com/spf13/cobra/doc"
	"github.ibm.com/skol/atkcli/cmd"
)

func main() {
	err := doc.GenMarkdownTree(cmd.RootCmd, "./docs")
	if err != nil {
		fmt.Println(err)
	}
}
