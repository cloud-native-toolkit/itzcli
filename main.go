package main

import (
	"github.com/cloud-native-toolkit/itzcli/cmd"
	"github.com/cloud-native-toolkit/itzcli/api"
)

var Version = "No Version Provided"

func main() {
	// Start the api server in the background  
	go api.StartServer()
	cmd.Execute(Version)
}
