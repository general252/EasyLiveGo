package main

import (
	"github.com/general252/go-rtsp/server"
	"log"
)

func main() {
	app := server.DefaultApp
	app.Run()

	log.Println("app end")
}
