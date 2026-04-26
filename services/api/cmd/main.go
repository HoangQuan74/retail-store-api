package main

import (
	"log"

	"github.com/hoangquan/retail-store-api/services/api/internal"
)

//	@title			Retail Store API
//	@version		1.0
//	@description	REST API for managing retail store products and categories.

//	@host		localhost:8080
//	@BasePath	/api/v1

func main() {
	app, err := internal.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
