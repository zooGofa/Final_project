package main

import (
	"final_project/pkg/db"
	"final_project/pkg/server"
	"log"
	"os"
)

func main() {
	// Allow switching web root via env if needed in future; default to ./web
	webDir := "./web"
	if v := os.Getenv("WEB_DIR"); v != "" {
		webDir = v
	}

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	if err := server.Run(webDir); err != nil {
		log.Fatal(err)
	}
}
