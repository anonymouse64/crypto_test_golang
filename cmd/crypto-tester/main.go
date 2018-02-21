package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"flag"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"strings"
	"time"

	md4 "golang.org/x/crypto/md4"
	ripemd160 "golang.org/x/crypto/ripemd160"
	sha3 "golang.org/x/crypto/sha3"
)

func timeHash(hasher hash.Hash, bytes []byte) ([]byte, time.Duration) {
	// Get start time
	start := time.Now()

	// Compute the time
	hashRes := hasher.Sum(bytes)

	// Return the hash result and the time since the start of the function
	return hashRes, time.Since(start)
}

type HasherType struct {
	Name   string
	Hasher hash.Hash
}

type HashResult struct {
	Algorithm string `json:"alg"`
	TimeNs    int64  `json:"time_ns"`
}

func main() {
	// Get the file from the flag
	fileStr := flag.String("file", "file.txt", "file to hash")

	// Parse command line flags
	flag.Parse()

	// Make all the hasher objects
	hasherTable := initializeHasherTable()

	// Read in the file as a byte array
	bytes, err := ioutil.ReadFile(*fileStr)
	if err != nil {
		log.Fatalf("Error reading a file : %+v\n", err)
	}

	var results []HashResult
	results = make([]HashResult, len(hasherTable))

	// For each hasher run the hash and print off how long it took
	for index, hasher := range hasherTable {
		// Compute the hash for each one
		_, timeElapsed := timeHash(hasher.Hasher, bytes)

		results[index] = HashResult{
			Algorithm: hasher.Name,
			TimeNs:    int64(timeElapsed / time.Nanosecond),
		}
	}

	jsonbytes, err := json.Marshal(results)
	if err != nil {
		log.Fatalf("failed to encode json")
	}

	fmt.Println(strings.Replace(string(jsonbytes[:]), "\"", "\\\"", -1))
}

func initializeHasherTable() []HasherType {
	return []HasherType{
		{"md4", md4.New()},
		{"md5", md5.New()},
		{"sha1", sha1.New()},
		{"sha256", sha256.New()},
		{"sha256_224", sha256.New224()},
		{"sha512", sha512.New()},
		{"sha512_224", sha512.New512_224()},
		{"sha512_256", sha512.New512_256()},
		{"sha512", sha512.New()},
		{"sha384", sha512.New384()},
		{"ripemd160", ripemd160.New()},
		{"sha3_224", sha3.New224()},
		{"sha3_256", sha3.New256()},
		{"sha3_384", sha3.New384()},
		{"sha3_512", sha3.New512()},
		// Disabled for now...
		// {"blake2s128", blake2s.New128(nil)},
		// {"blake2s256", blake2s.New256(nil)},
		// {"blake2b", blake2b.New(128, nil)},
		// {"blake2b256", blake2b.New256(nil)},
		// {"blake2b384", blake2b.New384(nil)},
		// {"blake2b512", blake2b.New512(nil)},
	}
}
