package main

import (
	"log"

	"github.com/hoangquan/retail-store-api/services/consumer/internal"
)

func main() {
	app, err := internal.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
