package main

import (
	"log"

	"github.com/kainguyen/retail-store-api/internal/app/api"
)

func main() {
	app, err := api.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
