package main

import (
	"log"

	"github.com/hoangquan/retail-store-api/internal/app/socket"
)

func main() {
	app, err := socket.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
