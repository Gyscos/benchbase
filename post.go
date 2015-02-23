package benchbase

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func PostBenchmark(host string, benchmark *Benchmark) error {
	b, err := json.Marshal(benchmark)
	if err != nil {
		return err
	}
	resp, err := http.Post(host+"/push", "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
