# go-mangler

[Documentation](https://pkg.go.dev/codeberg.org/gruf/go-mangler).

To put it simply is a bit of an odd library. It aims to provide incredibly fast, unique string outputs for all default supported input data types during a given runtime instance. See `mangler.String()`for supported types.

It is useful, for example, for use as part of larger abstractions involving hashmaps. That was my particular usecase anyways...

This package does make liberal use of the "unsafe" package.

Benchmarks are below. Please note the more important thing to notice here is the relative difference in benchmark scores, the actual `ns/op`,`B/op`,`allocs/op` accounts for running through ~80 possible test cases, including some not-ideal situations.

The choice of libraries in the benchmark are just a selection of libraries that could be used in a similar manner to this one, i.e. serializing in some manner.

```
$ go test -run=none -benchmem -gcflags=all='-l=4' -bench=.*
goos: linux
goarch: amd64
pkg: codeberg.org/gruf/go-mangler/v2
cpu: Intel(R) Core(TM) i7-9700T CPU @ 2.00GHz
BenchmarkMangle
BenchmarkMangle-8                         950618              1238 ns/op               0 B/op          0 allocs/op
BenchmarkMangleMulti
BenchmarkMangleMulti-8                    977418              1323 ns/op               0 B/op          0 allocs/op
BenchmarkMangleKnown
BenchmarkMangleKnown-8                   1961787               603.4 ns/op             0 B/op          0 allocs/op
BenchmarkJSON
BenchmarkJSON-8                           266931              4470 ns/op            2730 B/op        146 allocs/op
BenchmarkLoosy
BenchmarkLoosy-8                          353376              3643 ns/op            1008 B/op         98 allocs/op
BenchmarkFmt
BenchmarkFmt-8                            202864              7582 ns/op            1472 B/op        136 allocs/op
BenchmarkFxmackerCbor
BenchmarkFxmackerCbor-8                   419496              2855 ns/op            1640 B/op        150 allocs/op
BenchmarkMitchellhHashStructure
BenchmarkMitchellhHashStructure-8         113822             10711 ns/op           12046 B/op       1270 allocs/op
BenchmarkCnfStructhash
BenchmarkCnfStructhash-8                    8990            150967 ns/op          280594 B/op       3912 allocs/op
PASS
ok      codeberg.org/gruf/go-mangler/v2
```
