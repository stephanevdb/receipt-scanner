package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	// serve static files from the pb_public directory
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Static("/", "pb_public")
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
