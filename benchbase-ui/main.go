package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var port int

	flag.IntVar(&port, "p", 80, "Port to listen to.")

	flag.Parse()

	setupHandlers()

	log.Println("Now listening to port", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
