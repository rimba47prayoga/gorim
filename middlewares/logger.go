package middlewares

import (
	"fmt"
	"time"

	"github.com/rimba47prayoga/gorim.git"
)

type LoggerMiddleware struct {
	BaseMiddleware
}

func (m *LoggerMiddleware) Call(c gorim.Context) error {
	start := time.Now() // Start time recorded
	err := m.Next(c)
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