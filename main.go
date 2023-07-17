package main

import (
	"github.com/cloud-native-toolkit/itzcli/cmd"
)

// Version is used to print the version out to the user with the `itz version` command.
// This value is replaced in the LDFLAGS directive of the Makefile. See
// `LDFLAGS=-ldflags "-X main.Version=${ITZ_VER}"`
var Version = "No Version Provided"

func main() {
	cmd.Execute(Version)
}
