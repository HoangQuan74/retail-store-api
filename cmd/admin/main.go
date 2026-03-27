package main

import (
	"log"

	"github.com/hoangquan/retail-store-api/internal/app/admin"
)

func main() {
	app, err := admin.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
