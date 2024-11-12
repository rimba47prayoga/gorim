package middlewares

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now() // Start time recorded
		err := next(c)
		stop := time.Now()  // Stop time recorded

		// Calculate duration
		duration := stop.Sub(start).Seconds() * 1000

		// Get status and size
		status := c.Response().Status
		size := c.Response().Size

		// Log format similar to Django's, with duration
		logEntry := fmt.Sprintf("[%s] \"%s %s %s\" %d %d (Duration: %.2f ms)",
			stop.Format("02/Jan/2006 15:04:05"),
			c.Request().Method,
			c.Request().RequestURI,
			c.Request().Proto,
			status,
			size,
			duration,
		)

		// Print log entry
		fmt.Println(logEntry)
		return err
	}
}