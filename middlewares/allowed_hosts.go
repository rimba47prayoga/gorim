package middlewares

import (
	"net"
	"net/http"

	"gorim.org/gorim"
	"gorim.org/gorim/conf"
)

type AllowedHostsMiddleware struct {
	BaseMiddleware
}

func (m *AllowedHostsMiddleware) Call(c gorim.Context) error {
	host := c.Request().Host
	host, _, err := net.SplitHostPort(host)
	if err != nil {
		host = c.Request().Host
	}
	for _, allowedHost := range conf.ALLOWED_HOSTS {
		if host == allowedHost {
			return m.Next(c)
		}
	}
	// Reject requests from disallowed hosts
	return c.JSON(http.StatusForbidden, map[string]string{
		"error": "Host not allowed",
	})
}
