package main

import (
	"log"

	"github.com/kainguyen/retail-store-api/internal/app/consumer"
)

func main() {
	app, err := consumer.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
