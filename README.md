## Crypto Testing in Golang

These simple tools test the efficiency of various different hashing algorithms implemented in golang.

The first tool `crypto-tester` works by running different algorithms against the same file and outputting the results.
The second tool `sha3sum` will only run a single algorithm once and output the time taken.

### Building

The easiest way to get this is using `go get` and `go install`:

```
go get -u github.com/anonymouse64/crypto_test_golang/...
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

```
pi@raspberrypi:~ $ $GOPATH/bin/crypto-tester --help
Usage of /home/pi/go/bin/crypto-tester:
  -escape-quote
    	whether to escape quotes in the json output
  -file string
    	file to hash (random file is generated if this is omitted)
  -format string
    	format to output (json or yaml) (default "json")
  -size int
    	size of generated random file (default 10)
  -unit string
    	units to use (possible values : ns, us, ms, s) (default "ns")
```

This utility is intended to benchmark multiple different crypto hashing algorithms built into golang. You can specify what file to run against with the `-file=file.txt` option. If this option is omitted, a random 10 MB file will be generated and used (it is generated from `crypto/rand`'s Read: https://golang.org/pkg/crypto/rand/#Read.

The output can be output in either JSON or YAML. There is also an option to escape quotes if necessary. The units default to nanoseconds, but can be changed with the `-unit` option.

There are currently 14 different algorithms tested on the file in 2 different ways.
The first way is by just running the `.Sum` method on a pre-read byte array from the file.
The second way is by using snapcore's `osutil.FileDigest` function. 

Example with a specified file:

```
pi@raspberrypi:~ $ $GOPATH/bin/crypto-tester -file=file.txt- algorithm: md4
  byteshash: 15417
  filehash: 19789507
- algorithm: md5
  byteshash: 12344
  filehash: 6573204
- algorithm: sha1
  byteshash: 11980
  filehash: 15064057
- algorithm: sha256
  byteshash: 36198
  filehash: 41202966
- algorithm: sha256_224
  byteshash: 15624
  filehash: 20859607
- algorithm: sha512
  byteshash: 53437
  filehash: 34371013
- algorithm: sha512_224
  byteshash: 18958
  filehash: 36242777
- algorithm: sha512_256
  byteshash: 19687
  filehash: 33092425
- algorithm: sha384
  byteshash: 18750
  filehash: 32011127
- algorithm: ripemd160
  byteshash: 14011
  filehash: 25559693
- algorithm: sha3_224
  byteshash: 41719
  filehash: 23439701
- algorithm: sha3_256
  byteshash: 17083
  filehash: 25225267
- algorithm: sha3_384
  byteshash: 19687
  filehash: 31246859
- algorithm: sha3_512
  byteshash: 19062
  filehash: 43810091

```



#### sha3sum

```
pi@raspberrypi:~ $ $GOPATH/bin/sha3sum --help
Usage of /home/pi/go/bin/sha3sum:
  -alg string
    	algorithm to use, sha3_512, sha3_384, sha3_256, or sha3_224 (default "sha3_512")
  -file string
    	file to hash
  -size int
    	size of generated random file (default 10)
  -unit string
    	units to use (possible values : ns, us, ms, s) (default "s")
```

This utility is intended to benchmark specifically sha3 (or various variants of it) against a file, either randomly generated or specified with the `-file` option. The output unit can be changed with the `unit` file, and the specific algorithm to use can be specified with the `-alg` option.

Use a randomly generated file:

```
pi@raspberrypi:~ $ $GOPATH/bin/sha3sum 
6b792386272d67d01424e7b1714815289b95bf749a24f1c68b21825ccee270c788fcef7f51266e791b2895e0999ac886e828cae70f4e54d0ca0c891d2c6e41c1 /tmp/sha3sum_example107825369
Calculated in 1.550840 sec,  6.45 MBps
```

Specify a file:

```
pi@raspberrypi:~ $ $GOPATH/bin/sha3sum -file=file.txt
bb8f2e943d20d37d0984c07038a34d71f4b6ad67db5124cb28097e89365de9aa163e5f3c18e0f4227a6d4c0ab03a97b44154d4fccf957264014fa9de4614d56a file.txt
Calculated in 0.068106 sec,  3.96 MBps
```

Specify the size of the randomly generated file in Megabytes:
```
pi@raspberrypi:~ $ $GOPATH/bin/sha3sum -size=100
e01db5396e44d47435c81d4b0647463f43ee5923ce24013b714563868abfde8434ab974857d5eb2ccab6ee9bf51e34a730f3c29bf8a9c0ef995e1ddbe49a8d21 /tmp/sha3sum_example556721470
Calculated in 15.340368 sec,  6.52 MBps

```

