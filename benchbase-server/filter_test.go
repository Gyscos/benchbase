package main

import (
	"testing"

	"github.com/Gyscos/benchbase"
)

func TestFilter(t *testing.T) {
	b := benchbase.NewBenchmark()
	b.Conf["rev"] = "118"

	f := MakeFilter("rev>=110")

	if !f(b) {
		t.Error("Should accept the host.")
	}

	b.Conf["rev"] = "98"
	if f(b) {
		t.Error("Should refuse the host.")
	}
}
