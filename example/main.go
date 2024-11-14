package main

import (
	"example/api"
	"example/settings"

	"github.com/rimba47prayoga/gorim.git/cmd"
)


func main() {	
	//
	settings.Configure()
	api.APIRoutes()
	cmd.Execute()
}
