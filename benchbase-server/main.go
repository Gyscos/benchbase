// BenchBase is a benchmark database.
// It offers a simple HTTP API to push new benchmark results, and retrieve / sort / filter existing ones.

// Actions:
// Store: Benchmark
// Compare: (A, B ConfigurationFilter, SituationFilter) [][2]Benchmark
// List: (A ConfigurationFilter, SituationFilter) []Benchmark
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	var dataFilename string
	var port int

	flag.StringVar(&dataFilename, "d", "data.json", "File to save and load the database.")
	flag.IntVar(&port, "p", 6666, "Port to listen to.")
	flag.Parse()

	data := GetDatastore(dataFilename)

	// Makes the "push", "list" and "compare" HTTP handlers
	setupHandlers(data)

	// Save to disk every 5 minutes
	go func() {
		for _ = range time.Tick(5 * time.Minute) {
			data.SaveToDisk(dataFilename)
		}
	}()

	// Also save to disk when pressing ^C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			log.Println("Saving database...")
			data.SaveToDisk(dataFilename)
			log.Fatal("Exiting now.")
		}
	}()

	// And listen on the port
	log.Println("Now listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
