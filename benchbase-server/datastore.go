package main

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/Gyscos/benchbase"
)

// Datastore stores Benchmark and allow to get then through filters.
// Currently implemented as a simple list of benchmarks for each configuration.
type Datastore struct {
	data  map[string][]*benchbase.Benchmark
	mutex sync.RWMutex
}

type DatastoreDump struct {
	Data []*benchbase.Benchmark
}

// Returns a new empty datastore
func NewDatastore() *Datastore {
	return &Datastore{
		data: make(map[string][]*benchbase.Benchmark),
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
	for _, b := range flat.Data {
		d.Store(b)
	}

	return d, nil
}

func GetDatastore(filename string) *Datastore {
	d, err := LoadDatastore(filename)
	if err != nil {
		return NewDatastore()
	} else {
		return d
	}
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
	flat := d.List(TrueFilter, nil, 0)

	err = enc.Encode(&DatastoreDump{flat})

	return err
}

// Stores a benchmark in the database
func (d *Datastore) Store(benchmark *benchbase.Benchmark) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// We sort benchmark by configuration
	hash := benchmark.Conf.Hash()
	d.data[hash] = append(d.data[hash], benchmark)
}

// Return a list of benchmarks sorted by configuration
func (d *Datastore) List(filter Filter, ordering []string, max int) []*benchbase.Benchmark {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	result := make([]*benchbase.Benchmark, 0)
	for _, list := range d.data {
		for _, b := range list {
			if filter(b) {
				result = append(result, b)
			}
		}
	}
	return result
}
