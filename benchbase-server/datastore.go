package main

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/Gyscos/benchbase"
)

type SpecMap map[string]*SpecList

func (m *SpecMap) add(b *benchbase.Benchmark, id int) {
	for spec, value := range b.Conf {
		l := (*m)[spec]
		if l == nil {
			l = &SpecList{}
			(*m)[spec] = l
		}

		l.add(value, id)
	}
}

type SpecList struct {
	entries [][]int
	values  []string
}

func (s *SpecList) insert(i int, value string, id int) {
	s.entries = append(s.entries, nil)
	copy(s.entries[i+1:], s.entries[i:])
	s.entries[i] = []int{id}

	s.values = append(s.values, "")
	copy(s.values[i+1:], s.values[i:])
	s.values[i] = value
}

func (s *SpecList) add(value string, id int) {
	for i, v := range s.values {

		if v == value {
			s.entries[i] = append(s.entries[i], id)
			return
		} else if v > value {
			// Insert
			s.insert(i, value, id)
			return
		}
	}

	s.entries = append(s.entries, []int{id})
	s.values = append(s.values, value)
}

// Datastore stores Benchmark and allow to get then through filters.
// Currently implemented as a simple list of benchmarks for each configuration.
type Datastore struct {
	data    []*benchbase.Benchmark
	indices SpecMap

	mutex sync.RWMutex
}

type DatastoreDump struct {
	Data []*benchbase.Benchmark
}

// Returns a new empty datastore
func NewDatastore() *Datastore {
	return &Datastore{
		indices: make(map[string]*SpecList),
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

	err = enc.Encode(&DatastoreDump{d.data})

	return err
}

// Stores a benchmark in the database
func (d *Datastore) Store(benchmark *benchbase.Benchmark) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	i := len(d.data)
	d.data = append(d.data, benchmark)

	d.indices.add(benchmark, i)
}

// FilterIDs return a subset of the given ids, passing the filters, and sorted by the ordering list.
func (d *Datastore) FilterIDs(list []int, filters map[string]SpecFilter, ordering []string, max int) []int {
	if len(list) == 0 {
		return nil
	}

	// We are at the bottom of the recursion: just return the current candidates
	if len(ordering) == 0 {
		if max > 0 && len(list) > max {
			return list[:max]
		} else {
			return list
		}
	}

	// We will look at the sorted keys for the first spec, and recurse.
	spec := ordering[0]
	keys := d.indices[spec].values
	var keyIds []int
	// The filters can reduce the number of available keys.
	filter := filters[spec]
	if filter != nil {
		keyIds = filter(keys)
	} else {
		keyIds = invert(consecutive(0, len(keys)))
	}

	var result []int

	for _, keyId := range keyIds {
		// We go backward to get the [max] latest.
		newMax := 0
		if max > 0 {
			newMax = max - len(result)
		}
		r := d.FilterIDs(intersection(list, d.indices[spec].entries[keyId]), filters, ordering[1:], newMax)
		result = append(result, r...)
		if max > 0 && len(result) == max {
			break
		}
	}

	return result
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func (d *Datastore) completeOrdering(ordering []string) []string {
	result := append(make([]string, 0, len(d.indices)), ordering...)

	for spec, _ := range d.indices {
		if !contains(ordering, spec) {
			result = append(result, spec)
		}
	}

	return result
}

func (d *Datastore) List(filters map[string]SpecFilter, ordering []string, max int) []*benchbase.Benchmark {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	ids := consecutive(0, len(d.data))
	completeOrder := d.completeOrdering(ordering)
	ids = d.FilterIDs(ids, filters, completeOrder, max)
	result := make([]*benchbase.Benchmark, len(ids))
	for i, id := range ids {
		result[i] = d.data[id]
	}
	return result
}
