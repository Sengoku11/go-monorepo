// Example application that you can use as a playground
package main

import (
	"os"

	"github.com/Sengoku11/go-monorepo/pkg/loadenv"
	"github.com/Sengoku11/go-monorepo/pkg/logger"
)

func main() {
	if err := loadenv.Local(); err != nil {
		panic("cannot start locally: " + err.Error())
	}

	var log logger.Logger = logger.NewZerologLogger()

	log.Debug("this message appears when env DEBUG is true", "debug", os.Getenv("DEBUG"))
	log.Info("starting app")
}
