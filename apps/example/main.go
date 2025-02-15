// Example application that you can use as a playground
package main

import (
	"log"
	"os"

	"github.com/Sengoku11/go-monorepo/pkg/loadenv"
)

func main() {
	envErr := loadenv.Local()
	if envErr != nil {
		log.Fatal("cannot start locally: ", envErr)
	}

	environment := os.Getenv("ENVIRONMENT")
	log.Printf("environment variable ENVIRONMENT: %s", environment)
}
