// Example application that you can use as a playground
package main

import (
	"os"

	"github.com/Sengoku11/go-monorepo/pkg/loadenv"
	"github.com/Sengoku11/go-monorepo/pkg/logger"
)

func main() {
	var log logger.Logger = logger.NewZerologLogger()

	envErr := loadenv.Local()
	if envErr != nil {
		log.Fatal("cannot start locally", "error", envErr)
	}

	environment := os.Getenv("ENVIRONMENT")
	msg := "environment variable ENVIRONMENT: " + environment
	log.Info(msg)
}
