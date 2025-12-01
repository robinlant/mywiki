package main

import (
	"fmt"

	"github.com/robinlant/mywiki/internal/wiki/internal/web"
)

const addr string = ":8800"

func main() {
	fmt.Println("Starting a web server at", addr)

	web.Run(addr)
}
