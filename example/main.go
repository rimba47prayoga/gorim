package main

import (
	"example/db"
	"example/user"
	"log"

	"github.com/rimba47prayoga/gorim.git"
	"github.com/rimba47prayoga/gorim.git/settings"
)


func main() {
	server := gorim.New()
	settings.Configure()
	db.SetupDatabase()

	api := server.Group("/api/v1")
	user.RouterUser(api)

	err := server.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
