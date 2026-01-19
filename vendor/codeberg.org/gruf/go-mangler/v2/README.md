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
cpu: AMD Ryzen 7 7840U w/ Radeon  780M Graphics
BenchmarkMangle-16                       3229830               371.4 ns/op             0 B/op          0 allocs/op
BenchmarkMangleMulti-16                  3235609               370.3 ns/op             0 B/op          0 allocs/op
BenchmarkMangleKnown-16                  7368690               162.2 ns/op             0 B/op          0 allocs/op
BenchmarkJSON-16                          734290              1653 ns/op            2334 B/op        113 allocs/op
BenchmarkLoosy-16                        1117132              1074 ns/op             768 B/op         70 allocs/op
BenchmarkFmt-16                           413862              2487 ns/op            1514 B/op        146 allocs/op
BenchmarkFxmackerCbor-16                 1089543              1160 ns/op            1610 B/op        113 allocs/op
BenchmarkMitchellhHashStructure-16        184632              6109 ns/op           11461 B/op       1145 allocs/op
BenchmarkCnfStructhash-16                  12608             97234 ns/op          275670 B/op       3799 allocs/op
PASS
ok      codeberg.org/gruf/go-mangler/v2 16.613s
```
