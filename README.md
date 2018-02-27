## Crypto Testing in Golang

These simple tools test the efficiency of various different hashing algorithms implemented in golang.

The first tool `crypto-tester` works by running different algorithms against the same file and outputting the results.
The second tool `sha3sum` will only run a single algorithm once and output the time taken.

### Building

The easiest way to get this is using `go get` and `go install`:

```
go get github.com/anonymouse64/crypto_test_golang/...
go install github.com/anonymouse64/crypto_test_golang/...
```

This will install both of the included utilities into your `$GOPATH/bin`.

The tool can also be built locally using `go build`. The dependencies must be resolved and can be done using `go get`:

```
git clone https://github.com/anonymouse64/crypto_test_golang.git
cd crypto_test_golang
go get -d -t ./...
go build ./cmd/sha3sum/
go build ./cmd/crypto-tester/
```


### Usage

#### crypto-tester

This utility is intended to benchmark multiple different crypto hashing algorithms built into golang. You can specify what file to run against with the `-file=file.txt` option. If this option is omitted, a random 10 MB file will be generated and used (it is generated from `crypto/rand`'s Read: https://golang.org/pkg/crypto/rand/#Read.

There are currently 14 different algorithms tested on this file in 2 different ways.
The first way is by just running the `.Sum` method on a pre-read byte array from the file.
The second way is by using snapcore's `osutil.FileDigest` function. 

Example:

```
$GOPATH/bin/crypto-tester -file=file.txt
```

#### sha3sum

This utility is intended to benchmark specifically sha3 (or various variants of it) against a file, either randomly generated or specified with the `-file` option.

Example:

```
$GOPATH/bin/sha3sum
```