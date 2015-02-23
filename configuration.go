package benchbase

import "bytes"

// Configuration describes the kind of test used in a benchmark.
// The configuration defines the fields in the result times, so
// two results can only be compared if they have the same configuration.
type Configuration map[string]string

func NewConfiguration() Configuration {
	return make(map[string]string)
}

func (c *Configuration) Hash() string {
	var b bytes.Buffer

	for key, value := range *c {
		b.WriteString(key)
		b.WriteString("=")
		b.WriteString(value)
		b.WriteString(";")
	}

	return b.String()
}

func (c *Configuration) PartialHash(ignoredSpec string) string {
	var b bytes.Buffer

	for key, value := range *c {
		if key == ignoredSpec {
			continue
		}

		b.WriteString(key)
		b.WriteString("=")
		b.WriteString(value)
		b.WriteString(";")
	}

	return b.String()
}
