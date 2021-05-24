package main

import (
	"github.com/general252/go-rtsp/server"
	"log"
)

func main() {
	app := server.NewApp()
	app.Run()

	log.Println("app end")
}
