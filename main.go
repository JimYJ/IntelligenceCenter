package main

import (
	"IntelligenceCenter/router"
	"embed"
)

//go:embed static/dist
var static embed.FS

func main() {
	go router.Web(static)
	router.Api()
}
