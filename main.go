package main

import (
	"log"

	"github.com/robinlant/mywiki/internal/wiki/internal/store"
	"github.com/robinlant/mywiki/internal/wiki/internal/web"
)

var addr = ":8800"

var mongoCon = "mongodb://127.0.0.1:27017"
var db = "test"

func main() {
	st, disc := store.NewMongoStore(mongoCon, db)
	defer disc()
	log.Printf("[INFO] Starting a web server at %s", addr)
	web.Run(st, addr)
}
