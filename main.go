package main

import (
	api "gclustercheck/pkg/api"
	"log"

	air "github.com/aofei/air"
)

func main() {
	log.Println("Started xtradb-proxy-check")
	api.Init()
	err := air.Default.Serve()
	if err != nil {
		log.Panicln("Failed to start!", err)
	}
}
