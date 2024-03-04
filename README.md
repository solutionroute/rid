[![godoc](http://img.shields.io/badge/godev-reference-blue.svg?style=flat)](https://pkg.go.dev/github.com/mwyvr/rid?tab=doc)[![Test](https://github.com/mwyvr/rid/actions/workflows/test.yaml/badge.svg)](https://github.com/mwyvr/rid/actions/workflows/test.yaml)[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)![Coverage](https://img.shields.io/badge/coverage-92.6%25-brightgreen)

# rid

Package rid provides a performant, goroutine-safe generator of short
[k-sortable](https://en.wikipedia.org/wiki/K-sorted_sequence) unique IDs
suitable for use where inter-process ID generation coordination is not
required.

Using a non-standard character set (fewer vowels), IDs Base-32 encode as a
16-character URL-friendly, case-insensitive, representation like
`dfp7qt0v2pwt0v2x`.

An ID is a:

  - 4-byte timestamp value representing seconds since the Unix epoch, plus a
  - 6-byte random value; see the [Random Source](#random-source) discussion.

Built-in (de)serialization simplifies interacting with SQL databases and JSON.
`cmd/rid` provides the `rid` utility to generate or inspect IDs. Thanks to
fastrand[1], ID generation starts fast and remains so as cores are added.
De-serialization has also been optimized. See [Package
Benchmarks](#package-benchmarks).

Why `rid` as opposed to [alternatives](#package-comparisons)?

  - At 10 bytes binary, 16 bytes Base32 encoded, rid.IDs are case-insensitive
    and short, yet with 48 bits of uniqueness *per second*, are unique
    enough for many use cases.
  - IDs have a random component rather than potentially guessable
    monotonic counter.

_**Acknowledgement**: This package borrows heavily from rs/xid
(https://github.com/rs/xid), a zero-configuration globally-unique
high-performance ID generator which itself levers ideas from MongoDB
(https://docs.mongodb.com/manual/reference/method/ObjectId/)._

## Example:

```go
id := rid.New()
fmt.Printf("%s\n", id.String())
// Output: dfp7qt97menfv8ll

id2, err := rid.FromString("dfp7qt97menfv8ll")
if err != nil {
	fmt.Println(err)
}
fmt.Printf("%s %d %v\n", id2.Time(), id2.Random(), id2.Bytes())
// Output: 2022-12-28 09:24:57 -0800 PST 43582827111027 [99 172 123 233 39 163 106 237 162 115]
```

## CLI

Package `rid` also provides the `rid` tool for id generation and inspection. 

    $ rid 
	dfpb18y8dg90hc74

 	$ rid -c 2
	dfp9l9cgs05blztq
	dfp9l9d80yxdf804

    # produce 4 and inspect
	$ rid `rid -c 4`
	dfp9lmz9ksw87w48 ts:1672255955 rnd:256798116540552 2022-12-28 11:32:35 -0800 PST ID{ 0x63, 0xac, 0x99, 0xd3, 0xe9, 0x8e, 0x78, 0x83, 0xf0, 0x88 }
	dfp9lmxefym2ht2f ts:1672255955 rnd:190729433933902 2022-12-28 11:32:35 -0800 PST ID{ 0x63, 0xac, 0x99, 0xd3, 0xad, 0x77, 0xa8, 0x28, 0x68, 0x4e }
	dfp9lmt5zjy7km9n ts:1672255955 rnd: 76951796109621 2022-12-28 11:32:35 -0800 PST ID{ 0x63, 0xac, 0x99, 0xd3, 0x45, 0xfc, 0xbc, 0x78, 0xd1, 0x35 }
	dfp9lmxt5sms80m7 ts:1672255955 rnd:204708502569607 2022-12-28 11:32:35 -0800 PST ID{ 0x63, 0xac, 0x99, 0xd3, 0xba, 0x2e, 0x69, 0x94,  0x2, 0x87 }

## Random Source

Since cryptographically-secure IDs are not an objective for this package, other
approaches could be considered. With Go 1.19, `rid` utilized an internal runtime
`fastrand64` which provided single and multi-core performance benefits. Go
1.20 exposed `fastrand64` via the stdlib. As of rid v1.1.6, the package depends
on  Go 1.22 math/rand/v2 which provides Uint64N().

You may also enjoy reading:

- [Fast thread-safe randomness in Go](https://qqq.ninja/blog/post/fast-threadsafe-randomness-in-go/).
- For more information on fastrand (wyrand) see: https://github.com/wangyi-fudan/wyhash
 
To satisfy whether rid.IDs are unique enough for your use case, run
[eval/uniqcheck/main.go](eval/uniqcheck/main.go) with various values for number
of go routines and iterations, or, at the command line, produce 10,000,000 IDs
and use OS utilities to check:

    rid -c 10000000 | sort | uniq -d
    // None output

## Change Log

- 2023-03-02 v1.1.6: Package depends on math/rand/v2 and now requires Go 1.22+.
- 2023-01-23 Replaced stdlib Base32 with unrolled version for decoding performance.
- 2022-12-28 The "10byte" branch was merged to master; the "15byte-historical"
  branch will be left dormant. No major changes are now expected to this
  package; updates will focus on rounding out test coverage, addressing bugs,
  and clean up.

## Contributing

Contributions are welcome.

## Package Comparisons

Comparison table generated by [eval/compare/main.go](eval/compare/main.go):
| Package                                                   |BLen|ELen| K-Sort| Encoded ID and Next | Method | Components |
|-----------------------------------------------------------|----|----|-------|---------------------|--------|------------|
| [solutionroute/rid](https://github.com/solutionroute/rid) | 10 | 16 |  true | `dqjllq1sr0lrb93k`<br>`dqjllq40n8t6yx3r` | math/rand/v2 | 4 byte ts(sec) : 6 byte random |
| [rs/xid](https://github.com/rs/xid)                       | 12 | 20 |  true | `cnijjn34l33778tp3ing`<br>`cnijjn34l33778tp3io0` | counter | 4 byte ts(sec) : 2 byte mach ID : 2 byte pid : 3 byte monotonic counter |
| [segmentio/ksuid](https://github.com/segmentio/ksuid)     | 20 | 27 |  true | `2dCoMU7h8xaTlBhacYmys0CETTc`<br>`2dCoMQwFBYWoEi9mue6Kv1esiwy` | math/rand | 4 byte ts(sec) : 16 byte random |
| [google/uuid](https://github.com/google/uuid)             | 16 | 36 | false | `c9c5e75e-1ec1-4718-99d9-024c84422e27`<br>`e6387602-e980-44d2-a295-8c444d815d52` | crypt/rand | v4: 16 bytes random with version & variant embedded |
| [oklog/ulid](https://github.com/oklog/ulid)               | 16 | 26 |  true | `01HR3PM0YWWF2WK24EVAFHSTR8`<br>`01HR3PM0YWY56DH7AC0GX6YC4K` | crypt/rand | 6 byte ts(ms) : 10 byte counter random init per ts(ms) |
| [kjk/betterguid](https://github.com/kjk/betterguid)       | 17 | 20 |  true | `-Ns6PVER-PgtrwP2IgwR`<br>`-Ns6PVER-PgtrwP2IgwS` | counter | 8 byte ts(ms) : 9 byte counter random init per ts(ms) |

With only 48 bits of randomness per second, `rid` makes no attempt to weigh
in as a globally-unique ID generator. `rs/xid` is a feature comparable solid
alternative for such needs.

For a comparison of various Go-based unique ID solutions, see:
https://blog.kowalczyk.info/article/JyRZ/generating-good-unique-ids-in-go.html

## Package Benchmarks

A benchmark suite for the above noted packages can be found in
[eval/bench/bench_test.go](eval/bench/bench_test.go). All runs were done with scaling_governor set to `performance`:

    echo "performance" | sudo tee /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor

```
$ go test -cpu 1,2,4,8,16 -bench .
goos: linux
goarch: amd64
pkg: github.com/mwyvr/rid/eval/bench
cpu: AMD Ryzen 7 3800X 8-Core Processor             
BenchmarkRid              	19799767	       59.62 ns/op
BenchmarkRid-2            	23044785	       51.33 ns/op
BenchmarkRid-4            	43576563	       27.78 ns/op
BenchmarkRid-8            	60655580	       18.60 ns/op
BenchmarkRid-16           	57910780	       20.46 ns/op
BenchmarkXid              	22208196	       52.97 ns/op
BenchmarkXid-2            	36343159	      100.2 ns/op
BenchmarkXid-4            	20049046	       56.42 ns/op
BenchmarkXid-8            	27806431	       43.19 ns/op
BenchmarkXid-16           	56182581	       20.39 ns/op
BenchmarkKsuid            	2127823	      557.2 ns/op
BenchmarkKsuid-2          	1943469	      617.8 ns/op
BenchmarkKsuid-4          	2003335	      591.3 ns/op
BenchmarkKsuid-8          	1965288	      625.2 ns/op
BenchmarkKsuid-16         	1970457	      612.0 ns/op
BenchmarkGoogleUuid       	2211279	      539.6 ns/op
BenchmarkGoogleUuid-2     	3703702	      326.3 ns/op
BenchmarkGoogleUuid-4     	6998353	      182.9 ns/op
BenchmarkGoogleUuid-8     	12879537	       94.07 ns/op
BenchmarkGoogleUuid-16    	20498462	       60.23 ns/op
BenchmarkUlid             	 147532	     7795 ns/op
BenchmarkUlid-2           	 265328	     4460 ns/op
BenchmarkUlid-4           	 496263	     2434 ns/op
BenchmarkUlid-8           	 724950	     1606 ns/op
BenchmarkUlid-16          	 743781	     1667 ns/op
BenchmarkBetterguid       	13010038	       89.31 ns/op
BenchmarkBetterguid-2     	11938580	       93.20 ns/op
BenchmarkBetterguid-4     	7529553	      148.0 ns/op
BenchmarkBetterguid-8     	6638394	      179.2 ns/op
BenchmarkBetterguid-16    	4955013	      244.2 ns/op
PASS
```
