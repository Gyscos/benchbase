package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Gyscos/benchbase"
)

func sendError(w io.Writer, err string) error {
	enc := json.NewEncoder(w)

	return enc.Encode(struct {
		Error string
	}{
		err,
	})
}

func sendResult(w io.Writer, result interface{}) error {
	enc := json.NewEncoder(w)

	return enc.Encode(struct {
		Error  error
		Result interface{}
	}{
		nil,
		result,
	})
}

func setupHandlers(data *Datastore) {
	r := http.NewServeMux()

	r.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {

		var benchmark benchbase.Benchmark

		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&benchmark)
		if err != nil {
			log.Println("Error reading JSON:", err)
			sendError(w, "Could not read the benchmark: "+err.Error())
		}
		// Now add it to the datastore
		data.Store(&benchmark)

		sendResult(w, "OK")
	})

	r.HandleFunc("/compare", func(w http.ResponseWriter, r *http.Request) {
		// Parameters are:
		// * a compared spec
		// * a set of values for that spec
		// * a global filter

		spec := r.FormValue("spec")
		maxString := r.FormValue("max")
		max, _ := strconv.Atoi(maxString)
		sortJSON := r.FormValue("sort")
		var ordering []string
		if sortJSON != "" {
			err := json.Unmarshal([]byte(sortJSON), &ordering)
			if err != nil {
				log.Println("Bad sort JSON received:", err)
			}
		}

		// The global filter
		filterJson := r.FormValue("filter")
		filter := MakeFilter(filterJson)
		// Main data source
		data := data.List(filter, ordering, max)

		// Specs to ignore when projecting
		ignoreJSON := r.FormValue("ignore")
		var ignores []string
		if ignoreJSON != "" {
			err := json.Unmarshal([]byte(ignoreJSON), &ignores)
			if err != nil {
				log.Println("Bad ignore JSON received:", err)
			}
		}
		ignores = append(ignores, spec)

		// The individual filters
		valuesJSON := r.FormValue("values")
		var values []string
		err := json.Unmarshal([]byte(valuesJSON), &values)
		if err != nil {
			log.Println("Bad values JSON received:", err)
			sendError(w, "Error reading values: "+err.Error())
			return
		}
		// Dispatch according to spec value
		dispatched := Dispatch(data, spec, values)

		// Project onto similar configuration (except for spec)
		projected := Project(dispatched, ignores)
		err = sendResult(w, projected)
		if err != nil {
			log.Println("Error writing json:", err)
		}
	})

	r.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		filterJson := r.FormValue("filter")
		f := MakeFilter(filterJson)

		maxString := r.FormValue("max")
		max, _ := strconv.Atoi(maxString)
		sortJSON := r.FormValue("sort")
		var ordering []string
		if sortJSON != "" {
			err := json.Unmarshal([]byte(sortJSON), &ordering)
			if err != nil {
				log.Println("Bad sort JSON received:", err)
			}
		}

		benchlist := data.List(f, ordering, max)
		err := sendResult(w, benchlist)
		if err != nil {
			log.Println("Error writing json:", err)
		}
	})

	http.Handle("/", &MyServer{r})
}

// Small wrapper to handle cross-origin
type MyServer struct {
	r http.Handler
}

func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}
