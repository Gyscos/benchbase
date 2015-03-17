package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var filename string

	flag.StringVar(&filename, "f", "", "Database file.")

	flag.Parse()

	b, err := read(filename)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "    ")

	out.WriteTo(os.Stdout)
}

func read(filename string) ([]byte, error) {
	b, err := readDatabase(filename, true)
	if err == nil {
		return b, nil
	}

	return readDatabase(filename, false)
}

func readDatabase(filename string, compress bool) ([]byte, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var r io.Reader
	if compress {
		g, err := gzip.NewReader(f)
		if err != nil {
			// Cancel
			return nil, err
		}
		defer g.Close()

		r = g
	} else {
		r = f
	}

	return ioutil.ReadAll(r)
}
