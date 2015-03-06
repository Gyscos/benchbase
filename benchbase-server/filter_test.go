package main

import (
	"testing"

	"github.com/Gyscos/benchbase"
)

func TestFilter(t *testing.T) {
	b := benchbase.NewBenchmark()
	b.Conf["rev"] = "118"

	f, err := MakeFilters(`{"Rev":">=118"}`)
	if err != nil {
		t.Error(err)
	}

	if !f["Rev"](b.Conf["Rev"]) {
		t.Error("Should accept the host.")
	}

	b.Conf["rev"] = "98"
	if f(b) {
		t.Error("Should refuse the host.")
	}
}
