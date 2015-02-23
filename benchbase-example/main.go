package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Gyscos/benchbase"
)

func main() {
	bench := benchbase.NewBenchmark()

	bench.Conf["host"] = "alice"
	bench.Result["total"] = 2.006
	err := benchbase.PostBenchmark("http://localhost:6666", bench)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(&bench)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))

}
