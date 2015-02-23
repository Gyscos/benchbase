package benchbase

import (
	"bytes"
	"sort"
)

// Configuration describes the kind of test used in a benchmark.
// The configuration defines the fields in the result times, so
// two results can only be compared if they have the same configuration.
type Configuration map[string]string

func NewConfiguration() Configuration {
	return make(map[string]string)
}

func (c *Configuration) Hash() string {
	return c.PartialHash()
}

func (c *Configuration) PartialHash(ignoredSpecs ...string) string {
	var b bytes.Buffer

	var keys []string

	for key, _ := range *c {
		if !contains(ignoredSpecs, key) {
			keys = append(keys, key)
		}
	}

	sort.StringSlice(keys).Sort()

	for _, key := range keys {
		value := (*c)[key]

		b.WriteString(key)
		b.WriteString("=")
		b.WriteString(value)
		b.WriteString(";")
	}

	return b.String()
}

func contains(list []string, value string) bool {
	for _, s := range list {
		if s == value {
			return true
		}
	}
	return false
}
