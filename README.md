## Crypto Testing in Golang

This simple tool tests the efficiency of various different hashing algorithms implemented in golang.

### Building

The tool can be built with a simple `go build cmd/crypto-tester/main.go`, assuming the dependencies have been fulfilled. The dependencies can be fulfilled by running `go get -d -t ./...`, which will automatically install all of the depedencies for you.

### Usage

To run the test, just provide a single command line argument `-file=file.txt`. This file is read into a byte array and the various hashing algorithms are run against this file, which are then output to stdout via json.

There is an example file included with this, which is just a short snipped of text from the golang reference : https://golang.org/pkg/hash/#Hash