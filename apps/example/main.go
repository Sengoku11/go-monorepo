// Example application that you can use as a playground
package main

import (
	"fmt"
	"os"

	"github.com/Sengoku11/go-monorepo/pkg/loadenv"
	"github.com/Sengoku11/go-monorepo/pkg/logger"
)

func main() {
	var log logger.Logger = logger.NewZerologLogger()

	envErr := loadenv.Local()
	if envErr != nil {
		msg := fmt.Sprintf("cannot start locally: %s", envErr)
		log.Fatal(msg)
	}

	environment := os.Getenv("ENVIRONMENT")
	msg := "environment variable ENVIRONMENT: " + environment
	log.Info(msg)
}
