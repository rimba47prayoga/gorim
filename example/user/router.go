package user

import (
	"github.com/labstack/echo/v4"
	"github.com/rimba47prayoga/gorim.git/routers"
)

func RouterUser(group *echo.Group) {
	routeGroup := group.Group("/user")
	userRoute := routers.NewDefaultRouter[*UserViewSet](routeGroup)
	userRoute.Register(NewUserViewSet)
}
