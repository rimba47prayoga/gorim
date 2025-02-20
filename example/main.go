package main

import (
	"example/api"
	"example/settings"

	"gorim.org/gorim/cmd"
)


func main() {
	settings.Configure()
	api.APIRoutes()
	cmd.Execute()
}
