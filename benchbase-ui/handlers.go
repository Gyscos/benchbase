package main

import "net/http"

func setupHandlers() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
	})

	http.HandleFunc("/compare", func(w http.ResponseWriter, r *http.Request) {
	})
}
