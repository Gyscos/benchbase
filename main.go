// BenchBase is a benchmark database.
// It offers a simple HTTP API to push new benchmark results, and retrieve / sort / filter existing ones.

// Actions:
// Store: Benchmark
// Compare: (A, B ConfigurationFilter, SituationFilter) [][2]Benchmark
// List: (A ConfigurationFilter, SituationFilter) []Benchmark
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	data := NewDatastore()

	http.HandleFunc("push", func(w http.ResponseWriter, r *http.Request) {
		benchJson := r.FormValue("benchmark")
		var benchmark BenchMark
		err := json.Unmarshal([]byte(benchJson), &benchmark)
		if err != nil {
			log.Println("Error reading JSON:", err)
		}
		// Now add it to the datastore
		data.Store(benchmark)

	})
	http.HandleFunc("compare", func(w http.ResponseWriter, r *http.Request) {

	})
	http.HandleFunc("list", func(w http.ResponseWriter, r *http.Request) {
		filterJson := r.FormValue("filter")
		f := MakeFilter(filterJson)

		benchlist := data.List(f)
		enc := json.NewEncoder(w)
		err := enc.Encode(&benchlist)
		if err != nil {
			log.Println("Error writing json:", err)
		}
	})
}
