package benchbase

// Subject describes the subject of a benchmark.
type Subject map[string]string

func NewSubject() Subject {
	return make(map[string]string)
}
