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

	var inputFilename string
	var outputFilename string
	var port int
	var compress bool

	flag.BoolVar(&compress, "z", false, "Compress the database dump on disk.")
	flag.StringVar(&outputFilename, "o", "", "File to save the database to. Leave blank to use the input file.")
	flag.StringVar(&inputFilename, "i", "", "Input file. Leave blank to use the output file.")
	flag.IntVar(&port, "p", 2539, "Port to listen to.")
	flag.Parse()

	if outputFilename == "" {
		outputFilename = inputFilename
	}

	if inputFilename == "" {
		inputFilename = outputFilename
	}

	data := GetDatastore(inputFilename)

	// Makes the "push", "list" and "compare" HTTP handlers
	setupHandlers(data)

	// Save to disk every 5 minutes
	go func() {
		for _ = range time.Tick(5 * time.Minute) {
			data.SaveToDisk(outputFilename, compress)
		}
	}()

	// Also save to disk when pressing ^C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			log.Println("Saving database...")
			data.SaveToDisk(outputFilename, compress)
			log.Fatal("Exiting now.")
		}
	}()

	// And listen on the port
	log.Println("Now listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
