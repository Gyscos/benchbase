package main

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
)

func LoadDatastore(filename string) (*Datastore, error) {
	d, err := load(filename, true)
	if err == nil {
		return d, nil
	}

	return load(filename, false)
}

func load(filename string, compress bool) (*Datastore, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Decompress?
	var r io.Reader
	if compress {
		r, err = gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
	} else {
		r = f
	}

	dec := json.NewDecoder(r)
	var flat DatastoreDump
	err = dec.Decode(&flat)
	if err != nil {
		return nil, err
	}

	d := NewDatastore()
	for _, b := range flat.Data {
		d.Store(b)
	}

	return d, nil
}

func (d *Datastore) SaveToDisk(filename string, compress bool) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	f, err := os.Create(filename)
	if err != nil {
		return nil
	}
	defer f.Close()

	// Compress?
	var w io.WriteCloser
	if compress {
		w = gzip.NewWriter(f)
		defer w.Close()
	} else {
		w = f
	}

	enc := json.NewEncoder(w)

	err = enc.Encode(&DatastoreDump{d.data})

	return err
}
