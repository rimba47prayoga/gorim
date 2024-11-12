package api

import (
	"example/settings"
	"example/user"
)

func APIRoutes() {
	api := settings.Server.Group("/api/v1")
	user.RouterUser(api)
}
