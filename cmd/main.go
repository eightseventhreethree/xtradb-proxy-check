package main

import (
	"log"

	api "gclustercheck/pkg/api"

	air "github.com/aofei/air"
)

func main() {
	log.Println("Hello Log")
	api.Init()
	air.Default.Serve()
}
