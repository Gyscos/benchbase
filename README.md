Benchbase
=========

Benchbase is a library for setting up and accessing a benchmark database.
The benchbase-server runs a HTTP server, listening for requests.

A benchmark is defined as a configuration and results. The configuration is a
map[string]string, while a result is a map[string]float64.

Listing benchmark can include a filter, that will restrict the results to the
given configuration values.

Comparing benchmark takes a configuration key, a list of possible values, and
an optionnal global filter. The result will be, for each unique configuration,
a list of results for each value given.
