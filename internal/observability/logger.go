// Package observability is for logging the system wide observability
package observability

import (
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	Logger = log.New(os.Stdout,
		"",
		log.LstdFlags|log.LUTC|log.Lmicroseconds,
	)
}
