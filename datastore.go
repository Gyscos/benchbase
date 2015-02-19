package main

import (
	"encoding/json"
	"os"
	"sync"
)

type Datastore struct {
	data  map[Configuration][]BenchMark
	mutex sync.RWMutex
}

type DatastoreDump struct {
	data [][]BenchMark
}

// Returns a new empty datastore
func NewDatastore() *Datastore {
	return &Datastore{
		data: make(map[Configuration][]BenchMark),
	}
}

func LoadDatastore(filename string) (*Datastore, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Decompress?

	dec := json.NewDecoder(f)
	var flat DatastoreDump
	err = dec.Decode(&flat)
	if err != nil {
		return nil, err
	}

	d := NewDatastore()
	for _, list := range flat.data {
		for _, b := range list {
			d.Store(b)
		}
	}

	return d, nil
}

func (d *Datastore) SaveToDisk(filename string) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	f, err := os.Create(filename)
	if err != nil {
		return nil
	}
	defer f.Close()

	// Compress?

	enc := json.NewEncoder(f)
	flat := d.List(TrueFilter)

	err = enc.Encode(&DatastoreDump{flat})

	return err
}

// Stores a benchmark in the database
func (d *Datastore) Store(benchmark BenchMark) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// We sort benchmark by configuration
	d.data[benchmark.Conf] = append(d.data[benchmark.Conf], benchmark)
}

// Return a list of benchmarks sorted by configuration
func (d *Datastore) List(filter Filter) [][]BenchMark {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	result := make([][]BenchMark, 0)
	for _, list := range d.data {
		tmp := make([]BenchMark, 0)
		for _, b := range list {
			if filter(b) {
				tmp = append(tmp, b)
			}
		}
		if len(tmp) > 0 {
			result = append(result, tmp)
		}
	}
	return result
}
