package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Gyscos/benchbase"
)

func makeHandlers(data *Datastore) {
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
		// * a set of values for that spec
		// * a global filter

		spec := r.FormValue("spec")

		// The global filter
		filterJson := r.FormValue("filter")
		filter := MakeFilter(filterJson)
		// Main data source
		data := data.List(filter)

		// Specs to ignore when projecting
		ignoreJSON := r.FormValue("ignore")
		var ignores []string
		err := json.Unmarshal([]byte(ignoreJSON), &ignores)
		if err != nil {
			log.Println("Bad ignore JSON received:", err)
		}
		ignores = append(ignores, spec)

		// The individual filters
		valuesJSON := r.FormValue("values")
		var values []string
		err = json.Unmarshal([]byte(valuesJSON), &values)
		if err != nil {
			log.Println("Bad values JSON received:", err)
			return
		}
		// Dispatch according to spec value
		dispatched := Dispatch(data, spec, values)

		// Project onto similar configuration (except for spec)
		projected := Project(dispatched, ignores)

		enc := json.NewEncoder(w)
		err = enc.Encode(&projected)
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
}
