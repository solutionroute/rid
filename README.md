[![godoc](http://img.shields.io/badge/godev-reference-blue.svg?style=flat)](https://pkg.go.dev/github.com/solutionroute/rid?tab=doc)[![Test](https://github.com/solutionroute/rid/actions/workflows/test.yaml/badge.svg)](https://github.com/solutionroute/rid/actions/workflows/test.yaml)[![Go Coverage](https://img.shields.io/badge/coverage-98.3%25-brightgreen.svg?style=flat)](http://gocover.io/github.com/solutionroute/rid)[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# rid

Package `rid` provides a [k-sortable](https://en.wikipedia.org/wiki/K-sorted_sequence),
zero-configuration, unique ID generator. Binary IDs are Base32-encoded,
producing a 24-character case-insensitive URL-friendly representation like:
`062ekgz5k5f23ejagw2n7c9f`. Helper functions for Base64 encoding and decoding
are included.

Base32 encoding evenly aligns with 15 byte / 120 bit binary data. The 15-byte
binary representation of an ID is comprised of a:

- 6-byte timestamp value representing milliseconds since the Unix epoch
- 1-byte machine+process signature, derived from md5(machine ID + process ID)
- 6-byte random number using Go's runtime `fastrand` function. [1]

`rid` also implements a number of well-known interfaces to make use with JSON
and databases more convenient.

**Acknowledgement**: This package borrows _heavily_ from the at-scale capable
[rs/xid](https://github.com/rs/xid) package which itself levers ideas from
[MongoDB](https://docs.mongodb.com/manual/reference/method/ObjectId/).

Where this package differs, rid (15 bytes) | xid (12 bytes):

- 6-bytes of time, millisecond resolution | 4 bytes, second resolution
- 1-byte machine+process signature | 3 bytes machine ID, 2 bytes process ID
- 6-byte random number | 3-byte monotonic counter randomly initialized once 

## Usage

```go
	i := rid.New()
	fmt.Printf("%s\n", i)           // 062ekkxhmp31522vfjt7jv9t 
```

## Batteries included

`rid.ID` implements a number of common interfaces including:

- database/sql: driver.Valuer, sql.Scanner
- encoding: TextMarshaler, TextUnmarshaler
- encoding/json: json.Marshaler, json.Unmarshaler
- Stringer

Package `rid` also provides a command line tool `rid` allowing for id generation
and inspection. To install: `go install github.com/solutionroute/rid/cmd/...`

    $ rid 
    062ekjasgt18j0xgabq5zw45

    $ rid -c 2
    062ekjdxbc4yr0v0zyhv19zb
    062ekjdxbc4pesrn45jfz89k

    $ rid -c 2 -a  # use the alternate Base64 encoding:
    AYTuBDyA4ZXGYHMV7E1s
    AYTuBDyA4fh1zG-WbkUm

    # produce 4 and inspect
    $rid `rid -c 4`
    062ekjn39b2g7mvzwsxk2mx9 ts:1670369682250 rtsig:[0xc5] random:  4206918794033 | time:2022-12-06 15:34:42.25 -0800 PST ID{0x1,0x84,0xe9,0xca,0xa3,0x4a,0xc5,0x3,0xd3,0x7f,0xe6,0x7b,0x31,0x53,0xa9}
    062ekjn39b2tex8f39ht2vxk ts:1670369682250 rtsig:[0xc5] random:184121206399905 | time:2022-12-06 15:34:42.25 -0800 PST ID{0x1,0x84,0xe9,0xca,0xa3,0x4a,0xc5,0xa7,0x75,0xf,0x1a,0x63,0xa1,0x6f,0xb3}
    062ekjn39b2n2km1wn6qzaty ts:1670369682250 rtsig:[0xc5] random: 89397628587391 | time:2022-12-06 15:34:42.25 -0800 PST ID{0x1,0x84,0xe9,0xca,0xa3,0x4a,0xc5,0x51,0x4e,0x81,0xe5,0x4d,0x7f,0xab,0x5e}
    062ekjn39b2vxg1h326m5z9w ts:1670369682250 rtsig:[0xc5] random:209732666690882 | time:2022-12-06 15:34:42.25 -0800 PST ID{0x1,0x84,0xe9,0xca,0xa3,0x4a,0xc5,0xbe,0xc0,0x31,0x18,0x8d,0x42,0xfd,0x3c}

## Random Source

For random number generation `rid` uses a Go runtime `fastrand64` [1],
available in Go versions released post-spring 2022; it's non-deterministic,
goroutine safe, and fast.  For the purpose of *this* package, `fastrand64`
seems ideal.

Use of `fastrand` makes `rid` performant and scales well as cores/parallel
processes are added. While more testing will be done, no ID collisions have
been observed over numerous runs producing upwards of 300 million ID using
single and multiple goroutines.

[1] For more information on fastrand (wyrand) see: https://github.com/wangyi-fudan/wyhash
 and [Go's sources for runtime/stubs.go](https://cs.opensource.google/go/go/+/master:src/runtime/stubs.go;bpv=1;bpt=1?q=fastrand&ss=go%2Fgo:src%2Fruntime%2F).

## Package Comparisons

Comparison table generated by [eval/compare/main.go](eval/compare/main.go):

| Package                                                   |BLen|ELen| K-Sort| 0-Cfg | Encoded ID and Next | Method | Components |
|-----------------------------------------------------------|----|----|-------|-------|---------------------|--------|------------|
| [solutionroute/rid](https://github.com/solutionroute/rid)<br>Base32 (default) | 15 | 24 |  true |  true | `062exq1nk13qcc2bpjdbjx3b`<br>`062exq1nk13g08et49tn4h3f` | fastrand | 6 byte ts(ms) : 1 byte machine/pid signature : 8 byte random |
| [solutionroute/rid](https://github.com/solutionroute/rid)<br>Base64 (optional) | 15 | 20 |  true |  true | `AYTu3DWYR7jwhCW8ENzp`<br>`AYTu3DWYR1LQQgoG7P9B` | fastrand | 6 byte ts(ms) : 1 byte machine/pid signature : 8 byte random |
| [rs/xid](https://github.com/rs/xid)                       | 12 | 20 |  true |  true | `ce8hrfop26gfn8r71uh0`<br>`ce8hrfop26gfn8r71uhg` | counter | 4 byte ts(sec) : 2 byte mach ID : 2 byte pid : 3 byte monotonic counter |
| [segmentio/ksuid](https://github.com/segmentio/ksuid)     | 20 | 27 |  true |  true | `2IbeTgUpIvLr9xpVIYoD4pCrrM4`<br>`2IbeTmvvxbvLclWz7eSvQrRRTd5` | random | 4 byte ts(sec) : 16 byte random |
| [google/uuid](https://github.com/google/uuid)             | 16 | 36 | false |  true | `a1797151-5d1c-47c0-b92a-81af479372b2`<br>`2064a35e-f2f4-4580-bcda-15a651582b77` | crypt/rand | v4: 16 bytes random with version & variant embedded |
| [oklog/ulid](https://github.com/oklog/ulid)               | 16 | 26 |  true |  true | `01GKQDRDCRD1DC3DF73MP7H0S0`<br>`01GKQDRDCRWT49J9ZH681MNN4N` | crypt/rand | 6 byte ts(ms) : 10 byte counter random init per ts(ms) |
| [kjk/betterguid](https://github.com/kjk/betterguid)       | 17 | 20 |  true |  true | `-NIir2LN11DRU0HsEjcA`<br>`-NIir2LN11DRU0HsEjcB` | counter | 8 byte ts(ms) : 9 byte counter random init per ts(ms) |

If you don't need the k-sortable randomness this and other packages provide,
consider the well-tested and performant capable `rs/xid` package upon which 
`rid` is based. See https://github.com/rs/xid.

For a detailed comparison of various golang unique ID solutions, including `rs/xid`, see:
https://blog.kowalczyk.info/article/JyRZ/generating-good-unique-ids-in-go.html

## Package Benchmarks

A comparison with the above noted packages can be found in [eval/bench/bench_test.go](eval/bench/bench_test.go). Output:

### Intel 4-core Dell Latitude 7420 laptop

    $ go test -cpu 1,2,4,8 -benchmem  -run=^$   -bench  ^.*$ 
    goos: linux
    goarch: amd64
    pkg: github.com/solutionroute/rid/eval/bench
    cpu: 11th Gen Intel(R) Core(TM) i7-1185G7 @ 3.00GHz
    BenchmarkRid            	27007984	        41.90 ns/op	       0 B/op	       0 allocs/op
    BenchmarkRid-2          	54439544	        22.19 ns/op	       0 B/op	       0 allocs/op
    BenchmarkRid-4          	86903547	        13.66 ns/op	       0 B/op	       0 allocs/op
    BenchmarkRid-8          	132959510	         8.965 ns/op	       0 B/op	       0 allocs/op
    BenchmarkXid            	31221853	        37.28 ns/op	       0 B/op	       0 allocs/op
    BenchmarkXid-2          	35561181	        33.30 ns/op	       0 B/op	       0 allocs/op
    BenchmarkXid-4          	55113584	        27.53 ns/op	       0 B/op	       0 allocs/op
    BenchmarkXid-8          	71106020	        16.70 ns/op	       0 B/op	       0 allocs/op
    BenchmarkKsuid          	 3821538	       314.6 ns/op	       0 B/op	       0 allocs/op
    BenchmarkKsuid-2        	 3205950	       367.6 ns/op	       0 B/op	       0 allocs/op
    BenchmarkKsuid-4        	 3195728	       374.0 ns/op	       0 B/op	       0 allocs/op
    BenchmarkKsuid-8        	 3193402	       406.0 ns/op	       0 B/op	       0 allocs/op
    BenchmarkGoogleUuid     	 3561132	       334.0 ns/op	      16 B/op	       1 allocs/op
    BenchmarkGoogleUuid-2   	 4955325	       226.5 ns/op	      16 B/op	       1 allocs/op
    BenchmarkGoogleUuid-4   	 7119134	       160.6 ns/op	      16 B/op	       1 allocs/op
    BenchmarkGoogleUuid-8   	 7670070	       133.7 ns/op	      16 B/op	       1 allocs/op
    BenchmarkUlid           	  157432	      7488 ns/op	    5440 B/op	       3 allocs/op
    BenchmarkUlid-2         	  245198	      4758 ns/op	    5440 B/op	       3 allocs/op
    BenchmarkUlid-4         	  389346	      3082 ns/op	    5440 B/op	       3 allocs/op
    BenchmarkUlid-8         	  556402	      2108 ns/op	    5440 B/op	       3 allocs/op
    BenchmarkBetterguid     	13866880	        81.79 ns/op	      24 B/op	       1 allocs/op
    BenchmarkBetterguid-2   	10869040	       102.7 ns/op	      24 B/op	       1 allocs/op
    BenchmarkBetterguid-4   	 8379374	       138.3 ns/op	      24 B/op	       1 allocs/op
    BenchmarkBetterguid-8   	 6705165	       176.1 ns/op	      24 B/op	       1 allocs/op

### AMD 8-core desktop

    $ go test -cpu 1,2,4,8,16 -benchmem  -run=^$   -bench  ^.*$
    goos: linux
    goarch: amd64
    pkg: github.com/solutionroute/rid/eval/bench
    cpu: AMD Ryzen 7 3800X 8-Core Processor             
    BenchmarkRid              	19931982	        59.28 ns/op	       0 B/op	       0 allocs/op
    BenchmarkRid-2            	39499843	        29.90 ns/op	       0 B/op	       0 allocs/op
    BenchmarkRid-4            	78571719	        15.08 ns/op	       0 B/op	       0 allocs/op
    BenchmarkRid-8            	154435864	         7.715 ns/op	       0 B/op	       0 allocs/op
    BenchmarkRid-16           	279988606	         4.317 ns/op	       0 B/op	       0 allocs/op
    BenchmarkXid              	22248019	        52.32 ns/op	       0 B/op	       0 allocs/op
    BenchmarkXid-2            	37339971	       100.9 ns/op	       0 B/op	       0 allocs/op
    BenchmarkXid-4            	22793754	        52.99 ns/op	       0 B/op	       0 allocs/op
    BenchmarkXid-8            	43813854	        33.44 ns/op	       0 B/op	       0 allocs/op
    BenchmarkXid-16           	67285090	        16.89 ns/op	       0 B/op	       0 allocs/op
    BenchmarkKsuid            	 3252950	       362.7 ns/op	       0 B/op	       0 allocs/op
    BenchmarkKsuid-2          	 2198566	       783.1 ns/op	       0 B/op	       0 allocs/op
    BenchmarkKsuid-4          	 1443458	       832.5 ns/op	       0 B/op	       0 allocs/op
    BenchmarkKsuid-8          	 1441783	       838.8 ns/op	       0 B/op	       0 allocs/op
    BenchmarkKsuid-16         	 1407332	       857.6 ns/op	       0 B/op	       0 allocs/op
    BenchmarkGoogleUuid       	 2900432	       352.8 ns/op	      16 B/op	       1 allocs/op
    BenchmarkGoogleUuid-2     	 4841989	       214.2 ns/op	      16 B/op	       1 allocs/op
    BenchmarkGoogleUuid-4     	 9413534	       110.3 ns/op	      16 B/op	       1 allocs/op
    BenchmarkGoogleUuid-8     	18598221	        58.70 ns/op	      16 B/op	       1 allocs/op
    BenchmarkGoogleUuid-16    	29231677	        40.91 ns/op	      16 B/op	       1 allocs/op
    BenchmarkUlid             	  146024	      7890 ns/op	    5440 B/op	       3 allocs/op
    BenchmarkUlid-2           	  276771	      4396 ns/op	    5440 B/op	       3 allocs/op
    BenchmarkUlid-4           	  516540	      2348 ns/op	    5440 B/op	       3 allocs/op
    BenchmarkUlid-8           	  800305	      1476 ns/op	    5440 B/op	       3 allocs/op
    BenchmarkUlid-16          	  764970	      1484 ns/op	    5440 B/op	       3 allocs/op
    BenchmarkBetterguid       	14442830	        80.73 ns/op	      24 B/op	       1 allocs/op
    BenchmarkBetterguid-2     	10107843	       147.2 ns/op	      24 B/op	       1 allocs/op
    BenchmarkBetterguid-4     	 5394602	       260.7 ns/op	      24 B/op	       1 allocs/op
    BenchmarkBetterguid-8     	 4140949	       296.6 ns/op	      24 B/op	       1 allocs/op
    BenchmarkBetterguid-16    	 3173990	       379.9 ns/op	      24 B/op	       1 allocs/op

