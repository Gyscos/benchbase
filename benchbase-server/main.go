// BenchBase is a benchmark database.
// It offers a simple HTTP API to push new benchmark results, and retrieve / sort / filter existing ones.

// Actions:
// Store: Benchmark
// Compare: (A, B ConfigurationFilter, SituationFilter) [][2]Benchmark
// List: (A ConfigurationFilter, SituationFilter) []Benchmark
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Gyscos/benchbase"
)

func main() {
	var dataFilename = "data.json"
	data := GetDatastore(dataFilename)

	http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		var benchmark benchbase.Benchmark

		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&benchmark)
		if err != nil {
			log.Println("Error reading JSON:", err)
		}
		// Now add it to the datastore
		data.Store(benchmark)

		fmt.Fprintln(w, `{"success":true}`)
	})

	http.HandleFunc("/compare", func(w http.ResponseWriter, r *http.Request) {
		// Parameters are:
		// * a compared spec
		// * two filters on this spec
		// * a global filter

		spec := r.FormValue("spec")

		// The global filter
		filterJson := r.FormValue("filter")
		filter := MakeFilter(filterJson)

		// The individual filters
		valuesJSON := r.FormValue("values")
		var values []string
		err := json.Unmarshal([]byte(valuesJSON), &values)
		if err != nil {
			log.Println("Bad JSON received:", err)
			return
		}

		filters := make([]Filter, len(values))
		results := make([][][]benchbase.Benchmark, len(filters))
		for i, v := range values {
			filters[i] = AndFilter(filter, MakeSpecFilter(spec, v))
			results[i] = data.List(filters[i])
		}

		result := Compare(results...)

		enc := json.NewEncoder(w)
		err = enc.Encode(&result)
		if err != nil {
			log.Println("Error writing json:", err)
		}
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		filterJson := r.FormValue("filter")
		f := MakeFilter(filterJson)

		benchlist := data.List(f)
		enc := json.NewEncoder(w)
		err := enc.Encode(&benchlist)
		if err != nil {
			log.Println("Error writing json:", err)
		}
	})

	go func() {
		for _ = range time.Tick(30 * time.Second) {
			data.SaveToDisk(dataFilename)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			log.Println("Saving database...")
			data.SaveToDisk(dataFilename)
			log.Fatal("Exiting now.")
		}
	}()

	log.Println("Now listening on port 6666")
	log.Fatal(http.ListenAndServe(":6666", nil))
}
